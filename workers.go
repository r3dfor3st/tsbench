package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// workersType stores array of jobs channels, a results channel, done channel and a pointer
// to a structure containing the commandline params. Jobs channel is read by workers. When
// workers receive an array of parameters from the CSV file, they execute SQL query and hand
// over its execution time to the results channel. Once each worker is done processing (EOF
// is found), it sends a message to the done channel. This channel is used by the results
// channel listener to know when to exit the loop.
type workersType struct {
	jobsChannels   []chan []string
	resultsChannel chan time.Duration
	doneChannel    chan struct{}
	clip           *cliParamsType
}

// initWorkers initializes all lchannels that will be used by the async workers.
func initWorkers(clip *cliParamsType) (workers workersType) {
	workers.clip = clip
	workers.jobsChannels = make([]chan []string, clip.workers)
	for i := uint(0); i < clip.workers; i++ {
		workers.jobsChannels[i] = make(chan []string)
	}
	workers.resultsChannel = make(chan time.Duration)
	workers.doneChannel = make(chan struct{})
	return
}

/*
run starts all async workers (runner), asynchronously executes the stream processing callback
function (paramsProvider) and waits to collect all benchmarking results from the results channel.
paramsProvider function feeds the data to the worker (read from the CSV file). Check traverseInput
function for the details of how the CSV file is parsed. Once reading the CSV file is done (by some
error or EOF), worke channels are automatically closed. Also, when workers are done, they signal
this on the done channel.
*/
func (workers *workersType) run(paramsProvider func(), runner func(params <-chan []string, i int)) (results measurementsType, err error) {
	for i, jobChannel := range workers.jobsChannels {
		go func(jc <-chan []string, i int) {
			runner(jc, i)
			workers.doneChannel <- struct{}{}
		}(jobChannel, i)
	}
	go func() {
		paramsProvider()
		for _, jc := range workers.jobsChannels {
			close(jc)
		}
	}()

	workersDone := 0
	for {
		select {
		case m := <-workers.resultsChannel:
			results = append(results, m)
		case <-workers.doneChannel:
			workersDone++
			if workersDone == len(workers.jobsChannels) {
				return
			}
		case <-time.After(10 * time.Second):
			err = errors.New("results channel timed out")
			return
		}
	}
}

// benchmarkQueryExecution takes an array of parameters and executes SQL query with them.
// It registers the execution time and reports it to the results channel.
func (workers *workersType) benchmarkQueryExecution(params <-chan []string, i int) {
	var executionTime time.Duration
	var err error
	var db *sql.DB

	db, err = createDbConnection(workers.clip.connectionString())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer db.Close()

	for record := range params {
		executionTime, err = runQuery(db, workers.clip.sql, record)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		log.Println(i, executionTime)
		workers.resultsChannel <- executionTime
	}
}

// createDbConnection creates postgresql database connection from the provided connection string.
func createDbConnection(connectionString string) (db *sql.DB, err error) {
	return sql.Open("postgres", connectionString)
}

// runQuery executes given SQL query and returns the execution time.
func runQuery(db *sql.DB, sequel string, record []string) (duration time.Duration, err error) {
	itfRecord := make([]interface{}, len(record))
	for i := range record {
		itfRecord[i] = record[i]
	}
	start := time.Now()
	var rows *sql.Rows
	rows, err = db.Query(sequel, itfRecord...)
	if err != nil {
		return
	}
	duration = time.Since(start)
	rows.Close()
	return
}

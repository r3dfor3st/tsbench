package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	const workersCount uint = 11
	clip := cliParamsType{workers: workersCount}
	workers := initWorkers(&clip)

	results, _ := workers.run(
		func() {
			records := [][]string{
				[]string{"hello", "world"},
				[]string{"hello", "kitty"},
				[]string{"hasta", "la", "vista", "baby"}}
			for _, rec := range records {
				rand.Seed(time.Now().UTC().UnixNano())
				workers.jobsChannels[rand.Intn(int(workersCount))] <- rec
			}
		},
		func(params <-chan []string, i int) {
			for rec := range params {
				workers.resultsChannel <- time.Duration(len(rec))
			}
		})

	compare(len(results), 3, t)
	compare(sum(results), time.Duration(8), t)
}

func sum(arr []time.Duration) (ret time.Duration) {
	for _, rec := range arr {
		ret += rec
	}
	return
}

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"io"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	clip := initCliParams()
	err := clip.assert()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	workers := initWorkers(&clip)
	results, err := workers.run(
		func() {
			err := traverseInput(&clip, &workers, func(record []string) {
				workers.jobsChannels[channelFromHost(record[0], clip.workers)] <- record
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		},
		workers.benchmarkQueryExecution)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	results.dumpStats()
}

// channelFromHost determines which worker to use for processing params based on the content
// of the first field read from each input row. This allows distributing the load among
// different workers.
func channelFromHost(host string, workers uint) uint32 {
	h := fnv.New32a()
	h.Write([]byte(host))
	return h.Sum32() % uint32(workers)
}

// traverseInput walks either the input file or STDIN containing the CSV-formatted params.
// Each read line will be parsed and handed over to the runWithin function as a string slice.
// First line in the CSV file is considered header and will always be skipped. The number of fields
// in the header defines the valid number of fields for further rows. If for example header contains
// 3 fields and then some line is parsed and contains a different number of fields, error will be
// reported.
func traverseInput(clip *cliParamsType, workers *workersType, runWithin func(record []string)) (err error) {
	var reader *csv.Reader

	if len(clip.params) > 0 {
		var f *os.File
		f, err = os.Open(clip.params)
		if err != nil {
			return
		}
		defer f.Close()
		reader = csv.NewReader(f)
	} else {
		reader = csv.NewReader(bufio.NewReader(os.Stdin))
	}

	// Streams over the input, line by line, until EOF. When a line is read, calls a runWithin callback
	// (third traverseInput function's argument).
	lines := uint(0)
	for {
		var record []string
		record, err = reader.Read()
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}
		lines++
		if lines == 1 {
			continue
		}
		runWithin(record)
	}

	return
}

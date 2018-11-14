package main

import (
	"testing"
)

func TestChannelFromHost(t *testing.T) {
	compare(channelFromHost("host_000008", 10), uint32(4), t)
	compare(channelFromHost("host_000006", 5), uint32(1), t)
	compare(channelFromHost("host_000123", 6), uint32(0), t)
}

func TestTraverseInput(t *testing.T) {
	clip := cliParamsType{
		params:  "testdata/query_params.csv",
		workers: 10}
	workers := initWorkers(&clip)

	count := 0
	traverseInput(&clip, &workers, func(rec []string) {
		if count == 0 {
			compare(rec[0], "host_000008", t)
			compare(rec[1], "2017-01-01 08:59:22", t)
			compare(rec[2], "2017-01-01 09:59:22", t)
		}
		count++
	})
	compare(count, 200, t)
}

func compare(actual interface{}, expected interface{}, t *testing.T) {
	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

package main

import (
	"testing"
	"time"
)

func Example_measurementsType_DumpStats() {
	measurements := measurementsType{123, 456, 789}
	measurements.dumpStats()
	// Output:
	// total db queries: 3
	// total time: 1.368µs
	// min time: 123ns
	// max time: 789ns
	// avg time: 456ns
	// median time: 456ns
	// standard deviation: 271ns
}

func TestCount(t *testing.T) {
	measurements := measurementsType{123, 456, 789}
	compare(measurements.count(), 3, t)
}

func TestSum(t *testing.T) {
	measurements := measurementsType{123, 456, 789}
	compare(measurements.sum(), time.Duration(1368)*time.Nanosecond, t)
}

func TestMin(t *testing.T) {
	measurements := measurementsType{123, 456, 789}
	compare(measurements.min(), time.Duration(123)*time.Nanosecond, t)
}

func TestMax(t *testing.T) {
	measurements := measurementsType{123, 456, 789}
	compare(measurements.max(), time.Duration(789)*time.Nanosecond, t)
}

func TestAvg(t *testing.T) {
	measurements := measurementsType{123, 456, 789}
	compare(measurements.avg(), time.Duration(456)*time.Nanosecond, t)
}

func TestMedian(t *testing.T) {
	measurements := measurementsType{5, 3, 2, 4, 1, 6, 8}
	compare(measurements.median(), time.Duration(4)*time.Nanosecond, t)
}

func TestSdev(t *testing.T) {
	measurements := measurementsType{123, 456, 789}
	compare(measurements.sdev(), time.Duration(271)*time.Nanosecond, t)
}

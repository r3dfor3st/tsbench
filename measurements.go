package main

import (
	"fmt"
	"math"
	"time"
)

// measurementsType is a container meant to hold execution durations from a single
// tool run. It also contains methods to calculate different statistics.
type measurementsType []time.Duration

// dumpStats prints different stats to STDOUT.
func (m measurementsType) dumpStats() {
	if len(m) == 0 {
		return
	}

	fmt.Println("total db queries:", m.count())
	fmt.Println("total time:", m.sum())
	fmt.Println("min time:", m.min())
	fmt.Println("max time:", m.max())
	fmt.Println("avg time:", m.avg())
	fmt.Println("median time:", m.median())
	fmt.Println("standard deviation:", m.sdev())
}

// count returns total number of queries run.
func (m measurementsType) count() int {
	return len(m)
}

// sum returns total duration of all queries across all workers.
func (m measurementsType) sum() time.Duration {
	var sum time.Duration
	for _, dur := range m {
		sum += dur
	}
	return sum
}

// min returns minimum execution time.
func (m measurementsType) min() time.Duration {
	if m.count() == 0 {
		return 0.0
	}

	min := m[0]
	for _, dur := range m[1:] {
		if dur < min {
			min = dur
		}
	}
	return min
}

// max returns maximum execution time.
func (m measurementsType) max() time.Duration {
	if m.count() == 0 {
		return 0.0
	}

	max := m[0]
	for _, dur := range m[1:] {
		if dur > max {
			max = dur
		}
	}
	return max
}

// avg returns average execution time.
func (m measurementsType) avg() time.Duration {
	if m.count() == 0 {
		return 0.0
	}

	return m.sum() / time.Duration(m.count())
}

// median returns median execution time.
func (m measurementsType) median() (result time.Duration) {
	if m.count() == 0 {
		return 0.0
	}

	n := len(m) / 2
	if len(m)%2 == 0 {
		result = (m.nth(n-1) + m.nth(n)) / 2.0
	} else {
		result = m.nth(n)
	}
	return
}

// sdev returns standard deviation of execution time.
func (m measurementsType) sdev() (result time.Duration) {
	if m.count() == 0 {
		return 0.0
	}

	mean := m.avg()
	sum := 0.0
	for _, dur := range m {
		sum += math.Pow(float64((dur - mean).Nanoseconds()), 2)
	}
	result = time.Duration(math.Sqrt(sum / float64(len(m))))
	return
}

func (m measurementsType) nth(n int) (result time.Duration) {
	var pivot time.Duration
	var underPivot, overPivot, eqPivot measurementsType

	pivot = m[len(m)/2]

	for _, dur := range m {
		if dur < pivot {
			underPivot = append(underPivot, dur)
		} else if dur > pivot {
			overPivot = append(overPivot, dur)
		} else {
			eqPivot = append(eqPivot, dur)
		}
	}

	if n < len(underPivot) {
		result = underPivot.nth(n)
	} else if n < len(underPivot)+len(eqPivot) {
		result = pivot
	} else {
		result = overPivot.nth(n - len(underPivot) - len(eqPivot))
	}

	return
}

package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	JobStart = time.Now()

	totalJobs = uint64(0)
	sum       = uint64(0)
	priorSecs = make([]uint64, 0)
	calcsDone = uint64(0)
)

func tickerWorker() {
	t := time.NewTicker(AverageFrequency)

	for {
		select {
		case <-globalDone:
			t.Stop()
			return
		case <-t.C:
			calcCount := atomic.SwapUint64(&calcsDone, 0)
			priorSecs = append(priorSecs, calcCount)
			sum += calcCount
			totalJobs += calcCount

			if len(priorSecs) > AveragePoints {
				sum -= priorSecs[0]
				priorSecs = priorSecs[1:]
			}

			fmt.Printf("Averaging %d IDm tests per %s (processed %d total jobs over %s) on this machine\n", sum/uint64(len(priorSecs)), AverageFrequency, totalJobs, time.Now().Sub(JobStart))
		}
	}
}

package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	JobStart = time.Now()

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

			if len(priorSecs) > AveragePoints {
				sum -= priorSecs[0]
				priorSecs = priorSecs[1:]
			}

			fmt.Printf("Averaging %d IDm tests per %s\n", sum/uint64(len(priorSecs)), AverageFrequency)
		}
	}
}

package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

type matchMode uint8

const (
	contains = 0
	prefix   = 1
	align    = 2
)

func match(in string) bool {
	switch CurrentMatchMode {
	case contains:
		return strings.Contains(in, CurrentMatchTarget)
	case prefix:
		return strings.HasPrefix(in, CurrentMatchTarget)
	case align:
		for len(in) > 0 {
			if strings.HasPrefix(in, CurrentMatchTarget) {
				return true
			}

			in = in[4:]
		}

		return false
	}

	panic("invalid mode")
}

func largestMatch(in string) int { // implement as a sliding window
	start := 0
	max := 0

	for start < len(in) {
		count := 1

		for start+count < len(in) && strings.HasPrefix(CurrentMatchTarget, in[start:start+count]) {
			if CurrentMatchMode == align && start%4 != 0 {
				break // stop this search
			}

			if count > max {
				max = count
			}

			count++
		}

		start++
		if CurrentMatchMode == prefix {
			return max
		}
	}

	return max
}

type maxSubmit struct {
	count uint64
	value string
	idm   string
}

var (
	cMax  = uint64(0)
	maxCh = make(chan maxSubmit, runtime.NumCPU())
)

func maxWorker() {
	for {
		select {
		case <-globalDone:
			return
		case result := <-maxCh:
			if result.count <= cMax {
				continue
			}

			cMax = result.count
			fmt.Printf("new max result (%d matched): %s (idm %s) (took %s)\n", cMax, result.value, result.idm, time.Now().Sub(JobStart))
		}
	}
}

func maxMatch(in string, idm string) {
	n := uint64(largestMatch(in))
	if n > atomic.LoadUint64(&cMax) {
		maxCh <- maxSubmit{
			count: n,
			value: in,
			idm:   idm,
		}
	}
}

package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	globalDone     = make(chan any)
	globalDoneOnce = &sync.Once{}
)

func stringifyBytes(in []byte) string {
	out := ""

	for _, v := range in {
		out += fmt.Sprintf("%02X ", v)
	}

	return strings.TrimSpace(out)
}

func prepareInput(in string) string {
	return strings.TrimSpace(strings.ReplaceAll(in, " ", ""))
}

func done() {
	globalDoneOnce.Do(func() {
		close(globalDone)
	})
}

func main() {
	if len(os.Args) < 2 {
		panic("need a search term")
	}
	for _, v := range CurrentMatchTarget {
		if strings.Index(CardAlphabet, string(v)) == -1 {
			panic("search term contains invalid characters (no IOQV)")
		}
	}

	go maxWorker()
	go tickerWorker()

	seen := sync.Map{}
	results := sync.Map{}
	max, _ := strconv.ParseUint(prepareInput("0FFF FFFF FFFF FFFF"), 16, 64)

	wg := &sync.WaitGroup{}
	wg.Add(ThreadCount)

	for threadId := range ThreadCount {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-globalDone:
					return
				default:
				}

				// first, find a random ID
				id := rand.N(max)
				idStr := fmt.Sprintf("%016X", id)

				if !AllowConflicts {
					// push it to the map
					_, loaded := seen.LoadOrStore(idStr, true)
					if loaded {
						continue // find another
					}
				}

				// calculate
				res, err := UIDToKonami(idStr, NewEncrypter())
				if err != nil {
					fmt.Println(threadId, "failed", idStr, err)
					continue
				}

				if match(res) {
					fmt.Println("found", idStr, res)
					results.Store(idStr, res)
					done()
					return
				}

				maxMatch(res, idStr)
				atomic.AddUint64(&calcsDone, 1)
			}
		}()
	}

	wg.Wait()

	fmt.Println("found results: ")

	results.Range(func(key, res any) bool {
		konmai, _ := UIDToKonami(key.(string), NewEncrypter())
		fmt.Println(key, "konami ID:", konmai)
		return true
	})
}

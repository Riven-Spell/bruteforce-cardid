package main

import (
	"fmt"
	"os"
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

	results := sync.Map{}

	wg := &sync.WaitGroup{}
	wg.Add(ThreadCount)

	for threadId := range uint(ThreadCount) {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-globalDone:
					return
				default:
				}

				// first, find a random ID
				idStr := DivideLabour.GetJob(threadId)

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

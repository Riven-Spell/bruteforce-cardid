package main

import (
	"os"
	"runtime"
	"time"
)

// CurrentMatchMode has a few options,
// contains - exists at all in the 16 chars
// align - exists on a 4 char border
// prefix - exists from 0
var CurrentMatchMode matchMode = contains

// CurrentMatchTarget is filled from command line typically, but a string suffices.
var CurrentMatchTarget = os.Args[1]

// AllowConflicts -- potentially slower, lower memory usage.
var AllowConflicts = true

// ThreadCount -- determines the number of goroutines hammering away.
var ThreadCount = runtime.NumCPU()

// AveragePoints -- how many times does an average calc count get included
var AveragePoints = 10

// AverageFrequency how frequently are averages taken
var AverageFrequency = time.Second

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
var CurrentMatchMode matchMode = align

// CurrentMatchTarget is filled from command line typically, but a string suffices.
var CurrentMatchTarget = os.Args[1]

// DivideLabour -- define a labour division strategy, checkk division_of_labour.go
var DivideLabour LabourDivider = NewLabourDividerStrict(0, MachineDefinition{32, 12}, MachineDefinition{24, 6})

// ThreadCount -- determines the number of goroutines hammering away.
var ThreadCount = runtime.NumCPU() * 3

// AveragePoints -- how many times does an average calc count get included
var AveragePoints = 10

// AverageFrequency how frequently are averages taken
var AverageFrequency = time.Second * 15

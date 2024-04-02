package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/branchpred"
	"github.com/nsengupta5/Branch-Predictor/instruction"
)

const (
	defaultHistoryLength = 8
	defaultTableSize     = 256
)

// This file contains the main function to run the branch predictor
// It reads in the command line arguements and initializes the branch predictor

func main() {
	traceFile := os.Args[1]
	algorithm := os.Args[2]

	instructions, err := instruction.ReadTraceFile(traceFile)
	if err != nil {
		panic(err)
	}

	var tableSize uint64 = getTableSize(algorithm)

	simulator := branchpred.NewBranchPredictor(algorithm, tableSize, defaultHistoryLength)

	result := simulator.Predict(instructions)
	fmt.Println(result)
}

func getTableSize(algorithm string) uint64 {
	if algorithm != "two-bit" {
		return 256
	}

	var tableSize uint64
	if len(os.Args) < 4 {
		panic("Please provide table size for two-bit predictor")
	}
	tableSize, err := strconv.ParseUint(os.Args[3], 10, 64)
	if err != nil {
		panic(err)
	}

	switch tableSize {
	case 512, 1024, 2048, 4096:
	default:
		panic("Invalid table size. Please provide a valid table size")
	}
	return tableSize
}

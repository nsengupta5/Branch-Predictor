package main

import (
	"fmt"
	"os"

	"github.com/nsengupta5/Branch-Predictor/branchpred"
	"github.com/nsengupta5/Branch-Predictor/instruction"
)

// This file contains the main function to run the branch predictor
// It reads in the command line arguements and initializes the branch predictor

func main() {
	traceFile := os.Args[1]
	algorithm := os.Args[2]

	simulator := branchpred.NewBranchPredictor(algorithm)
	instructions, err := instruction.ReadTraceFile(traceFile)
	if err != nil {
		panic(err)
	}

	result := simulator.Predict(instructions)
	fmt.Printf("Accuracy: %.2f%%\n", result)
}

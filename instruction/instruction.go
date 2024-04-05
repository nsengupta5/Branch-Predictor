package instruction

// This file contains the Instruction struct and functions to read a trace file.
// The main purpose of the file is to read from a trace file and obtain a scice
// of instructions that can be used to predict branches.

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Instruction represents a single instruction in a trace file
type Instruction struct {
	PCAddress     string
	TargetAddress string
	BranchKind    rune
	Direct        bool
	Conditional   bool
	Taken         bool
}

// ReadTraceFile reads a trace file and returns a slice of instructions
// It takes in an maxLines argument to limit the number of lines read
// from the file. If maxLines is negative, all lines are read.
func ReadTraceFile(traceFile string, maxLines int) ([]Instruction, error) {
	file, err := os.Open(traceFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var instructions []Instruction
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++

		// If maxLines is set and we've read the max number of lines, break
		if maxLines >= 0 && lineCount > int(maxLines) {
			break
		}

		line := scanner.Text()
		fields := strings.Fields(line)

		direct, _ := strconv.ParseBool(fields[3])
		conditional, _ := strconv.ParseBool(fields[4])
		taken, _ := strconv.ParseBool(fields[5])

		instruction := Instruction{
			PCAddress:     fields[0],
			TargetAddress: fields[1],
			BranchKind:    rune(fields[2][0]),
			Direct:        direct,
			Conditional:   conditional,
			Taken:         taken,
		}

		instructions = append(instructions, instruction)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return instructions, nil
}

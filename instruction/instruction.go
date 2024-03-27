package instruction

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	PCAddress     string
	TargetAddress string
	BranchKind    rune
	Direct        bool
	Conditional   bool
	Taken         bool
}

func ReadTraceFile(traceFile string) ([]Instruction, error) {
	file, err := os.Open(traceFile)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var instructions []Instruction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
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

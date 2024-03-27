package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type AlwaysTaken struct{}

func NewAlwaysTaken() *AlwaysTaken {
	return &AlwaysTaken{}
}

func (at *AlwaysTaken) Predict(instructions []instruction.Instruction) float64 {
	totalBranches := 0
	mispredictions := 0

	for _, instruction := range instructions {
		if instruction.Conditional {
			totalBranches++
			if !instruction.Taken {
				mispredictions++
			}
		}
	}

	prediction := Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction.Accuracy()
}

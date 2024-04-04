package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type AlwaysTaken struct{}

func NewAlwaysTaken() *AlwaysTaken {
	return &AlwaysTaken{}
}

func (at *AlwaysTaken) Predict(instructions []instruction.Instruction) utils.Prediction {
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

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

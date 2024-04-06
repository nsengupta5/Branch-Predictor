package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type Profiled struct {
	// Add fields here
}

// NewProfiled creates a new instance of the Profiled struct
func NewProfiled(condif *utils.ProfiledConfig, metadata *utils.MetaData) *Profiled {
	// Add code here

	return &Profiled{
		// Add fields here
	}
}

// Predict predicts the outcome of the branches in the given instructions
func (p *Profiled) Predict(instructions []instruction.Instruction) utils.Prediction {
	totalBranches := 0
	mispredictions := 0

	for _, instruction := range instructions {
		// Check if the instruction is a branch instruction
		if instruction.Conditional {

		}
	}

	return utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}
}

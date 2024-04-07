package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

// AlwaysTaken is a struct that implements the Algorithm interface
type AlwaysTaken struct {
	config   *utils.AlwaysTakenConfig
	metadata *utils.MetaData
}

// NewAlwaysTaken creates a new instance of the AlwaysTaken struct
func NewAlwaysTaken(config *utils.AlwaysTakenConfig, metadata *utils.MetaData) *AlwaysTaken {
	return &AlwaysTaken{
		config:   config,
		metadata: metadata,
	}
}

// Predict predicts the outcome of the branches in the given instructions
func (at *AlwaysTaken) Predict(instructions []instruction.Instruction) utils.Prediction {
	totalBranches := 0
	mispredictions := 0

	for _, instruction := range instructions {
		// Check if the instruction is a branch instruction
		if instruction.Conditional {
			totalBranches++

			// Checks if the branch instruction is mispredicted
			isMispredicted := false

			// Misprediction occurs if the branch is not taken
			if !instruction.Taken {
				mispredictions++
				isMispredicted = true
			}

			// Update the metadata
			// The AlwaysTaken algorithm does not use any states,
			// so the state is set to 4, (not a valid state for other algorithms)
			at.metadata.Update(false, false, 0, isMispredicted, 4)
		}
	}

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

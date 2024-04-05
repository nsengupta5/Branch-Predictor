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
			at.UpdateMetaData(instruction, isMispredicted)
		}
	}

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

// Updates the metadata for the given instruction, based on whether it was mispredicted or not
func (at *AlwaysTaken) UpdateMetaData(instruction instruction.Instruction, isMispredicted bool) {
	// If the metadata already exists for the given instruction, update it
	if at.metadata.Exists(instruction.PCAddress) {
		at.metadata.Update(instruction, isMispredicted)
	} else {
		at.metadata.AddBranch(instruction.PCAddress)
		at.metadata.Update(instruction, isMispredicted)
	}
}

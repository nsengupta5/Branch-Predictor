package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type AlwaysTaken struct {
	config   *utils.AlwaysTakenConfig
	metadata *utils.MetaData
}

func NewAlwaysTaken(config *utils.AlwaysTakenConfig, metadata *utils.MetaData) *AlwaysTaken {
	return &AlwaysTaken{
		config:   config,
		metadata: metadata,
	}
}

func (at *AlwaysTaken) Predict(instructions []instruction.Instruction) utils.Prediction {
	totalBranches := 0
	mispredictions := 0

	for _, instruction := range instructions {
		if instruction.Conditional {
			totalBranches++
			isMispredicted := false

			if !instruction.Taken {
				mispredictions++
				isMispredicted = true
			}

			at.metadata.Update(instruction, isMispredicted)
		}
	}

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

func (at *AlwaysTaken) GetName() string {
	return at.config.Name
}

func (at *AlwaysTaken) UpdateMetaData(instruction instruction.Instruction, isMispredicted bool) {
	if at.metadata.Exists(instruction.PCAddress) {
		at.metadata.Update(instruction, isMispredicted)
	} else {
		at.metadata.AddBranch(instruction.PCAddress)
		at.metadata.Update(instruction, isMispredicted)
	}
}

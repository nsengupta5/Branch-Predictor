package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

// Algorithm is an interface that defines the methods that a branch predictor algorithm should implement
type Algorithm interface {
	Predict(il []instruction.Instruction) utils.Prediction
	UpdateMetaData(instruction instruction.Instruction, isMispredicted bool)
}

// BranchPredictor is a struct that contains the branch prediction algorithm and its metadata
type BranchPredictor struct {
	Algorithm Algorithm
	Metadata  *utils.MetaData
}

// NewBranchPredictor creates a new instance of the BranchPredictor struct
func NewBranchPredictor(config utils.Config) *BranchPredictor {
	switch cfg := config.(type) {
	case utils.AlwaysTakenConfig:
		metadata := utils.NewMetaData(0)
		algo := NewAlwaysTaken(&cfg, metadata)
		return &BranchPredictor{Algorithm: algo, Metadata: metadata}

	case utils.TwoBitConfig:
		var tableSize uint64 = cfg.TableSize
		metadata := utils.NewMetaData(tableSize)
		algo := NewTwoBit(&cfg, metadata)
		return &BranchPredictor{Algorithm: algo, Metadata: metadata}

	case utils.GShareConfig:
		var tableSize uint64 = cfg.TableSize
		metadata := utils.NewMetaData(tableSize)
		algo := NewGshare(&cfg, metadata)
		return &BranchPredictor{Algorithm: algo, Metadata: metadata}

	default:
		return nil
	}
}

// Predict predicts the outcome of the predictions made using the algorithm
func (bp *BranchPredictor) Predict(instructions []instruction.Instruction) utils.Prediction {
	return bp.Algorithm.Predict(instructions)
}

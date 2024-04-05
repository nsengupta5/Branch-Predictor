package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type Config interface{}

type Algorithm interface {
	Predict(il []instruction.Instruction) utils.Prediction
	UpdateMetaData(instruction instruction.Instruction, isMispredicted bool)
}

type BranchPredictor struct {
	Algorithm Algorithm
	Metadata  *utils.MetaData
}

func NewBranchPredictor(config Config) *BranchPredictor {
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

func (bp *BranchPredictor) Predict(instructions []instruction.Instruction) utils.Prediction {
	return bp.Algorithm.Predict(instructions)
}

package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type Algorithm interface {
	Predict(il []instruction.Instruction) Prediction
}

type BranchPredictor struct {
	Algorithm Algorithm
}

func NewBranchPredictor(algorithm string, tableSize uint64) *BranchPredictor {
	switch algorithm {
	case "always-taken":
		return &BranchPredictor{Algorithm: NewAlwaysTaken()}
	case "two-bit":
		return &BranchPredictor{Algorithm: NewTwoBit(tableSize)}
	case "gshare":
		return &BranchPredictor{Algorithm: NewGshare(tableSize)}
	default:
		return nil
	}
}

func (bp *BranchPredictor) Predict(instructions []instruction.Instruction) Prediction {
	return bp.Algorithm.Predict(instructions)
}

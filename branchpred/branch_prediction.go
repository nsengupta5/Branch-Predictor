package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type Algorithm interface {
	Predict(il []instruction.Instruction) float64
}

type BranchPredictor struct {
	Algorithm Algorithm
}

func NewBranchPredictor(algorithm string) *BranchPredictor {
	switch algorithm {
	case "always-taken":
		return &BranchPredictor{Algorithm: NewAlwaysTaken()}
	default:
		return nil
	}
}

func (bp *BranchPredictor) Predict(instructions []instruction.Instruction) float64 {
	return bp.Algorithm.Predict(instructions)
}

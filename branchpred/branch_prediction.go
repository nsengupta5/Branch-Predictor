package branchpred

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type Algorithm interface {
	Predict(il []instruction.Instruction) utils.Prediction
	GetName() string
	UpdateMetaData(instruction instruction.Instruction, isMispredicted bool)
}

type BranchPredictor struct {
	Algorithm Algorithm
	Metadata  *utils.MetaData
}

func NewBranchPredictor(algorithm string, tableSize uint64) *BranchPredictor {
	metaData := utils.NewMetaData(tableSize)

	switch algorithm {
	case "always-taken":
		return &BranchPredictor{Algorithm: NewAlwaysTaken(metaData), Metadata: metaData}
	case "two-bit":
		return &BranchPredictor{Algorithm: NewTwoBit(metaData, tableSize), Metadata: metaData}
	case "gshare":
		return &BranchPredictor{Algorithm: NewGshare(metaData, tableSize), Metadata: metaData}
	default:
		return nil
	}
}

func (bp *BranchPredictor) Predict(instructions []instruction.Instruction) utils.Prediction {
	return bp.Algorithm.Predict(instructions)
}

func (bp *BranchPredictor) ExportMetaData() utils.MetaData {
	metaData, err := json.Marshal(bp.Metadata)
	if err != nil {
		panic(err)
	}

	var filepath string = fmt.Sprintf("metadata/%s.json", bp.Algorithm.GetName())
	err = os.WriteFile(filepath, metaData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Metadata exported to", filepath)
	return *bp.Metadata
}

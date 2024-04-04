package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type TwoBit struct {
	name            string
	predictionTable map[uint64]utils.State
	indexBitCount   uint64
}

func NewTwoBit(tableSize uint64) *TwoBit {
	indexBitCount := uint64(math.Log2(float64(tableSize)))

	return &TwoBit{
		name:            "two-bit",
		predictionTable: make(map[uint64]utils.State, tableSize),
		indexBitCount:   indexBitCount,
	}
}

func (tb *TwoBit) Predict(instructions []instruction.Instruction) utils.Prediction {
	totalBranches := 0
	mispredictions := 0

	mask := uint64((1 << tb.indexBitCount) - 1)

	for _, instruction := range instructions {
		if instruction.Conditional {
			totalBranches++
			pcAddress, err := strconv.ParseUint(instruction.PCAddress, 16, 64)
			if err != nil {
				panic(err)
			}
			pcAddress = pcAddress & mask
			taken := instruction.Taken

			if _, ok := tb.predictionTable[pcAddress]; !ok {
				tb.predictionTable[pcAddress] = utils.StronglyNotTaken
			}

			switch tb.predictionTable[pcAddress] {
			case utils.StronglyNotTaken:
				if taken {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.WeaklyNotTaken
				}
			case utils.WeaklyNotTaken:
				if taken {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.StronglyTaken
				} else {
					tb.predictionTable[pcAddress] = utils.StronglyNotTaken
				}
			case utils.WeaklyTaken:
				if taken {
					tb.predictionTable[pcAddress] = utils.StronglyTaken
				} else {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.StronglyNotTaken
				}
			case utils.StronglyTaken:
				if !taken {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.WeaklyTaken
				}
			}
		}
	}

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

func (tb *TwoBit) GetName() string {
	return tb.name
}

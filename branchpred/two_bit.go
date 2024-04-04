package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type TwoBit struct {
	predictionTable map[uint64]State
	keys            []uint64
	indexBitCount   uint64
}

func NewTwoBit(tableSize uint64) *TwoBit {
	indexBitCount := uint64(math.Log2(float64(tableSize)))

	return &TwoBit{
		predictionTable: make(map[uint64]State, tableSize),
		keys:            make([]uint64, 0, tableSize),
		indexBitCount:   indexBitCount,
	}
}

func (tb *TwoBit) Predict(instructions []instruction.Instruction) Prediction {
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
				tb.predictionTable[pcAddress] = StronglyNotTaken
			}

			switch tb.predictionTable[pcAddress] {
			case StronglyNotTaken:
				if taken {
					mispredictions++
					tb.predictionTable[pcAddress] = WeaklyNotTaken
				}
			case WeaklyNotTaken:
				if taken {
					mispredictions++
					tb.predictionTable[pcAddress] = StronglyTaken
				} else {
					tb.predictionTable[pcAddress] = StronglyNotTaken
				}
			case WeaklyTaken:
				if taken {
					tb.predictionTable[pcAddress] = StronglyTaken
				} else {
					mispredictions++
					tb.predictionTable[pcAddress] = StronglyNotTaken
				}
			case StronglyTaken:
				if !taken {
					mispredictions++
					tb.predictionTable[pcAddress] = WeaklyTaken
				}
			}
		}
	}

	prediction := Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

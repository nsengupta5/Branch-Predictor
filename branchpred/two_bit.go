package branchpred

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type State uint8

const (
	StronglyNotTaken State = iota // 00
	WeaklyNotTaken                // 01
	WeaklyTaken                   // 10
	StronglyTaken                 // 11
)

type TwoBit struct {
	predictionTable map[string]State
	twoBitTableSize uint64
	keys            []string
}

func NewTwoBit(tableSize uint64) *TwoBit {
	return &TwoBit{
		predictionTable: make(map[string]State, tableSize),
		twoBitTableSize: tableSize,
		keys:            make([]string, 0, tableSize),
	}
}

func (tb *TwoBit) Predict(instructions []instruction.Instruction) float64 {
	totalBranches := 0
	mispredictions := 0

	for _, instruction := range instructions {
		if instruction.Conditional {
			totalBranches++
			pcAddress := instruction.PCAddress
			taken := instruction.Taken

			if _, ok := tb.predictionTable[pcAddress]; !ok {
				if uint64(len(tb.predictionTable)) == tb.twoBitTableSize {
					key := tb.keys[0]
					delete(tb.predictionTable, key)
					tb.keys = tb.keys[1:]
				}
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
					tb.predictionTable[pcAddress] = WeaklyNotTaken
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

	return prediction.Accuracy()
}

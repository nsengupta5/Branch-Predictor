package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type Gshare struct {
	globalHistoryRegister uint64
	patternHistoryTable   map[uint64]State
	historyLength         uint64
}

func NewGshare(tableSize uint64) *Gshare {
	historyLength := uint64(math.Log2(float64(tableSize)))

	return &Gshare{
		globalHistoryRegister: 0,
		patternHistoryTable:   make(map[uint64]State, tableSize),
		historyLength:         historyLength,
	}
}

func (gs *Gshare) Predict(instructions []instruction.Instruction) Prediction {
	totalBranches := 0
	mispredictions := 0

	mask := uint64((1 << gs.historyLength) - 1)

	for _, instruction := range instructions {
		if instruction.Conditional {
			totalBranches++
			pcAddress, err := strconv.ParseUint(instruction.PCAddress, 16, 64)
			if err != nil {
				panic(err)
			}

			pcAddress = pcAddress & mask
			taken := instruction.Taken

			gs.updateGlobalHistoryRegister(taken)

			index := gs.getIndex(pcAddress)
			if _, ok := gs.patternHistoryTable[index]; !ok {
				gs.patternHistoryTable[index] = StronglyNotTaken
			}

			switch gs.patternHistoryTable[index] {
			case StronglyNotTaken:
				if taken {
					mispredictions++
					gs.patternHistoryTable[index] = WeaklyNotTaken
				}
			case WeaklyNotTaken:
				if taken {
					mispredictions++
					gs.patternHistoryTable[index] = WeaklyTaken
				} else {
					gs.patternHistoryTable[index] = StronglyNotTaken
				}
			case WeaklyTaken:
				if !taken {
					mispredictions++
					gs.patternHistoryTable[index] = WeaklyNotTaken
				} else {
					gs.patternHistoryTable[index] = StronglyTaken
				}
			case StronglyTaken:
				if !taken {
					mispredictions++
					gs.patternHistoryTable[index] = WeaklyTaken
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

func (gs *Gshare) getIndex(pcAddress uint64) uint64 {
	var addressBits uint64 = uint64(pcAddress) & ((1 << gs.historyLength) - 1)
	return addressBits ^ gs.globalHistoryRegister
}

func (gs *Gshare) updateGlobalHistoryRegister(taken bool) {
	gs.globalHistoryRegister <<= 1
	if taken {
		gs.globalHistoryRegister |= 1
	}
	gs.globalHistoryRegister &= (1 << gs.historyLength) - 1
}

package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type Gshare struct {
	name                  string
	metadata              *utils.MetaData
	globalHistoryRegister uint64
	patternHistoryTable   map[uint64]utils.State
	historyLength         uint64
}

func NewGshare(metadata *utils.MetaData, tableSize uint64) *Gshare {
	historyLength := uint64(math.Log2(float64(tableSize)))

	return &Gshare{
		name:                  "gshare",
		globalHistoryRegister: 0,
		patternHistoryTable:   make(map[uint64]utils.State, tableSize),
		historyLength:         historyLength,
		metadata:              metadata,
	}
}

func (gs *Gshare) Predict(instructions []instruction.Instruction) utils.Prediction {
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
			isMispredicted := false

			gs.updateGlobalHistoryRegister(taken)

			index := gs.getIndex(pcAddress)
			if _, ok := gs.patternHistoryTable[index]; !ok {
				gs.patternHistoryTable[index] = utils.StronglyNotTaken
			}

			switch gs.patternHistoryTable[index] {
			case utils.StronglyNotTaken:
				if taken {
					mispredictions++
					gs.patternHistoryTable[index] = utils.WeaklyNotTaken
					isMispredicted = true
				}
			case utils.WeaklyNotTaken:
				if taken {
					mispredictions++
					gs.patternHistoryTable[index] = utils.WeaklyTaken
					isMispredicted = true
				} else {
					gs.patternHistoryTable[index] = utils.StronglyNotTaken
				}
			case utils.WeaklyTaken:
				if !taken {
					mispredictions++
					gs.patternHistoryTable[index] = utils.WeaklyNotTaken
					isMispredicted = true
				} else {
					gs.patternHistoryTable[index] = utils.StronglyTaken
				}
			case utils.StronglyTaken:
				if !taken {
					mispredictions++
					gs.patternHistoryTable[index] = utils.WeaklyTaken
					isMispredicted = true
				}
			}
			gs.UpdateMetaData(instruction, isMispredicted)
		}
	}

	prediction := utils.Prediction{
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

func (gs *Gshare) GetName() string {
	return gs.name
}

func (gs *Gshare) UpdateMetaData(instruction instruction.Instruction, isMispredicted bool) {
	if gs.metadata.Exists(instruction.PCAddress) {
		gs.metadata.Update(instruction, isMispredicted)
	} else {
		gs.metadata.AddBranch(instruction.PCAddress)
		gs.metadata.Update(instruction, isMispredicted)
	}
}

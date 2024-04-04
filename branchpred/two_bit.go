package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type TwoBit struct {
	config          *utils.TwoBitConfig
	metadata        *utils.MetaData
	indexBitCount   uint64
	predictionTable map[uint64]utils.State
	initialState    utils.State
}

func NewTwoBit(config *utils.TwoBitConfig, metadata *utils.MetaData) *TwoBit {
	tableSize := config.TableSize
	indexBitCount := uint64(math.Log2(float64(tableSize)))
	initialState, ok := utils.StateMap[config.InitialState]
	if !ok {
		panic("Invalid initial state")
	}

	return &TwoBit{
		config:          config,
		metadata:        metadata,
		indexBitCount:   indexBitCount,
		predictionTable: make(map[uint64]utils.State, tableSize),
		initialState:    initialState,
	}
}

func (tb *TwoBit) Predict(instructions []instruction.Instruction) utils.Prediction {
	totalBranches := 0
	mispredictions := 0

	mask := uint64((1 << tb.indexBitCount) - 1)

	for _, instruction := range instructions {
		if instruction.Conditional {
			totalBranches++
			pcAddressInt, err := strconv.ParseUint(instruction.PCAddress, 16, 64)
			if err != nil {
				panic(err)
			}
			pcAddress := pcAddressInt & mask
			taken := instruction.Taken
			isMispredicted := false

			if _, ok := tb.predictionTable[pcAddress]; !ok {
				tb.predictionTable[pcAddress] = tb.initialState
			}

			switch tb.predictionTable[pcAddress] {
			case utils.StronglyNotTaken:
				if taken {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.WeaklyNotTaken
					isMispredicted = true
				}
			case utils.WeaklyNotTaken:
				if taken {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.StronglyTaken
					isMispredicted = true
				} else {
					tb.predictionTable[pcAddress] = utils.StronglyNotTaken
				}
			case utils.WeaklyTaken:
				if taken {
					tb.predictionTable[pcAddress] = utils.StronglyTaken
				} else {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.StronglyNotTaken
					isMispredicted = true
				}
			case utils.StronglyTaken:
				if !taken {
					mispredictions++
					tb.predictionTable[pcAddress] = utils.WeaklyTaken
					isMispredicted = true
				}
			}

			tb.UpdateMetaData(instruction, isMispredicted)
		}
	}

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

func (tb *TwoBit) GetName() string {
	return tb.config.Name
}

func (tb *TwoBit) UpdateMetaData(instruction instruction.Instruction, isMispredicted bool) {
	if tb.metadata.Exists(instruction.PCAddress) {
		tb.metadata.Update(instruction, isMispredicted)
	} else {
		tb.metadata.AddBranch(instruction.PCAddress)
		tb.metadata.Update(instruction, isMispredicted)
	}
}

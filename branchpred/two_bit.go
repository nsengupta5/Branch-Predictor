package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

// TwoBit is a struct that implements the Algorithm interface
type TwoBit struct {
	config          *utils.TwoBitConfig
	metadata        *utils.MetaData
	indexBitCount   uint64
	predictionTable map[uint64]utils.State
	initialState    utils.State
}

// NewTwoBit creates a new instance of the TwoBit struct. It initializes the prediction table
// that is used to store the state of each branch instruction. The indexBitCount used as the
// key for the prediction table is calculated based on the table size provided in the configuration.
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

// Predict predicts the outcome of the branches in the given instructions
func (tb *TwoBit) Predict(instructions []instruction.Instruction) utils.Prediction {
	totalBranches := 0
	mispredictions := 0

	// Create a mask to extract the index bits from the lower bits of the PC address
	mask := uint64((1 << tb.indexBitCount) - 1)

	for _, instruction := range instructions {
		// Check if the instruction is a branch instruction
		if instruction.Conditional {
			totalBranches++
			pcAddressInt, err := strconv.ParseUint(instruction.PCAddress, 16, 64)
			if err != nil {
				panic(err)
			}

			// Extract the lower index bits from the PC address
			pcAddress := pcAddressInt & mask
			taken := instruction.Taken
			isMispredicted := false

			// Add the PC address to the prediction table if it does not exist
			if _, ok := tb.predictionTable[pcAddress]; !ok {
				tb.predictionTable[pcAddress] = tb.initialState
			}

			// Update the state of the branch instruction based on the prediction table
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

			// Update the metadata
			tb.UpdateMetaData(instruction, isMispredicted)
		}
	}

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

// UpdateMetaData updates the metadata for the given instruction, based on whether it was mispredicted or not
func (tb *TwoBit) UpdateMetaData(instruction instruction.Instruction, isMispredicted bool) {
	if tb.metadata.Exists(instruction.PCAddress) {
		tb.metadata.Update(instruction, isMispredicted)
	} else {
		tb.metadata.AddBranch(instruction.PCAddress)
		tb.metadata.Update(instruction, isMispredicted)
	}
}

package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

// Gshare is a struct that implements the Algorithm interface
type Gshare struct {
	config                *utils.GShareConfig
	metadata              *utils.MetaData
	globalHistoryRegister uint64
	patternHistoryTable   map[uint64]utils.State
	historyLength         uint64
	initialState          utils.State
}

// NewGshare creates a new instance of the Gshare struct. It initializes the pattern history table
// that is used to store the state of each branch instruction. The history length used as the
// key for the pattern history table is calculated based on the table size provided in the configuration.
// Finally, it initializes the global history register to 0, which is the number of bits used to store
// the global history of the branch instructions.
func NewGshare(config *utils.GShareConfig, metadata *utils.MetaData) *Gshare {
	tableSize := config.TableSize
	historyLength := uint64(math.Log2(float64(tableSize)))
	initialState := utils.StateMap[config.InitialState]

	return &Gshare{
		config:                config,
		globalHistoryRegister: 0,
		patternHistoryTable:   make(map[uint64]utils.State, tableSize),
		historyLength:         historyLength,
		metadata:              metadata,
		initialState:          initialState,
	}
}

// Predict predicts the outcome of the branches in the given instructions
func (gs *Gshare) Predict(instructions []instruction.Instruction) utils.Prediction {
	totalBranches := 0
	mispredictions := 0

	// Create a mask to extract the index bits from the lower bits of the PC address
	mask := uint64((1 << gs.historyLength) - 1)

	for _, instruction := range instructions {
		// Check if the instruction is a branch instruction
		if instruction.Conditional {
			totalBranches++
			pcAddress, err := strconv.ParseUint(instruction.PCAddress, 16, 64)
			if err != nil {
				panic(err)
			}

			// Extract the lower index bits from the PC address
			pcAddress = pcAddress & mask
			taken := instruction.Taken
			isMispredicted := false

			// Update the global history register based on the outcome of the branch instruction
			gs.updateGlobalHistoryRegister(taken)

			// Get the index for the pattern history table and initialize the state if it does not exist
			index := gs.getIndex(pcAddress)
			if _, ok := gs.patternHistoryTable[index]; !ok {
				gs.patternHistoryTable[index] = gs.initialState
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

			// Update the metadata
			gs.UpdateMetaData(instruction, isMispredicted)
		}
	}

	prediction := utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}

	return prediction
}

// getIndex calculates the index for the pattern history table based on the PC address
// and the global history register. It uses the lower bits of the PC address and XORs it
// with the global history register to get the index.
func (gs *Gshare) getIndex(pcAddress uint64) uint64 {
	var addressBits uint64 = uint64(pcAddress) & ((1 << gs.historyLength) - 1)
	return addressBits ^ gs.globalHistoryRegister
}

// updateGlobalHistoryRegister updates the global history register based on the outcome of the branch instruction
// It shifts the global history register to the left by 1 bit and if the branch is taken, it sets the least significant
// bit to 1. It then masks the global history register to keep only the lower bits based on the history length.
func (gs *Gshare) updateGlobalHistoryRegister(taken bool) {
	gs.globalHistoryRegister <<= 1
	if taken {
		gs.globalHistoryRegister |= 1
	}
	gs.globalHistoryRegister &= (1 << gs.historyLength) - 1
}

// UpdateMetaData updates the metadata for the given instruction, based on whether it was mispredicted or not
func (gs *Gshare) UpdateMetaData(instruction instruction.Instruction, isMispredicted bool) {
	if gs.metadata.Exists(instruction.PCAddress) {
		gs.metadata.Update(instruction, isMispredicted)
	} else {
		gs.metadata.AddBranch(instruction.PCAddress)
		gs.metadata.Update(instruction, isMispredicted)
	}
}

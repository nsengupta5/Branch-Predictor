package branchpred

import (
	"math"
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type TwoBitProfiled struct {
	config           *utils.TwoBitProfiledConfig
	metadata         *utils.MetaData
	indexBitCount    uint64
	predictionTable  map[uint64]utils.State
	initialState     utils.State
	strataSize       uint64
	strataProportion float64
}

// NewProfiled creates a new instance of the Profiled struct
func NewTwoBitProfiled(config *utils.TwoBitProfiledConfig, metadata *utils.MetaData) *TwoBitProfiled {
	tableSize := config.TableSize
	indexBitCount := uint64(math.Log2(float64(tableSize)))
	initialState, ok := utils.StateMap[config.InitialState]
	if !ok {
		panic("Invalid initial state")
	}

	return &TwoBitProfiled{
		config:           config,
		metadata:         metadata,
		indexBitCount:    indexBitCount,
		predictionTable:  make(map[uint64]utils.State, tableSize),
		initialState:     initialState,
		strataSize:       config.StrataSize,
		strataProportion: config.StrataProportion,
	}
}

// getProfilerInstructions returns a subset of the instructions to be used for profiling
// The subset of the instruction is obtained using stratified sampling, where the instructions
// are divided into strata and a sample of the instructions is taken from each stratum
func (tbp *TwoBitProfiled) getProfilerInstructions(instructions []instruction.Instruction) []instruction.Instruction {
	totalStrata := uint64(len(instructions)) / tbp.strataSize
	sampledInstructions := make([]instruction.Instruction, 0)

	for i := uint64(0); i < totalStrata; i++ {
		sampleSize := uint64(float64(tbp.strataSize) * tbp.strataProportion)
		start := i * tbp.strataSize
		end := start + sampleSize

		if end > uint64(len(instructions)) {
			end = uint64(len(instructions))
		}

		// Take the first sampleSize instructions from each stratum
		sampledInstructions = append(sampledInstructions, instructions[start:end]...)
	}

	return sampledInstructions
}

// Profile profiles the branches in the given instructions
func (tbp *TwoBitProfiled) Profile(instructions []instruction.Instruction) {
	profilerInstructions := tbp.getProfilerInstructions(instructions)
	totalBranches := 0
	mispredictions := 0

	// Create a mask to extract the index bits from the lower bits of the PC address
	mask := uint64((1 << tbp.indexBitCount) - 1)

	for _, instruction := range profilerInstructions {
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
			direct := instruction.Direct
			isMispredicted := false

			// Add the PC address to the prediction table if it does not exist
			if _, ok := tbp.predictionTable[pcAddress]; !ok {
				tbp.predictionTable[pcAddress] = tbp.initialState
			}

			// Update the state of the branch instruction based on the prediction table
			switch tbp.predictionTable[pcAddress] {
			case utils.StronglyNotTaken:
				if taken {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.WeaklyNotTaken
					isMispredicted = true
				}
			case utils.WeaklyNotTaken:
				if taken {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.StronglyTaken
					isMispredicted = true
				} else {
					tbp.predictionTable[pcAddress] = utils.StronglyNotTaken
				}
			case utils.WeaklyTaken:
				if taken {
					tbp.predictionTable[pcAddress] = utils.StronglyTaken
				} else {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.StronglyNotTaken
					isMispredicted = true
				}
			case utils.StronglyTaken:
				if !taken {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.WeaklyTaken
					isMispredicted = true
				}
			}

			// Update the metadata
			tbp.metadata.Update(taken, direct, pcAddress, isMispredicted, tbp.predictionTable[pcAddress])
		}
	}

	// Initialize the most frequent state for each branch in the metadata
	tbp.metadata.InitializeMostFreqState()

	// Reset the prediction table
	tbp.predictionTable = make(map[uint64]utils.State, len(tbp.predictionTable))
}

// Predict predicts the outcome of the branches in the given instructions
func (tbp *TwoBitProfiled) Predict(instructions []instruction.Instruction) utils.Prediction {
	tbp.Profile(instructions)
	totalBranches := 0
	mispredictions := 0

	// Create a mask to extract the index bits from the lower bits of the PC address
	mask := uint64((1 << tbp.indexBitCount) - 1)

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

			// Add the PC address to the prediction table if it does not exist
			if _, ok := tbp.predictionTable[pcAddress]; !ok {
				branch, ok := tbp.metadata.BranchAddress[pcAddress]
				// If the branch address does not exist in the metadata, use the initial state
				if !ok {
					tbp.predictionTable[pcAddress] = tbp.initialState
				} else {
					// Otherwise, use the most frequent state for the branch address
					// found from the profiling phase
					mostFreqState := branch.MostFreqState
					tbp.predictionTable[pcAddress] = mostFreqState
				}
			}

			// Update the state of the branch instruction based on the prediction table
			switch tbp.predictionTable[pcAddress] {
			case utils.StronglyNotTaken:
				if taken {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.WeaklyNotTaken
				}
			case utils.WeaklyNotTaken:
				if taken {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.StronglyTaken
				} else {
					tbp.predictionTable[pcAddress] = utils.StronglyNotTaken
				}
			case utils.WeaklyTaken:
				if taken {
					tbp.predictionTable[pcAddress] = utils.StronglyTaken
				} else {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.StronglyNotTaken
				}
			case utils.StronglyTaken:
				if !taken {
					mispredictions++
					tbp.predictionTable[pcAddress] = utils.WeaklyTaken
				}
			}
		}
	}

	return utils.Prediction{
		Mispredictions: mispredictions,
		Count:          totalBranches,
	}
}

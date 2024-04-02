package branchpred

import (
	"strconv"

	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type Gshare struct {
	globalHistoryRegister uint32
	patternHistoryTable   map[uint64]State
	historyLength         uint32
	gshareTableSize       uint64
}

func NewGshare(historyLength uint32, tableSize uint64) *Gshare {
	return &Gshare{
		globalHistoryRegister: 0,
		patternHistoryTable:   make(map[uint64]State, tableSize),
		historyLength:         uint32(historyLength),
		gshareTableSize:       tableSize,
	}
}

func (gs *Gshare) Predict(instructions []instruction.Instruction) Prediction {
	return Prediction{}
}

func (gs *Gshare) getIndex(pcAddress string) uint32 {

	pcAddressInt, err := strconv.ParseUint(pcAddress, 16, 32)
	if err != nil {
		panic(err)
	}
	var addressBits uint32 = uint32(pcAddressInt) & ((1 << gs.historyLength) - 1)
	return addressBits ^ gs.globalHistoryRegister
}

func (gs *Gshare) updateGlobalHistoryRegister(taken bool) {
	gs.globalHistoryRegister <<= 1
	if taken {
		gs.globalHistoryRegister |= 1
	}
	gs.globalHistoryRegister &= (1 << gs.historyLength) - 1
}

package branchpred

type State uint8

const (
	StronglyNotTaken State = iota // 00
	WeaklyNotTaken                // 01
	WeaklyTaken                   // 10
	StronglyTaken                 // 11
)

type TwoBit struct {
	predictionTable map[uint32]State
}

func NewTwoBit(tableSize int) *TwoBit {
	predictor := &TwoBit{
		predictionTable: make(map[uint32]State),
	}

	for i := 0; i < tableSize; i++ {
		predictor.predictionTable[uint32(i)] = StronglyNotTaken
	}

	return predictor
}

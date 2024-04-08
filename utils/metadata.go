package utils

// MetaData represents the metadata produced by the simulator
type MetaData struct {
	BranchAddress map[uint64]BranchInfo `json:"branch_address"`
}

// BranchInfo represents the branch specific information when making predictions
type BranchInfo struct {
	TakenRecord         []bool  `json:"taken_record"`
	DirectRecord        []bool  `json:"direct_record"`
	MispredictionRecord []bool  `json:"misprediction_record"`
	StateRecord         []State `json:"state_record"`
	MostFreqState       State   `json:"most_freq_state"`
}

// NewMetaData creates a new metadata object
func NewMetaData(tableSize uint64) *MetaData {
	return &MetaData{
		BranchAddress: make(map[uint64]BranchInfo, tableSize),
	}
}

// AddBranch adds a new branch to the metadata
func (md *MetaData) AddBranch(branchAddress uint64) {
	md.BranchAddress[branchAddress] = BranchInfo{
		TakenRecord:         make([]bool, 0),
		DirectRecord:        make([]bool, 0),
		MispredictionRecord: make([]bool, 0),
	}
}

// Update updates the metadata with the new information
func (md *MetaData) Update(taken bool, direct bool, address uint64, misprediction bool, state State) {
	branchInfo, ok := md.BranchAddress[address]

	// Create a new branch if it does not exist
	if !ok {
		md.AddBranch(address)
		branchInfo = md.BranchAddress[address]

	}
	branchInfo.TakenRecord = append(branchInfo.TakenRecord, taken)
	branchInfo.DirectRecord = append(branchInfo.DirectRecord, direct)
	branchInfo.MispredictionRecord = append(branchInfo.MispredictionRecord, misprediction)

	branchInfo.StateRecord = append(branchInfo.StateRecord, state)

	md.BranchAddress[address] = branchInfo
}

func (md *MetaData) InitializeMostFreqState() {
	for address, branch := range md.BranchAddress {
		counts := make(map[State]int)
		for _, state := range branch.StateRecord {
			counts[state]++
		}

		maxVal := 0
		var mostFreqState State
		for state, count := range counts {
			if count > maxVal {
				maxVal = count
				mostFreqState = state
			}
		}

		updatedBranch := branch
		updatedBranch.MostFreqState = mostFreqState
		md.BranchAddress[address] = updatedBranch
	}
}

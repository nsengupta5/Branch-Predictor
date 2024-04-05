package utils

import (
	"github.com/nsengupta5/Branch-Predictor/instruction"
)

type MetaData struct {
	BranchAddress map[string]BranchInfo `json:"branch_address"`
}

type BranchInfo struct {
	TakenRecord         []bool `json:"taken_record"`
	DirectRecord        []bool `json:"direct_record"`
	MispredictionRecord []bool `json:"misprediction_record"`
}

func NewMetaData(tableSize uint64) *MetaData {
	return &MetaData{
		BranchAddress: make(map[string]BranchInfo, tableSize),
	}
}

func (md *MetaData) AddBranch(branchAddress string) {
	md.BranchAddress[branchAddress] = BranchInfo{
		TakenRecord:         make([]bool, 0),
		DirectRecord:        make([]bool, 0),
		MispredictionRecord: make([]bool, 0),
	}
}

func (md *MetaData) Update(is instruction.Instruction, misprediction bool) {
	branchInfo, ok := md.BranchAddress[is.PCAddress]
	if !ok {
		md.AddBranch(is.PCAddress)
		branchInfo = md.BranchAddress[is.PCAddress]

	}
	branchInfo.TakenRecord = append(branchInfo.TakenRecord, is.Taken)
	branchInfo.DirectRecord = append(branchInfo.DirectRecord, is.Direct)
	branchInfo.MispredictionRecord = append(branchInfo.MispredictionRecord, misprediction)

	md.BranchAddress[is.PCAddress] = branchInfo
}

// Exists checks if the branch address exists in the metadata
func (md *MetaData) Exists(branchAddress string) bool {
	_, ok := md.BranchAddress[branchAddress]
	return ok
}

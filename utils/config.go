package utils

import (
	"errors"
)

// Config represents the configuration of the branch predictor
type Config interface {
	Validate() error
}

// AlwaysTakenConfig represents the configuration of the always taken branch predictor
type AlwaysTakenConfig struct {
}

// TwoBitConfig represents the configuration of the two-bit branch predictor
type TwoBitConfig struct {
	TableSize    uint64 `json:"table_size"`
	InitialState string `json:"initial_state"`
}

// GShareConfig represents the configuration of the gshare branch predictor
type GShareConfig struct {
	TableSize    uint64 `json:"table_size"`
	InitialState string `json:"initial_state"`
}

// TwoBitProfiledConfig represents the configuration of the two-bit profiled branch predictor
type TwoBitProfiledConfig struct {
	TableSize        uint64  `json:"table_size"`
	InitialState     string  `json:"initial_state"`
	StrataSize       uint64  `json:"strata_size"`
	StrataProportion float64 `json:"strata_proportion"`
}

// No validation required for always taken predictor
func (atc AlwaysTakenConfig) Validate() error {
	return nil
}

// Validate the two-bit configuration, ensuring that the table size is valid
func (tbc TwoBitConfig) Validate() error {
	switch tbc.TableSize {
	case 512, 1024, 2048, 4096:
		return nil
	default:
		return errors.New("Invalid table size. Please provide a valid table size")
	}
}

// Validate the gshare configuration, ensuring that the table size is valid
func (gsc GShareConfig) Validate() error {
	switch gsc.TableSize {
	case 512, 1024, 2048, 4096:
		return nil
	default:
		return errors.New("Invalid table size. Please provide a valid table size")
	}
}

// Validate the gshare profiled configuration, ensuring that the table size and train proportion are valid
func (tbpc TwoBitProfiledConfig) Validate() error {
	switch tbpc.TableSize {
	case 512, 1024, 2048, 4096:
	default:
		return errors.New("Invalid table size. Please provide a valid table size")
	}

	if tbpc.StrataProportion < 0 || tbpc.StrataProportion > 1 {
		return errors.New("Invalid train proportion. Please provide a valid train proportion")
	}

	return nil
}

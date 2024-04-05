package utils

import (
	"errors"
)

type Config interface {
	Validate() error
}

type AlwaysTakenConfig struct {
}

type TwoBitConfig struct {
	TableSize    uint64 `json:"table_size"`
	InitialState string `json:"initial_state"`
}

type GShareConfig struct {
	TableSize    uint64 `json:"table_size"`
	InitialState string `json:"initial_state"`
}

func (atc AlwaysTakenConfig) Validate() error {
	return nil
}

func (tbc TwoBitConfig) Validate() error {
	switch tbc.TableSize {
	case 8, 512, 1024, 2048, 4096:
		return nil
	default:
		return errors.New("Invalid table size. Please provide a valid table size")
	}
}

func (gsc GShareConfig) Validate() error {
	switch gsc.TableSize {
	case 8, 512, 1024, 2048, 4096:
		return nil
	default:
		return errors.New("Invalid table size. Please provide a valid table size")
	}
}

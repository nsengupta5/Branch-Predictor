package utils

import (
	"errors"
	"fmt"
)

type Config interface {
	Validate() error
	String() string
}

type AlwaysTakenConfig struct {
	Name string `json:"name"`
}

type TwoBitConfig struct {
	Name         string `json:"name"`
	TableSize    uint64 `json:"table_size"`
	InitialState string `json:"initial_state"`
}

type GShareConfig struct {
	Name         string `json:"name"`
	TableSize    uint64 `json:"table_size"`
	InitialState string `json:"initial_state"`
	HistorySize  uint64 `json:"history_size"`
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

func (atc AlwaysTakenConfig) String() string {
	return fmt.Sprintf("Name: %s\n", atc.Name)
}

func (tbc TwoBitConfig) String() string {
	return fmt.Sprintf("Name: %s\n", tbc.Name) +
		fmt.Sprintf("Table Size: %d\n", tbc.TableSize) +
		fmt.Sprintf("Initial State: %s\n", tbc.InitialState)
}

func (gsc GShareConfig) String() string {
	return fmt.Sprintf("Name: %s\n", gsc.Name) +
		fmt.Sprintf("Table Size: %d\n", gsc.TableSize) +
		fmt.Sprintf("Initial State: %s\n", gsc.InitialState) +
		fmt.Sprintf("History Size: %d\n", gsc.HistorySize)
}

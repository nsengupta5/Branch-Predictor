package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nsengupta5/Branch-Predictor/branchpred"
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type BPConfig struct {
	Algorithm string          `json:"algorithm"`
	MaxLines  int             `json:"max_lines"`
	Config    json.RawMessage `json:"config"`
}

func main() {
	traceFile := os.Args[1]
	configFile := os.Args[2]

	var bpConfig BPConfig = getBPConfig(configFile)

	instructions, err := instruction.ReadTraceFile(traceFile, bpConfig.MaxLines)
	if err != nil {
		panic(err)
	}

	var config utils.Config = getAlgoConfig(&bpConfig)
	simulator := branchpred.NewBranchPredictor(config)

	result := simulator.Predict(instructions)
	printResults(&bpConfig, result)
	simulator.ExportMetaData()
}

func getAlgoConfig(bpConfig *BPConfig) utils.Config {
	switch bpConfig.Algorithm {
	case "always-taken":
		var config utils.AlwaysTakenConfig
		err := json.Unmarshal(bpConfig.Config, &config)
		if err != nil {
			panic(err)
		}
		return config
	case "two-bit":
		var config utils.TwoBitConfig
		err := json.Unmarshal(bpConfig.Config, &config)
		if err != nil {
			panic(err)
		}
		return config
	case "gshare":
		var config utils.GShareConfig
		err := json.Unmarshal(bpConfig.Config, &config)
		if err != nil {
			panic(err)
		}
		return config
	default:
		panic("Invalid algorithm")
	}
}

func getBPConfig(configFile string) BPConfig {
	file, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var bpConfig BPConfig
	err = json.Unmarshal(file, &bpConfig)
	if err != nil {
		panic(err)
	}

	return bpConfig
}

func printResults(bpConfig *BPConfig, result utils.Prediction) {
	fullResult := result.GeneratePredictionFull()
	output := make(map[string]interface{})
	output["config"] = bpConfig
	output["result"] = fullResult

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	outputFile := fmt.Sprintf("outputs/results/%s.json", bpConfig.Algorithm)

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(jsonOutput)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Results written to %s\n", outputFile)
}

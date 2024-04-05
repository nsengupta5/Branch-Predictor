package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nsengupta5/Branch-Predictor/branchpred"
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

type BPConfig struct {
	Algorithm string          `json:"algorithm"`
	MaxLines  int             `json:"max_lines"`
	Configs   json.RawMessage `json:"configs"`
}

func main() {
	traceFile := os.Args[1]
	configFile := os.Args[2]

	var bpConfig BPConfig = getBPConfig(configFile)

	instructions, err := instruction.ReadTraceFile(traceFile, bpConfig.MaxLines)
	if err != nil {
		panic(err)
	}

	var results []map[string]interface{}
	var configs []utils.Config = getAlgoConfig(&bpConfig)
	for _, config := range configs {
		simulator := branchpred.NewBranchPredictor(config)
		result := simulator.Predict(instructions)
		output := make(map[string]interface{})

		fullResult := result.GeneratePredictionFull()
		output["config"] = config
		output["result"] = fullResult
		results = append(results, output)
	}
	exportResults(results, &bpConfig, traceFile)
}

func getAlgoConfig(bpConfig *BPConfig) []utils.Config {
	switch bpConfig.Algorithm {
	case "always-taken":
		var atConfigs []utils.AlwaysTakenConfig
		err := json.Unmarshal(bpConfig.Configs, &atConfigs)
		if err != nil {
			panic(err)
		}

		var configs []utils.Config
		for _, config := range atConfigs {
			configs = append(configs, config)
		}
		return configs
	case "two-bit":
		var tbConfigs []utils.TwoBitConfig
		err := json.Unmarshal(bpConfig.Configs, &tbConfigs)
		if err != nil {
			panic(err)
		}

		var configs []utils.Config
		for _, config := range tbConfigs {
			configs = append(configs, config)
		}
		return configs
	case "gshare":
		var gsConfigs []utils.GShareConfig
		err := json.Unmarshal(bpConfig.Configs, &gsConfigs)
		if err != nil {
			panic(err)
		}

		var configs []utils.Config
		for _, config := range gsConfigs {
			configs = append(configs, config)
		}
		return configs
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

func exportResults(results []map[string]interface{}, bpConfig *BPConfig, traceFile string) {
	// Split the trace file name to get the trace file name without the extension
	traceFile = strings.Split(traceFile, "/")[1]
	traceFile = strings.Split(traceFile, ".")[0]

	outputs := map[string]interface{}{
		"run": map[string]interface{}{
			"algorithm":  bpConfig.Algorithm,
			"max_lines":  bpConfig.MaxLines,
			"trace_file": traceFile,
		},
		"stats": results,
	}

	jsonOutput, err := json.MarshalIndent(outputs, "", "  ")
	if err != nil {
		panic(err)
	}

	outputFile := fmt.Sprintf("outputs/results/%s/%s.json", bpConfig.Algorithm, traceFile)
	err = os.WriteFile(outputFile, jsonOutput, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Results written to %s\n", outputFile)
}

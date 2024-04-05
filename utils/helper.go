package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type BPConfig struct {
	Algorithm string          `json:"algorithm"`
	MaxLines  int             `json:"max_lines"`
	Configs   json.RawMessage `json:"configs"`
}

func GetAlgoConfig(bpConfig *BPConfig) []Config {
	switch bpConfig.Algorithm {
	case "always-taken":
		var atConfigs []AlwaysTakenConfig
		err := json.Unmarshal(bpConfig.Configs, &atConfigs)
		if err != nil {
			panic(err)
		}

		var configs []Config
		for _, config := range atConfigs {
			configs = append(configs, config)
		}
		return configs
	case "two-bit":
		var tbConfigs []TwoBitConfig
		err := json.Unmarshal(bpConfig.Configs, &tbConfigs)
		if err != nil {
			panic(err)
		}

		var configs []Config
		for _, config := range tbConfigs {
			configs = append(configs, config)
		}
		return configs
	case "gshare":
		var gsConfigs []GShareConfig
		err := json.Unmarshal(bpConfig.Configs, &gsConfigs)
		if err != nil {
			panic(err)
		}

		var configs []Config
		for _, config := range gsConfigs {
			configs = append(configs, config)
		}
		return configs
	default:
		panic("Invalid algorithm")
	}
}

func GetBPConfig(configFile string) BPConfig {
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

func ExportResults(results []map[string]interface{}, bpConfig *BPConfig, traceFile string) {
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
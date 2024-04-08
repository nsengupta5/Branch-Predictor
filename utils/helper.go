package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// BPConfig represents the configuration of the branch predictor
type BPConfig struct {
	Algorithm    string          `json:"algorithm"`
	MaxLinesList []int           `json:"max_lines"`
	Configs      json.RawMessage `json:"configs"`
}

type Output struct {
	Algorithm string                   `json:"algorithm"`
	TraceFile string                   `json:"trace_file"`
	Runs      []map[string]interface{} `json:"runs"`
}

// GetAlgoConfig returns the configuration for the specified algorithm
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

	case "two-bit-profiled":
		var tbpConfigs []TwoBitProfiledConfig
		err := json.Unmarshal(bpConfig.Configs, &tbpConfigs)
		if err != nil {
			panic(err)
		}

		var configs []Config
		for _, config := range tbpConfigs {
			configs = append(configs, config)
		}

		return configs

	default:
		panic("Invalid algorithm")
	}
}

// GetBPConfig reads the branch predictor configuration from the specified file
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

// ExportResults exports the results of the branch predictor simulation to a JSON file
func ExportResults(runs []map[string]interface{}, bpConfig *BPConfig, traceFile string) {

	// Split the trace file name to get the trace file name without the extension
	traceFile = strings.Split(traceFile, "/")[1]
	traceFile = strings.Split(traceFile, ".")[0]

	output := Output{
		Algorithm: bpConfig.Algorithm,
		TraceFile: traceFile,
		Runs:      runs,
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
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

// ExportMetadata exports the metadata of the branch predictor simulation to a JSON file
func ExportMetadata(metadataRuns []map[string]interface{}, bpConfig *BPConfig, traceFile string) {
	// Export the metadata to a JSON file
	jsonOutput, err := json.MarshalIndent(metadataRuns, "", "    ")
	if err != nil {
		panic(err)
	}

	// Split the trace file name to get the trace file name without the extension
	traceFile = strings.Split(traceFile, "/")[1]
	traceFile = strings.Split(traceFile, ".")[0]

	outputFile := fmt.Sprintf("outputs/metadata/%s/%s.json", bpConfig.Algorithm, traceFile)
	err = os.WriteFile(outputFile, jsonOutput, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Metadata written to %s\n", outputFile)
}

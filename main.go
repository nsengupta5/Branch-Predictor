package main

import (
	"os"

	"github.com/nsengupta5/Branch-Predictor/branchpred"
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

func main() {
	// Read the trace file and configuration file from the command line arguments
	traceFile := os.Args[1]
	configFile := os.Args[2]

	var bpConfig utils.BPConfig = utils.GetBPConfig(configFile)

	var runs []map[string]interface{}
	var metadataRuns []map[string]interface{}

	// Run the branch predictor simulation for each maxLines value
	for _, maxLines := range bpConfig.MaxLinesList {
		// Read the instructions from the trace file
		instructions, err := instruction.ReadTraceFile(traceFile, maxLines)
		if err != nil {
			panic(err)
		}

		var results []map[string]interface{}
		var metadatas []map[string]interface{}
		var configs []utils.Config = utils.GetAlgoConfig(&bpConfig)

		// Run the branch predictor simulation for each configuration
		for _, config := range configs {
			simulator := branchpred.NewBranchPredictor(config)
			result := simulator.Predict(instructions)
			output := make(map[string]interface{})
			metadata := make(map[string]interface{})

			fullResult := result.GeneratePredictionFull()
			output["config"] = config
			output["result"] = fullResult

			metadata["config"] = config
			metadata["result"] = simulator.Metadata

			results = append(results, output)
			metadatas = append(metadatas, metadata)
		}

		outputs := map[string]interface{}{
			"max_lines": maxLines,
			"stats":     results,
		}

		metadataOutputs := map[string]interface{}{
			"max_lines": maxLines,
			"stats":     metadatas,
		}
		runs = append(runs, outputs)
		metadataRuns = append(metadataRuns, metadataOutputs)
	}

	// Export the results to a JSON file
	utils.ExportResults(runs, &bpConfig, traceFile)
	// utils.ExportMetadata(metadataRuns, &bpConfig, traceFile)
}

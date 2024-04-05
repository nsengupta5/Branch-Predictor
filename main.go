package main

import (
	"os"

	"github.com/nsengupta5/Branch-Predictor/branchpred"
	"github.com/nsengupta5/Branch-Predictor/instruction"
	"github.com/nsengupta5/Branch-Predictor/utils"
)

func main() {
	traceFile := os.Args[1]
	configFile := os.Args[2]

	var bpConfig utils.BPConfig = utils.GetBPConfig(configFile)

	instructions, err := instruction.ReadTraceFile(traceFile, bpConfig.MaxLines)
	if err != nil {
		panic(err)
	}

	var results []map[string]interface{}
	var configs []utils.Config = utils.GetAlgoConfig(&bpConfig)
	for _, config := range configs {
		simulator := branchpred.NewBranchPredictor(config)
		result := simulator.Predict(instructions)
		output := make(map[string]interface{})

		fullResult := result.GeneratePredictionFull()
		output["config"] = config
		output["result"] = fullResult
		results = append(results, output)
	}
	utils.ExportResults(results, &bpConfig, traceFile)
}

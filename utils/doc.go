// Package utils provides the components of the branch prediction simulator that are
// are shared across multiple packages and do not hold the logic of the simulator.
// These components include the configuration of the simulator, the states utilized
// by the simulator, the metadata produced from the simulator and the predictions structure
// that holds the predictions made by the simulator.
//
// The package contains the following files:
//
// config.go: Contains the configuration of the simulator. The configuration includes
// the parameters that the algorithm uses to make predictions.
//
// metadata.go: Contains the metadata produced by the simulator. The metadata includes
// information of the branch specific information when making predictions.
//
// predictions.go: Contains the predictions made by the simulator. The predictions include
// the number of mispredictions, the number of correct predictions, the accuracy of the
// predictions and the misprediction rate.
//
// state.go: Contains the states utilized by the simulator. The states include the states
// of the branch predictor, specifically for the two-bit and gshare predictors.
package utils

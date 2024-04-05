package utils

import "fmt"

// Prediction represents the predictions made by the simulator
type Prediction struct {
	Mispredictions int
	Count          int
}

// PredictionFull represents the full prediction details made by the simulator
type PredictionFull struct {
	Mispredictions    int    `json:"mispredictions"`
	Count             int    `json:"count"`
	Correct           int    `json:"correct"`
	Accuracy          string `json:"accuracy"`
	MispredictionRate string `json:"misprediction_rate"`
}

// Returns the accuracy of the predictions
func (p Prediction) Accuracy() float64 {
	return p.PredictionRate() * 100
}

// Returns the number of correct predictions
func (p Prediction) Correct() int {
	return p.Count - p.Mispredictions
}

// Returns the misprediction rate of the predictions
func (p Prediction) MispredictionRate() float64 {
	return float64(p.Mispredictions) / float64(p.Count) * 100
}

// Returns the prediction rate of the predictions
func (p Prediction) PredictionRate() float64 {
	return float64(p.Correct()) / float64(p.Count)
}

// GeneratePredictionFull generates the full prediction details
func (p Prediction) GeneratePredictionFull() PredictionFull {
	return PredictionFull{
		Mispredictions:    p.Mispredictions,
		Count:             p.Count,
		Correct:           p.Correct(),
		Accuracy:          fmt.Sprintf("%.2f%%", p.Accuracy()),
		MispredictionRate: fmt.Sprintf("%.2f%%", p.MispredictionRate()),
	}
}

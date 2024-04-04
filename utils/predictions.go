package utils

import "fmt"

type Prediction struct {
	Mispredictions int
	Count          int
}

type PredictionFull struct {
	Mispredictions    int    `json:"mispredictions"`
	Count             int    `json:"count"`
	Correct           int    `json:"correct"`
	Accuracy          string `json:"accuracy"`
	MispredictionRate string `json:"misprediction_rate"`
}

func (p Prediction) Accuracy() float64 {
	return p.PredictionRate() * 100
}

func (p Prediction) Correct() int {
	return p.Count - p.Mispredictions
}

func (p Prediction) MispredictionRate() float64 {
	return float64(p.Mispredictions) / float64(p.Count) * 100
}

func (p Prediction) PredictionRate() float64 {
	return float64(p.Correct()) / float64(p.Count)
}

func (p Prediction) GeneratePredictionFull() PredictionFull {
	return PredictionFull{
		Mispredictions:    p.Mispredictions,
		Count:             p.Count,
		Correct:           p.Correct(),
		Accuracy:          fmt.Sprintf("%.2f%%", p.Accuracy()),
		MispredictionRate: fmt.Sprintf("%.2f%%", p.MispredictionRate()),
	}
}

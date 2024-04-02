package branchpred

import "fmt"

type Prediction struct {
	Mispredictions int
	Count          int
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

func (p Prediction) String() string {
	return fmt.Sprintf("Total Branches: %d\n", p.Count) +
		fmt.Sprintf("Correct: %d\n", p.Correct()) +
		fmt.Sprintf("Mispredictions: %d\n", p.Mispredictions) +
		fmt.Sprintf("Accuracy: %.2f%%\n", p.Accuracy()) +
		fmt.Sprintf("Misprediction Rate: %.2f%%\n", p.MispredictionRate())
}

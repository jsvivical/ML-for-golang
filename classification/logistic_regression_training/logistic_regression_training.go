package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gonum/matrix/mat64"
)

/*
if err != nil{
	log.Fatal(err)
}
*/
func logisticRegression(features *mat64.Dense, labels []float64, numSteps int, learningRate float64) []float64 {
	_, numWeights := features.Dims()
	weights := make([]float64, numWeights)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for idx, _ := range weights {
		weights[idx] = r.Float64()
	}

	for i := 0; i < numSteps; i++ {
		var sumError float64
		for idx, label := range labels {
			featureRow := mat64.Row(nil, idx, features)

			pred := logistic(featureRow[0]*weights[0] + featureRow[1]*weights[1])
			predError := label - pred
			sumError += math.Pow(predError, 2)

			for j := 0; j < len(featureRow); j++ {
				weights[j] += learningRate * predError * pred * (1 - pred) * featureRow[j]
			}
		}
	}
	return weights
}

func logistic(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func main() {
	f, err := os.Open("training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 2

	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	featureData := make([]float64, 2*(len(rawCSVData)-1))
	labels := make([]float64, len(rawCSVData)-1)

	var featureIndex int

	for idx, record := range rawCSVData {
		if idx == 0 {
			continue
		}

		featureVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		featureData[featureIndex] = featureVal

		featureData[featureIndex+1] = 1.0

		featureIndex += 2

		labelVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		labels[idx-1] = labelVal
	}

	features := mat64.NewDense(len(rawCSVData)-1, 2, featureData)

	weights := logisticRegression(features, labels, 100, 0.3)

	formula := " p = 1 / ( 1 + exp(-m1 * FICO.score - m2))"

	fmt.Printf("\n\n%s\nm1 = %0.2f\nm2 = %0.2f\n\n", formula, weights[0], weights[1])
}

package main

/**
 *  Sales = m1*TV + m2*Radio + b
 */

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

func Error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\multi_linear_regression\\training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	trainingData, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}
	//여기에서는 TV와 라디오 수치 + y절편으로
	//판매량에 대한 모델링을 시도한다.

	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")
	r.SetVar(1, "Radio")

	for i, record := range trainingData {

		//헤더는 건너뛴다.
		if i == 0 {
			continue
		}
		//판매량 값을 구문 분석해 읽는다.
		yVal, err := strconv.ParseFloat(record[3], 64)
		Error(err)
		tvVal, err := strconv.ParseFloat(record[0], 64)
		Error(err)
		radioVal, err := strconv.ParseFloat(record[1], 64)
		Error(err)

		//이 값들을 regression 값에 추가한다.
		r.Train(regression.DataPoint(yVal, []float64{tvVal, radioVal}))
	} //for 문 끝
	//회귀분석 모델을 훈련(학습)/적합시킨다.

	r.Run()

	fmt.Printf("\nRegression Formula : %v\n\n", r.Formula)

	//predict 메소드를 활용해 회귀분석 모델을 사용해 예측을 수행
	f, err = os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\multi_linear_regression\\test.csv")
	Error(err)

	reader = csv.NewReader(f)
	reader.FieldsPerRecord = 4
	testData, err := reader.ReadAll()
	Error(err)

	/**
	 * 루프를 통해 테스트 데이터를 모두 읽고 y를 예측한 다음
	 * 평균 절대값 오차(MAE)를 사용해 예측 결과를 평가한다.
	 */

	var MAE float64
	for i, record := range testData {
		//헤더는 건너뜀
		if i == 0 {
			continue
		}

		yObserved, err := strconv.ParseFloat(record[3], 64)
		Error(err)

		tvVal, err := strconv.ParseFloat(record[0], 64)
		Error(err)

		radioVal, err := strconv.ParseFloat(record[1], 64)
		Error(err)
		//훈련된 모델을 이용해서 y를 예측한다.
		yPredicted, err := r.Predict([]float64{tvVal, radioVal})
		//평균 절대값 오차에 값을 추가한다.
		MAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}
	fmt.Printf("MAE = %0.2f\n\n", MAE)

}

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/berkmancenter/ridge"
	"github.com/gonum/matrix/mat64"
)

/**
 * 비선형 및 다른 유형이 회귀분석
 * Sales = m1*TV^1 + m2 * TV^2 + ... + b
 * */

func Error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func predict(tv, radio, newspaper float64) float64 {
	return 3.038 + 0.047*tv + 0.177*radio + 0.001*newspaper
}

func main() {
	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\non_linear_regression\\training.csv")
	Error(err)
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4

	//csv 레코드를 모두 읽는다.
	rawCSVData, err := reader.ReadAll()
	Error(err)

	//featureData는 최종적으로 수치를 나타내는 행렬을 만드는데 사용될 모든 float값을 저장한다.

	featureData := make([]float64, 4*len(rawCSVData))
	yData := make([]float64, len(rawCSVData))

	//featureIndex 및 yIndex는 행렬 값의 현재 인덱스를 추적하는데 사용된다.
	var featureIndex, yIndex int

	for idx, record := range rawCSVData {
		//헤더는 건너뜀
		if idx == 0 {
			continue
		}

		for i, val := range record {
			//값을 flaot로 변환한다.
			valParsed, err := strconv.ParseFloat(val, 64)
			Error(err)

			if i < 3 {
				//모델에 y절편을 추가한다.
				if i == 0 {
					featureData[featureIndex] = 1
					featureIndex++
				}
				featureData[featureIndex] = valParsed
				featureIndex++

			}
			if i == 3 { //float값을 저장하는 슬라이스에 float 값을 추가한다.
				yData[yIndex] = valParsed
				yIndex++
			}
		}
	}
	//회귀분석 모델에 입력될 행렬들을 만든다.
	features := mat64.NewDense(len(rawCSVData), 4, featureData)
	y := mat64.NewVector(len(rawCSVData), yData)

	/**
	 *
	 * 다음으로 독립 변수 및 종속 변수 행렬을 갖는 새 ridge.RidgeRegression을 생성하고 Regress() 메소드를 호출해서 회귀분석 모델을 훈련시킨다.
	 *그 다음 훈련을 거친 회귀분석 공식을 출력할 수 있다.
	 *
	 */

	//1.0의 패널티 항을 갖는 새 RidgeRegression값을 생성한다.
	r := ridge.New(features, y, 1.0)

	//회귀분석 모델을 훈련시킨다.
	r.Regress()

	c1 := r.Coefficients.At(0, 0)
	c2 := r.Coefficients.At(1, 0)
	c3 := r.Coefficients.At(2, 0)
	c4 := r.Coefficients.At(3, 0)
	fmt.Printf("\n\nRegression Formula : \n")
	fmt.Printf("y = %0.3f + %0.3f TV + %0.3fRadio + %0.3f Newpaper\n\n", c1, c2, c3, c4)

	//이제 다음과 같이 predict 함수를 생성해 이 능선 회귀분석 공식을 테스트 해볼 수 있다.
	/**
	 * predict는 TV, 라디오, 신문 값을 기반으로 예측을 수행하기 위해
	 * 훈련된 회귀분석 모델을 사용한다.
	 */

	f, err = os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\non_linear_regression\\test.csv")
	Error(err)
	defer f.Close()

	reader = csv.NewReader(f)

	reader.FieldsPerRecord = 4
	testData, err := reader.ReadAll()
	Error(err)

	/**
	 * 루프를 통해 y를 예측하는 홀드아웃 데이터를 읽고
	 * 평균 절대값 오차를 활용해 수치를 평가한다.
	 */

	var MAE float64
	for i, record := range testData {
		if i == 0 {
			continue
		}

		//판매량을 구문 분석해 읽는다. \
		yObserverd, err := strconv.ParseFloat(record[3], 64)
		Error(err)

		tvVal, err := strconv.ParseFloat(record[0], 64)
		Error(err)

		radioVal, err := strconv.ParseFloat(record[1], 64)
		Error(err)

		newspaperVal, err := strconv.ParseFloat(record[2], 64)
		Error(err)
		//훈련된 모델을 사용해 y값을 예측한다
		yPredicted := predict(tvVal, radioVal, newspaperVal)
		//평균 절대값 오차에 값을 추가한다.
		MAE += math.Abs(yObserverd-yPredicted) / float64(len(testData))
	}

	fmt.Printf("\n\nMAE = %0.3f\n\n", MAE)

}

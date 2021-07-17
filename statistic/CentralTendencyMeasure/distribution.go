package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/montanaflynn/stats"
	"gonum.org/v1/gonum/stat"
)

func main() {
	irisFile, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\statistic\\distribution\\iris.data")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()

	//CSV파일에서 데이터프레임 생성하기
	irisDF := dataframe.ReadCSV(irisFile)

	fmt.Println(irisDF)

	//이 변수에 대한 측정값을 확인하기 위해
	//"sepal_length" 열에서 float값을 가져온다

	sepalLength := irisDF.Col("petal_length").Float()
	// 변수의 평균(Mean) 계산하기
	meanVal := stat.Mean(sepalLength, nil)

	//변수의 최빈값(Mode) 계산하기
	modeVal, modeCount := stat.Mode(sepalLength, nil)

	//변수의 중앙값(median) 계산하기
	medianVal, err := stats.Median(sepalLength)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n꽃받침 길이 요약\n")
	fmt.Printf("평균값 : %.2f\n", meanVal)
	fmt.Printf("최빈값 : %.2f\n", modeVal)
	fmt.Printf("최빈값의 개수 : %d\n", int(modeCount))
	fmt.Printf("중앙값 : %.2f\n", medianVal)

	/************************************************************************************
	참고
	평균값과 중앙값이 비슷하지 않은 경우 높은 값과 낮은 값은 각각 평균값을 높게 또는 낮게 끌어당긴다.
	이런 영향은 중앙값에는 눈에 띄지 않는데, 이런 현상을 기울어진 분포라한다
	*************************************************************************************/

}

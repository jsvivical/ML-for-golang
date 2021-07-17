package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

func main() {
	irisFile, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\statistic\\variance\\iris.data")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()

	irisDF := dataframe.ReadCSV(irisFile)

	//이 변수에 대한 측정값을 확인하기 위해
	//"sepal_length"열에서 float값을 가져온다.

	sepalLength := irisDF.Col("petal_length").Float()

	//변수의 최솟값
	minVal := floats.Min(sepalLength)
	//변수의 최댓값
	maxVal := floats.Max(sepalLength)
	//변수의 범위
	rangeVal := maxVal - minVal
	//변수의  분산
	varianceVal := stat.Variance(sepalLength, nil)
	//변수의 표준편차
	stdDevVal := stat.StdDev(sepalLength, nil)
	//값을 정렬
	inds := make([]int, len(sepalLength))
	floats.Argsort(sepalLength, inds)
	//분위수를 계산
	quant25 := stat.Quantile(0.25, stat.Empirical, sepalLength, nil)
	quant50 := stat.Quantile(0.50, stat.Empirical, sepalLength, nil)
	quant75 := stat.Quantile(0.75, stat.Empirical, sepalLength, nil)

	fmt.Printf("변수의 최솟값 : %.2f\n\n ", minVal)
	fmt.Printf("변수의 최댓값 : %.2f\n\n", maxVal)
	fmt.Printf("변수의 범위 : %.2f\n\n", rangeVal)
	fmt.Printf("변수의  분산 : %.2f\n\n", varianceVal)
	fmt.Printf("변수의 표준편차 : %.2f\n\n", stdDevVal)
	fmt.Printf("25분위 : %.2f\n\n", quant25)
	fmt.Printf("50분위 : %.2f\n\n", quant50)
	fmt.Printf("75분위 : %.2f\n\n", quant75)

	for index, v := range sepalLength {
		fmt.Printf("%10.1f", v)
		if index%10 == 9 {
			fmt.Println()
		}
	}

	for index, v := range inds {
		fmt.Printf("%10d", v)
		if index%10 == 9 {
			fmt.Println()
		}
	}
}

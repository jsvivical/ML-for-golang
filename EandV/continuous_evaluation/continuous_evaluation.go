/***************************************************************************
머신러닝이 자신의 역하릉ㄹ 얼마나 잘 수행하고 있는지를 측정하는 프로세스를
평가(evaluation)이라고 함.
일반적으로 평가해야하는 결과에는 몇 가지 유형이 있다.
1.연속형(continuous) : 총 판매액, 주가, 온도와 같이 연속적인 수치를 가질 수 있는 결과
2.범주형(Categorical) : 유한한 여러 범주에서 하나의 값을 가질 수 있는 사기성 여부, 활동, 이름과 같은 결과

각 유형의 결과에는 각각에 해당하는 측정방법이 있다/
***************************************************************************/

// 1. 연속형 측정 방법

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
)

func main() {

	/*
		1.첫번째 단계 : 오차 구하기
			오차 : 예측된 값에서 실제 관찰된 값이 얼마나 벗어나 있는지에 대한 정보를 준다
			그러나 모든 오차를 개별적으로 확인하는 것은 실용적이지 않다.
			따라서 오차를 전체적으로 이해할 필요가 있다.
			MSE(Mean Squared Error,평균 제곱근 오차), MAE(Mean Absolute Error, 평균 절대값 오차)는
			오차를 종합적으로 보여준다.main()

				MSE : 평균 제곱 편차는 모든 오차를 제곱하고 평균을 구한 값, 제곱이기 때문에 오차값이 더 민감.
				MAE : 모든 오차의 절대값을 평균을 구한 값 , 오차가 변수와 동일한 단위를 유지하기 때문에 예측값과 직접 비교 가능
	*/

	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\EandV\\continuous_evaluation\\continuous_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//열린 파일을 읽는 새 CSV reader를 생성한다.
	reader := csv.NewReader(f)

	//observed 및 predicted 변수는 연속형 데이터 파일로부터
	// 읽어온 관찰된 값 및 예측값에 대한 정보를 저장하는데 사용

	var observed []float64
	var predicted []float64

	line := 1

	for {
		//행을 읽는다.
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		//헤더는 건너뛴다.
		if line == 1 {
			line++
			continue
		}

		//첫번째 행을 읽음 -> 관찰된 값
		observedVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("예상하지 못한 유형으로 인한 %d줄 읽기 실패", line)
			continue
		}
		//두번째 행을 읽음 -> 예측값
		predictedVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("예상치 못한 유형으로 인한 %d줄 읽기 실패", line)
			continue
		}
		//예상되는 유형의 경우 슬라이스에 해당 레코드를 추가한다.
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++

	}

	var MAE float64
	var MSE float64

	for i, oVal := range observed {
		MAE += math.Abs(oVal-predicted[i]) / float64(len(observed))

		MSE += math.Pow(oVal-predicted[i], 2) / float64(len(observed))
	}

	fmt.Printf("\nMAE = %0.2f\n\n", MAE)
	fmt.Printf("\nMSE = %0.2f\n\n", MSE)

	/*
		이 값이 좋은지 아닌지 판단하기 위해서는 이 결과를 관찰된 데이터의 값과 비교해야한다.
		특히 MAE는 2.55이고 관찰된 값의 평균은 14.0이기 때문에 MAE는 관찰된 값의 평균의 약 20%이다.
		이는 문맥을 기준으로 봤을 때 그리 좋지않다.
	*/

	/**********************************************************************************************
	MSE, MAE와 함께 R-Square 또는 결정계수도 연속형값 모델을 측정하는 방법으로 사용된다.
	R-Square는 예측값에서 수지반 관찰값에서 분산의 비율을 측정한다.
	예를 들어, 주식 가격, 금리, 질병의 진행 상황을 예측하려고 할 때 이런 값들은 본직적으로 동일하지 않다.(값이 변한다)
	관측된 값을 통해 이런 변동성을 예측하는 머신 러닝 모델을 만들려고 할 때 수집한 분산의 비율이 R-Square로 표시된다.
	**********************************************************************************************/

	rSquared := stat.RSquaredFrom(observed, predicted, nil)

	fmt.Printf("\n\nR-Square = %0.2f\n\n", rSquared)

	/*
	   R-Square는 비율이고, 비율에서는 높은 비율이 좋다.
	   따라서 37%라는 결과는 좋지 못한 결과.
	*/

}

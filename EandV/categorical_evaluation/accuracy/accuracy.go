package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\EandV\\categorical_evaluation\\accuracy\\labeled.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	/* 	observed 및 predicted 변수는  labeled 데이터 파일에서 읽어온
	   	관찰값과 예측값을 저장하는데 사용된다 */
	var observed []int
	var predicted []int

	// CSV의 레코드에서 값을 읽어노는 작업을 반복하고 정확도를 계산하기 위해 관찰된 값과 예측값을 비교한다.
	//line 변수는 로그를 위해 행의 수를 기록한다.
	line := 1
	for {
		//열에서 예기치 않은 유형에 대한 레코드를 읽는다.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		// 헤더생략
		if line == 1 {
			line++
			continue
		}

		//관측된 값과 예측된 값을 읽음
		observedVal, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("예상하지 못한 유형으로 인한 %d줄 읽기 실패", line)
			continue
		}

		predictedVal, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("예상하지 못한 유형으로 인한 %d줄 읽기 실패", line)
			continue
		}

		//예상되는 유형인 경우 슬라이스에 해당 레코드를 추가
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)

		line++
	}
	//truePosNeg 변수는 TP와 TN값의 횟수를 저장하는 데 사용된다.
	//즉, 예측이 맞은 경우
	var truePosNeg int

	//true positive/ negative 횟수를 누적시킨다.
	for i, oVal := range observed {
		if oVal == predicted[i] {
			truePosNeg++
		}
	}

	accuracy := float64(truePosNeg) / float64(len(observed))

	fmt.Printf("\n정확도 : %0.2f\n\n", accuracy)

	//범주가 2개 이상일 때도 다양한 방법을 통해 계산 가능
	//예를 들어 하나를 양성, 나머지를 음성으로 간주하는 방법 등
	// 정밀도와 재현율도 비슷하게 계산 가능
	//classes  변수는 labeled 데이터에서 3가지 가능한 클래스를 포함한다.

	classes := []int{0, 1, 2}
	//각 클래스에 대해 루프를 통해 작업한다.
	for _, class := range classes {
		//이 변수들은 true positive의 횟수와 false positive, false negative의 횟수를 저장하는데 사용
		var truePos, falsePos, falseNeg int

		for idx, oVal := range observed {
			switch oVal {
			//예측한 값이 해당 클래스였는지를 확인
			case class:
				if predicted[idx] == class {
					truePos++
					continue
				}
				falseNeg++

			//관찰된 값이 다른 클래스인 경우
			//예측값이 false positive 였는지 확인한다.
			default:
				if predicted[idx] == class {
					falsePos++
				}
			}
		}
		precision := float64(truePos) / float64(truePos+falsePos)
		recall := float64(truePos) / float64(truePos+falseNeg)

		fmt.Printf("\n\n정밀도(클래스 %d) : %0.2f\n\n", class, precision)
		fmt.Printf("\n\n재현율(클래스 %d) : %0.2f\n\n", class, recall)
	}

}

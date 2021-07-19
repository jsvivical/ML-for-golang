/**
 * 주어진 자료의 신용점수를 기반으로 특정 금리 이하로 대출을 받을 수 있을지 여부를 예측하는 로지스틱 회귀분석 모델을 만들어라.
 *예를 들어, 12%이하의 금리로 대출을 받고자 한다고 가정해보자,그러면 로지스틱 회귀분석 모델을 부여된 신용점수를 기반으로
 * 해당 대출 금리 이하로 대출을 받을 수 있을지 또는 받을 수 없을지를 알려준다.
 *
 *
 *
 * 데이터 정리 및 데이터 프로파일링
 * 대출 샘플을 살펴보면, 필요한, 정확한 형태로 작성돼 있지 않다. 따라서 다음과 같은 작업을 수행한다.
 * 	1. 이자율 및 FICO 저수 열에서 숫자가 아닌 문자를 제거한다.
 * 	2. 부여된 이자율 기준치를 기반으로 이자율을 두 클래스로 나눈다. 1.0을 첫번째 클래스로 사용하고 0.0을 두번째 사용한다
 * 	3.  FICO 신용 점수에 대해 하나의 값을 선택한다. 다양한 신용 점수가 주어지지만 하나의 값만 필요하다. 평균, 최소, 최대 신용 점수를 선택하는 것이 자연스러운
 * 선택이며 예제에서는 최솟값을 사용할 예정이다.
 * 4.여기에서는 FICO 신용점수를 표준화해  사용한다. 이를 통해 점수값을 0.0에서 1.0 사이로 분포시킬 수 있다.
 *이렇게 하면 데이터의 가독성이 떨어지기 때문에 정당화가 필요하다.package logisticregression
 *
 * 정당화에 좋은 방법은 로지스틱 회귀분석을 훈련할 때 경사하강법을 사용하는 것이다. 경사하강법은 정규화된 데이터에서 더 좋은 성능을 낸다.
 * */

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/gonum/matrix/mat64"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	SCORE_MAX = 830.0
	SCORE_MIN = 640.0
)

func Error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func main() {

	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\classification\\logistic_regression\\loan_data.csv")
	Error(err)
	defer f.Close()

	//열린 파일을 읽는 새 CSV 리더를 생성한다.
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 2

	//CSV레코드를 모두 읽는다.
	rawCSVData, err := reader.ReadAll()
	Error(err)

	//출력파일을 생성한다.
	f, err = os.Create("clean_loan_data.csv")
	Error(err)

	//CSV writer를 생성한다.
	w := csv.NewWriter(f)

	//열을 순차적으로 이동하면서 구문 분석된 값을 쓴다.
	for i, record := range rawCSVData {

		//헤더 처리
		if i == 0 {
			//헤더를 출력파일에 씀
			if err := w.Write(record); err != nil {
				log.Fatal(err)
			}
			continue
		}

		//읽어온 값을 저장하기 위한 슬라이스(slice)를 초기화한다.
		outRecord := make([]string, 2)

		//FICO 점수를 구문분석하고 표준화한다.
		score, err := strconv.ParseFloat(strings.Split(record[0], "-")[0], 64) //"-"이전의 숫자를 읽음
		Error(err)
		outRecord[0] = strconv.FormatFloat((score-SCORE_MIN)/(SCORE_MAX-SCORE_MIN), 'f', 4, 64)

		//이자율 클래스를 구문 분석한다.
		rate, err := strconv.ParseFloat(strings.TrimSuffix(record[1], "%"), 64)
		Error(err)

		if rate <= 12.0 {
			outRecord[1] = "1.0"
			//레코드를 출력 파일에 쓴다.
			if err := w.Write(outRecord); err != nil {
				log.Fatal(err)
			}
			continue
		}
		outRecord[1] = "0.0"

		//레코드를 출력파일에 쓴다.
		if err := w.Write(outRecord); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	/******************************************************************************************************************
	 * 데이터 정리 및 데이터 프로파일링 끝
	 * 이제 FICO 점수와 이자율 데이터에 대한 히스토그램을 생성하고 요약 통계를 계산해서 데이터를 좀 더 직관적으로 살펴본다
	 *
	 * 히스토그램 생성
	 *
	 * ****************************************************************************************************************/

	loanDataFile, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\classification\\logistic_regression\\clean_loan_data.csv")
	Error(err)
	defer loanDataFile.Close()

	//CSV파일로부터 data프레임을 생성한다.
	loanDF := dataframe.ReadCSV(loanDataFile)

	//Describe 메소드를 사용해 모든 열에 대한 요약 통계를 한번에 계산
	loanSummary := loanDF.Describe()
	fmt.Println(loanSummary)

	//데이터 집합의 모든 열에 대한 히스토그램을 생성한다.
	for _, colName := range loanDF.Names() {

		//plotter.Values  값을 생성하고 dataframe에서 해당하는 값으로 plotter.Values를 채운다.
		plotVals := make(plotter.Values, loanDF.Nrow())

		for i, floatVal := range loanDF.Col(colName).Float() {
			plotVals[i] = floatVal
		}

		//도표를 만들고 제목을 설정
		p := plot.New()
		p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

		//원하는 값에 대한 히스토그램을 설정한다.
		h, err := plotter.NewHist(plotVals, 16)
		Error(err)

		//히스토그램을 정규화한다.
		h.Normalize(1)

		//히스토그램을 도표에 추가한다.
		p.Add(h)

		//도표 저장
		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}

	}
	f.Close()
	/***************************************************************************************************************
	 *평균 신용 점수가 706.1.로 다소 높은 것을 볼 수있고, 평균이 0.5 근처에 나타나 클래스 1과 0 사이의 균형이 상당히 좋게 설정되었음.
	 *하지만 클래스 0의 경우가 조금 더 많다.
	 * 즉 12% 이하로 대출을 받을 수 없다는 것을 의미한다.
	 * 이제 이 데이터를 훈련시키기 위해 훈련용 데이터와 테스트용 데이터로 분리한다.
	 *
	 * *************************************************************************************************************/

	//정리된 대출 데이터 집합 파일을 연다.
	f, err = os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\classification\\logistic_regression\\clean_loan_data.csv")
	Error(err)

	//CSV 파일로부터 dataframe을 생성한다.
	//열의 유형을 추론한다.
	loanDF = dataframe.ReadCSV(f)

	//각 집합에서 항목의 수를 계산한다.
	trainingNum := (4 * loanDF.Nrow()) / 5 //80%
	testNum := loanDF.Nrow() / 5           //20%
	if trainingNum+testNum < loanDF.Nrow() {
		trainingNum++
	}

	//훈련용 데이터 집합과 테스트 데이터 집합에서 사용할 인덱스르 생성한다.
	trainingIdx := make([]int, trainingNum)
	testIdx := make([]int, testNum)

	//루프를 통해 훈련용 인덱스를 저장한다.
	for i := 0; i < trainingNum; i++ {
		trainingIdx[i] = i
	}

	//루프를 통해 테스트용 인덱스를 저장한다.
	for i := 0; i < testNum; i++ {
		testIdx[i] = i + trainingNum
	}

	//훈련용, 테스트용 데이터의 dataframe을 생성한다
	trainingDF := loanDF.Subset(trainingIdx)
	testDF := loanDF.Subset(testIdx)

	//데이터를 파일에 쓸 때 사용할 맵을 생성한다.
	setMap := map[int]dataframe.DataFrame{
		0: trainingDF,
		1: testDF,
	}

	//파일을 각각 생성한다.
	for idx, setName := range []string{"training.csv", "test.csv"} {
		//필터링을 거친 데이터 집합 파일을 저장한다.
		f, err := os.Create(setName)
		Error(err)

		//버퍼 writer를 생성한다.
		w := bufio.NewWriter(f)

		//dataframe을 CSV파일로 쓴다.
		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}
	f.Close()
	/*********************************************************************************************************************************
	 * 로지스틱 회귀 분석 작성 완료하고
	 * 데이터 훈련시키기
	 * ********************************************************************************************************************************/

	f, err = os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\데이터분석\\Machine-Learning-With-Go\\Chapter05\\logistic_regression\\example6\\training.csv")
	Error(err)
	defer f.Close()

	reader = csv.NewReader(f)
	reader.FieldsPerRecord = 2

	rawCSVData, err = reader.ReadAll()
	Error(err)

	// featureData와 labels 변수는 최종적으로 모델을 훈련시키는데 사용될 float값을 저장하는 데 사용된다.

	featureData := make([]float64, 2*len(rawCSVData))
	labels := make([]float64, len(rawCSVData))

	//featureIndex 변수는 수치를 저장하는 행렬 값의 현재 인덱스를 추적하는데 사용된다.
	var featureIndex int

	//열을 순차적으로 이동하면서 슬라이스에 float값을 저장한다.
	for idx, record := range rawCSVData {
		if idx == 0 {
			continue
		}

		//FICO 점수 수치를 추가한다.
		featureVal, err := strconv.ParseFloat(record[0], 64)
		Error(err)
		featureData[featureIndex] = featureVal

		//y절편을 추가한다.
		featureData[featureIndex+1] = 1.0

		//수치 열에 대한 인덱스를 증가시킨다.
		featureIndex += 2

		//클래스 레이블을 추가한다.
		labelVal, err := strconv.ParseFloat(record[1], 64)
		Error(err)
		labels[idx-1] = labelVal
	}

	//앞에서 저장한 수치로 행렬을 만든다
	features := mat64.NewDense(len(rawCSVData), 2, featureData)

	weights := logisticRegression(features, labels, 100, 0.3)
	fmt.Println(weights)
	//표준출력을 통해 로지스틱 회귀분석 모델 공식을 출력
	formula := "p = 1 / (1 + exp(-m1 * FICO.score - m2))"
	fmt.Printf("\n\n%s\n\nm1 = %0.2f\nm2 = %0.2f\n\n", formula, weights[0], weights[1])

	/*********************************************************************************************************************
	 * 훈련 완료
	 * 이제 predict 함수를 작성하고 예측하고 평가함
	 ***********************************************************************************************************************/

	f, err = os.Open("test.csv")
	Error(err)
	defer f.Close()

	reader = csv.NewReader(f)

	//observed 와 predicted 변수는 레이블에 지정된 데이터 파일로부터 읽어온 관찰값 및 예측값을 저장하는데 사용
	var observed, predicted []float64

	//line변수는 로그를 위해 열의 수를 추적하는데 사용
	line := 1

	//열에서 예기치 않은 유형을 찾기 위해 레코드를 읽는다.
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		//헤더는 건넌 뜀
		if line == 1 {
			line++
			continue
		}

		//관찰값을 읽는다.
		observedVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		score, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		predictedVal := predict(score, weights)

		//기대하는 유형인 경우, 해당 레코드를 슬라이스(slice)에 추가
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++
	}
	//아래의 변수는 true positivedhk true negative의 횟수를 저장 (즉 예상이 맞은 경우)
	var truePosNeg int

	//true Positive/ Negative 횟수를 누적시킨다.
	for idx, oVal := range observed {
		if oVal == predicted[idx] {
			truePosNeg++
		}
	}
	accuracy := float64(truePosNeg) / float64(len(observed))

	fmt.Printf("\n\n정확도 : %0.2f\n\n", accuracy)

	fmt.Print("신용점수를 입력하세요  : ")
	var fico, yourFico, yourScore float64
	n, err := fmt.Scanf("%f ", &fico)
	if err != nil {
		log.Fatal(n, err)
	}

	for {
		yourFico = (fico - SCORE_MIN) / (SCORE_MAX - SCORE_MIN)
		yourScore = predict(yourFico, weights)

		if yourScore == 1.0 {
			if fico == 9999 {
				break
			}
			fmt.Println("당신은 12%이하의 이자율로 대출을 받을 수 있을 것으로 예상됩니다.")
		} else {
			fmt.Println("당신은 12%이하의 이자율로 이자율로 대출을 받기 어려울 것 같습니다.")
		}

		fmt.Print("신용점수를 입력하세요  : ")
		n, err := fmt.Scanf("%f ", &fico)
		if err != nil {
			log.Fatal(n, err)
		}

	}
}

/**************************************************************************************************************************************
 * 훈련용, 테스트용 데이터 분리 완료
 * 다음으로 로지스틱 회귀 분석 모델을 훈련시키고 테스트 함
 * 로지스틱 회귀분석 모델을 훈련시키는 함수 생성
 * 그 함수는 다음의 기능을 수행
 * 1. FICO 점수 데이터를 독립 변수로 수락한다
 * 2. 로지스틱 회귀분석 모델에 y절편을 추가한다.
 * 3. 로지스틱 회귀분석 모델의 계수를 초기화하고 최적화한다. (확률적 경사하강법  적용 : SGD, stochastic gradient descent)
 * 4. 훈련된 모델을 정의하는 최적화된 가중치를 반환한다.
 *
 *최적화 방법에 대한 구현은 다음과 같은 매개변수를 사용한다.
 *features : gonum mat64.Dense 행렬에 대한 포인터를 나타낸다. 이 행렬에는 예제에서 사용할 독립변수에 대한 열(FICO 점수)과 y절편을 나타내는 1.0의
					 값을 가진 열이 포함된다.
 *
 * labels : features에 해당하는 모든 클래스 레이블을 포함하는 float로 이루어진 슬라이스
 *
 * numSteps : 최적화를 위한 최대 반복 횟수
 *
 *learningRate : 조정 가능한 매개변수로서 최적화의 수렴을 돕는 데 사용된다.
 *
**************************************************************************************************************************************/

func logisticRegression(features *mat64.Dense, labels []float64, numSteps int, learningRate float64) []float64 {

	// Initialize random weights.
	_, numWeights := features.Dims()
	weights := make([]float64, numWeights)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for idx, _ := range weights {
		weights[idx] = r.Float64()
	}

	// Iteratively optimize the weights.
	for i := 0; i < numSteps; i++ {

		// Initialize a variable to accumulate error for this iteration.
		var sumError float64

		// Make predictions for each label and accumlate error.
		for idx, label := range labels {

			// Get the features corresponding to this label.
			featureRow := mat64.Row(nil, idx, features)

			// Calculate the error for this iteration's weights.
			pred := logistic(featureRow[0]*weights[0] + featureRow[1]*weights[1])
			predError := label - pred
			sumError += math.Pow(predError, 2)

			// Update the feature weights.
			for j := 0; j < len(featureRow); j++ {
				weights[j] += learningRate * predError * pred * (1 - pred) * featureRow[j]
			}
		}
	}
	fmt.Println(weights)

	return weights
}

func logistic(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func predict(score float64, weights []float64) float64 {
	p := 1 / (1 + math.Exp(-1*weights[0]*score-weights[1]))

	if p >= 0.5 {
		return 1.0
	}

	return 0.0
}

/* ****************************************************************************************************************************
선형 회귀분석은 가장 간단한 머신 러닝 모델 중 하나
간단하고 해석하기 쉬운 모델일수록 무결성을 관리하기가 쉽기때문에 매우 중요한 이점을 가짐
선형 회귀분석 모델은 해석 가능하기 때문에 데이터 과학자에게 안전하고 생산적인 옵션을 제공할 수 있다.

개요
선형 회귀분석에서는 직선의 방정식을 사용해 종속 변수 y를 독립변수 x로 모델링 한다.
 y = mx + b
 여기서 기울기 m과 b를 알면 반응을 예측할 수 있다.
 이 때 가장 널리 사용되고 간단한 방법은 최소자승법이다.

 최소자승법

 최소자승법을 통해 m과 b를 구하려면 먼저 m과 b의 값을 선택한다. 그런 다음 알고있는 지점(데이터)과 예제 선과의 수직 수직거리를 측정한다.
 이 거리를 오차(error) 또는 잔차(residuals)라고 한다. 평가 및 검증 에서 나오는 오차와 유사하다.
이어서 이런 오차의 제곱의 합을 더한다.
이 오차제곱의 합이 최소화될 때까지 m과 b를 조정한다. 다시말해 훈련을 거친 선형 회귀분석 직선은 이 제곱의 합을 최소화시키는 선이다.
이 오차 제곱의 합을 최소화하는데 가장 보편적으로 사용되는 최적화 기법은 경사하강법이다.

선형 회귀분석 가정 및 함정

선형관계 : 선형 회귀분석은 종속 변수가 독립변수에 선형적으로 의존한다는 것을 가정한다.
정규성 : 변수가 정규분포에 따라 분포되어야 한다.
다중공선성 없음 : 다중공선성은 독립변수들이 실제로는 독립적이지 않다는 것을 의미하는 용어다.
자기상관없음 : 자기상관은 특정 변수가 자기자신 또는 자기자신에서 일부 변경된 버전에 종속된다는 것을 의미하는 용어다
등분산성 :
 ******************************************************************************************************************************/

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	/**
	 * 데이터 프로파일링
	 * 이해할 수 있는 모델을 만들고 결과를 확인하기 위해서는 모든 머신러닝 모델의 제작과정을 프로파일링에서부터 시작
	 * 각각의 변수들이 어떻게 분포되어 있는지 그리고 변수들의 범위와  변동성에 대한 이해가 필요하다.
	 * github.com/go-gota/gota/dataframe 사용
	 */
	advertFile, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\회귀분석\\linear\\Advertising.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer advertFile.Close()

	advertDF := dataframe.ReadCSV(advertFile)

	//Describe 메소드를 사용해 모든 열에 대한 요약 통계를 한번에 계산
	advertSummary := advertDF.Describe()
	fmt.Println(advertSummary)

	/* *************************************************************
		column    TV                  Radio         Newspaper    Sales
	 0: mean     147.042500 23.264000 30.554000  14.022500
	 1: median   149.750000 22.900000 25.750000  12.900000
	 2: stddev   85.854236  14.846809 21.778621  5.217457
	 3: min      0.700000   0.000000  0.300000   1.600000
		******************************************************************/

	//시각적인 이해를 위해 각각의 열에 있는 값에 대한 히스토그램을 만든다

	for _, colName := range advertDF.Names() {

		//plotter.Values값을 생성하고 dataframe의 각 열에 있는 값으로
		//plotter.Values를 채운다
		plotVals := make(plotter.Values, advertDF.Nrow())
		for i, floatVal := range advertDF.Col(colName).Float() {
			plotVals[i] = floatVal
		}
		//도표를 만들고 제목을 설정한다.

		p := plot.New()

		p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)
		//표준법선으로부터 그려지는 값의 히스토그램을 생성한다.

		h, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			log.Fatal(err)
		}
		h.Normalize(1)
		p.Add(h)

		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
	}
	//위의 값들은 정규분포를 따르지 않기 때문에  다음의 결정을 내려야한다.
	//1. 정규분포를 따르는 변수로 변환한 다음, 이렇게 변환된 변수를 선형 회귀분석 모델에 사용
	//2.문제를 해결할 수 있는 다른 데이터를 사용
	//3.선형 회귀분석 가정에 대한 문제를 무시하고 모델을 만든다(권장)

	//목표하는 열을 추출한다.
	yVals := advertDF.Col("Sales").Float()
	for _, colName := range advertDF.Names() {

		//pts변수는 도표에 대한 값을 저장한다.
		pts := make(plotter.XYs, advertDF.Nrow())

		//pts변수를 데이터로 채운다.
		for i, floatVal := range advertDF.Col(colName).Float() {
			pts[i].X = floatVal
			pts[i].Y = yVals[i]
		}

		//도표를 생성한다.
		p := plot.New()
		p.X.Label.Text = colName
		p.Y.Label.Text = "y"
		p.Add(plotter.NewGrid())

		//산점도
		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Fatal(err)
		}
		s.GlyphStyle.Radius = vg.Points(3)

		p.Add(s)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_scatter.png"); err != nil {
			log.Fatal(err)
		}

	}

	//출력 결과 TV와 판매량과의 관계가 가장 분명하다.
	//예측 능력을 가진 선형 회귀분석 모델을 생성할 수 있는지 확인하기 위해 계속 작업
	//훈련용 데이터 집합 및 테스트 집합 만들기

	//각 집합에서 요소의 수를 계산, 80 20 분할 사용
	trainingNum := (4 * advertDF.Nrow()) / 5 // 80
	testNum := advertDF.Nrow() / 5           //20
	if trainingNum+testNum > advertDF.Nrow() {
		trainingNum++
	}
	//훈련용 인덱스와 테스트용 인덱스를 저장할 배열을 생성한다.
	trainingIdx := make([]int, trainingNum)
	testIdx := make([]int, testNum)

	for i := 0; i < trainingNum; i++ {
		trainingIdx[i] = i //앞에서부터 80%
	}

	for i := 0; i < testNum; i++ {
		testIdx[i] = trainingNum + i //이후 20%
	}

	//각 데이터 집합에 대한 데이터프레임을 생성한다.
	trainingDF := advertDF.Subset(trainingIdx)
	testDF := advertDF.Subset(testIdx)

	//데이터를 파일에 쓸 때 사용될 맵을 생성한다.
	setMap := map[int]dataframe.DataFrame{
		0: trainingDF,
		1: testDF,
	}

	//각각의 파일을 생성한다.
	for idx, setName := range []string{"training.csv", "test.csv"} {
		f, err := os.Create(setName)
		if err != nil {
			log.Fatal(err)
		}
		//buffered writer를 생성한다.
		w := bufio.NewWriter(f)

		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}

	}

	//모델 훈련(학습) 시키기
	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\회귀분석\\linear\\training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	//CSV레코드를 모두 읽는다.
	reader.FieldsPerRecord = 4
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	/********************************************************************************************************************************
	 * 여기서는 TV수치와 y절편에 대해 판매량(y)를 모델링 하려고한다.
	 * 따라서 github.com/sajari/regression을 사용해 모델을 훈련(학습) 시키기 위해 필요한 구조체를 생성
	 *********************************************************************************************************************************/

	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")

	//루프를 통해 CSV에서 레코드를 읽고 regression 값에 훈련(학습) 데이터를 추가한다.

	for i, record := range trainingData {
		//헤더는 건너 뜀
		if i == 0 {
			continue
		}

		//판매량 회귀분석 측정값 또는 "y" 값을 구문 분석해 읽는다.
		yVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		//TV값을 구분 분석해 읽는다.
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		//이 값들을 regression에 추가한다.
		r.Train(regression.DataPoint(yVal, []float64{tvVal}))
	}
	//회귀분석 모델을 훈련(학습)/적합한다.
	r.Run()

	//훈련(학습)을 거친 모델 매개변수를 출력한다.
	fmt.Printf("\n\n회귀분석 공식 : %v\n\n", r.Formula)

	//훈련된 모델 평가하기

	f, err = os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\회귀분석\\linear\\test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader = csv.NewReader(f)
	//모든 CSV레코드를 읽는다
	reader.FieldsPerRecord = 4
	testData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//루프를 통해 y를 예측하는 테스트 데이터를 읽고
	//평균 절대값 오차를 활용해 예측된 수치를 평가한다.

	var MAE float64
	for i, record := range testData {
		//헤더는 건너뛴다
		if i == 0 {
			continue
		}

		//관찰된 판매량 또는 "y"값을 구문 분석해 읽는다.
		yObserved, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		//tv값을 구문 분석해 읽는다.
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		//훈련된 모델을 사용해 예측을 수행한다.
		yPredicted, err := r.Predict([]float64{tvVal})

		//MAE에 추가한다.

		MAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}

	//표준출ㄹ력으로 MAE출력
	fmt.Printf("\n\nMAE = %0.2f\n\n", MAE)
	//MAE = 3.01

	//결과 평가
	/**
	 * 평균 판매량은 14.02, 표준편차는 5.02이었음.
	 * 따라서 MAE는 판매량값의 표준편차보다 작고 평균값의 약 20% 정도이기 때문에
	 *  우리가 제작한 모델은 예측 능력을 어느 정도 갖췄다고 할 수 있다
	 */

}

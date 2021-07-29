/************************************************************************************************************************************
 * 자동 회귀 모델의 가정 및 문제점
 *	고정적 : AR모델은 사용되는 시계열이 고정적이라고 가정한다. AR모델을 사용하려면 데이터에서 어떤 경향성이 나타나서는 안된다.
 * 에르고딕성 : 이 용어는 기본적으로 평균 및 분산과 같은 시계열의 통계적 속성이 시간에 따라 변하지 않아야 한다는 것을 의미
 *
 * 이 프로그램은 항공 승객 데이터를 사용하지만, 이 데이터가 고정적이지 않기 때문에 AR모델의 가정을 충족하지 못함.
 * 따라서 차분(Differencing)이라는 기법을 사용해서 시계열 데이터를 고정적으로 만든다.
************************************************************************************************************************************/

package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const dataURL = "https://raw.githubusercontent.com/vincentarelbundock/Rdatasets/master/csv/datasets/AirPassengers.csv"

func acf(x []float64, lag int) float64 {

	xAdj := x[lag:len(x)]
	xLag := x[0:(len(x) - lag)]

	//numerator 변수는 누적된 분자의 값을 저장하는데 사용됨
	//denominator 변수는 누적된 분모의 값을 저장
	var numerator, denominator float64

	xBar := stat.Mean(x, nil)
	//numerator 계산
	for idx, xVal := range xAdj {
		numerator += ((xVal - xBar) * (xLag[idx] - xBar))
	}
	//denominator 계산
	for _, xVal := range x {
		denominator += math.Pow(xVal-xBar, 2)
	}

	return numerator / denominator

}

func pacf(x []float64, lag int) float64 {
	//pacf 함수는 주어진 특정 주기 이전의 값에서 시계열의 편 자기상관을 계산
	//regression모델을 사용해 학습시키기 위한 regression.Regression 값을 생성
	var r regression.Regression
	r.SetObserved("x")

	//현재 및 중간 이전의 값을 모두 정의
	for i := 0; i < lag; i++ {
		r.SetVar(i, "x"+strconv.Itoa(i))
	}

	//데이터 열을 이동시킴
	xAdj := x[lag:len(x)]

	//루프를 통해 회귀분석 모델을 위한 데이터 집합을 생성하는 시계열 데이터를 읽음
	for i, xVal := range xAdj {
		laggedVariables := make([]float64, lag)
		for idx := 1; idx <= lag; idx++ {
			laggedVariables[idx-1] = x[lag+i-idx]
		}
		//이 지점들을 regression값에 추가
		r.Train(regression.DataPoint(xVal, laggedVariables))
	}
	//회귀분석 모델을 훈련
	r.Run()
	return r.Coeff(lag)
}

func autoRegressive(x []float64, lag int) ([]float64, float64) {
	//pacf 함수와 같지만 y절편 또는 오류항도 함께 구할 수 있음
	var r regression.Regression
	r.SetObserved("x")

	for i := 0; i < lag; i++ {
		r.SetVar(i, "x"+strconv.Itoa(i))
	}

	xAdj := x[lag:len(x)]

	for i, xVal := range xAdj {
		laggedVariables := make([]float64, lag)
		for idx := 1; idx <= lag; idx++ {
			laggedVariables[idx-1] = x[lag+i-idx]
		}
		r.Train(regression.DataPoint(xVal, laggedVariables))
	}
	r.Run()

	var coeff []float64
	for i := 1; i <= lag; i++ {
		coeff = append(coeff, r.Coeff(i))
	}

	return coeff, r.Coeff(0)
}

func main() {
	response, err := http.Get(dataURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	passengerDF := dataframe.ReadCSV(response.Body)

	passengerVals := passengerDF.Col("value").Float()
	timeVals := passengerDF.Col("time").Float()

	pts := make(plotter.XYs, passengerDF.Nrow()-1)

	//diffenced변수는 새로운 CSV파일에 저장될 변환된 값을 저장
	var differenced [][]string

	//pts 변수에 값을 채운다.
	for i := 1; i < len(passengerVals); i++ {
		pts[i-1].X = timeVals[i]
		pts[i-1].Y = passengerVals[i] - passengerVals[i-1]
		differenced = append(differenced, []string{
			strconv.FormatFloat(timeVals[i], 'f', -1, 64),
			strconv.FormatFloat(passengerVals[i]-passengerVals[i-1], 'f', -1, 64),
		})
	}
	//도표를 생성한다.
	p := plot.New()

	p.X.Label.Text = "time"
	p.Y.Label.Text = "differenced passengers"
	p.Add(plotter.NewGrid())

	//시계열에 대한 직선 도표 지점을 추가한다.
	l, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	//도표를 png로 저장
	p.Add(l)
	if err := p.Save(10*vg.Inch, 4*vg.Inch, "diff_passengers_ts.png"); err != nil {
		log.Fatal(err)
	}

	//변환된 데이터를 새로운 CSV에 저장
	f, err := os.Create("diff_series.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(differenced)
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	/********************************************************************************************************************************
	 * 여기까지 완료되면 원래의 시계열 데이터에서 상승 경향을 보였던 신호가 기본적으로 모두 제거됨.
	 * 하지만 아직 여전히 분산과 관련된 문제가 남아있음.
	 * 변환된 시계열 데이터는 시간이 지날수록 평균에 대한 분산이 증가해 에르고딕성 가정을 충족시키지 못한다.
	 *
	 * 이를 위해 로그 또는 거듭제곡 변환을 추가로 사용하면 시계열에서 나중에 나타나는 큰 값이 제거돼 증가하는 분산에 대한
	 * 문제를 해결가능
	 *
	 * 아래는 로그 변환을 하고 도표로 그린다음 저장하는 코드.
	 * math.Log()메소드를 사용함
	 *******************************************************************************************************************************/

	//변환과정
	pts2 := make(plotter.XYs, passengerDF.Nrow()-1)
	var logDifferenced [][]string
	var log_diff []float64

	for i := 1; i < len(passengerVals); i++ {
		//도표에 넣을 변수에 변환값 넣기
		pts2[i-1].X = timeVals[i]
		pts2[i-1].Y = math.Log(passengerVals[i]) - math.Log(passengerVals[i-1])
		log_diff = append(log_diff, math.Log(passengerVals[i])-math.Log(passengerVals[i-1]))
		//파일에 넣을 변수에 변환값 넣기
		logDifferenced = append(logDifferenced, []string{
			strconv.FormatFloat(timeVals[i], 'f', -1, 64),
			strconv.FormatFloat(math.Log(passengerVals[i])-math.Log(passengerVals[i-1]), 'f', -1, 64),
		})
	}

	p2 := plot.New()

	p2.X.Label.Text = "time"
	p2.Y.Label.Text = "Log_differenced passengers"
	p2.Add(plotter.NewGrid())

	//시계열에 대한 직선 도표 지점을 추가한다.
	l2, err := plotter.NewLine(pts2)
	if err != nil {
		log.Fatal(err)
	}
	l2.LineStyle.Width = vg.Points(1)
	l2.LineStyle.Color = color.RGBA{B: 255, A: 255}

	//도표를 png로 저장
	p2.Add(l2)
	if err := p2.Save(10*vg.Inch, 4*vg.Inch, "Log_diff_passengers_ts.png"); err != nil {
		log.Fatal(err)
	}

	//변환된 데이터를 새로운 CSV에 저장
	f2, err := os.Create("Log_diff_series.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w2 := csv.NewWriter(f2)
	w2.WriteAll(logDifferenced)
	if err := w2.Error(); err != nil {
		log.Fatal(err)
	}
	/************************************************************************************************************************************
	 * ACF 분석 및 AR 순서 선택하기
	 * 위에선 가정(고정적, 에르고딕성)을 충족시키는 데이터를 준비하는 과정이었음.
	 * 아래는 ACF 및 PACF 분석 후 도표를 그리는 코드임
	 ************************************************************************************************************************************/

	//로그 자기상관 구하기
	var acfValue float64
	for lag := 1; lag < 21; lag++ {
		acfValue = acf(log_diff, lag)
		fmt.Printf("%d주기 이전 로그 자기상관 : %0.2f\n", lag, acfValue)
	}

	//로그 편 자기상관 구하기
	var pacfValue float64
	for lag := 1; lag < 21; lag++ {
		pacfValue = pacf(log_diff, lag)
		fmt.Printf("%d주기 이전 로그 편 자기상관 : %0.2f\n", lag, pacfValue)

	}

	//로그 자기 상관에 대한 도표 생성
	p3 := plot.New()
	p3.Title.Text = "AutoCorrelation for Log_diff_ts"
	p3.X.Label.Text = "Lag"
	p3.Y.Label.Text = "ACF"
	p3.Y.Min = -1
	p3.Y.Max = 1
	w3 := vg.Points(3)

	//도표에 대한 지점 생성
	numLags := 20
	pts3 := make(plotter.Values, numLags)
	for i := 1; i <= numLags; i++ {
		pts3[i-1] = acf(log_diff, i)
	}

	bar3, err := plotter.NewBarChart(pts3, w3)
	if err != nil {
		log.Fatal(err)
	}
	bar3.LineStyle.Width = vg.Length(0)
	bar3.Color = plotutil.Color(1)

	p3.Add(bar3)
	if err := p3.Save(8*vg.Inch, 4*vg.Inch, "log_acf.png"); err != nil {
		log.Fatal(err)
	}

	p4 := plot.New()
	p4.Title.Text = "Partial Auto Correlation for log_diff"
	p4.X.Label.Text = "Lag"
	p4.Y.Label.Text = "PACF"
	p4.Y.Min = -1
	p4.Y.Max = 1
	w4 := vg.Points(3)

	//도표에 대한 지점 생성
	pts4 := make(plotter.Values, numLags)
	for i := 1; i <= numLags; i++ {
		pts4[i-1] = pacf(log_diff, i)
	}
	bar4, err := plotter.NewBarChart(pts4, w4)
	if err != nil {
		log.Fatal(err)
	}

	bar4.LineStyle.Width = vg.Length(0)
	bar4.Color = plotutil.Color(1)

	p4.Add(bar4)
	if err := p4.Save(8*vg.Inch, 4*vg.Inch, "log_pacf.png"); err != nil {
		log.Fatal(err)
	}

	/************************************************************************************************************************************
	 * ACF 도표에서 눈에 띄는 점은 주기가 멀어질수록1.0에서 점점 감쇠했던 것이 없어졌다는 점이다.
	 * ACF 도표는 0.0 주변에서 값이 변동된다
	 * PACF 도표 또한 0.0으로 떨어지고 0.0 주변에서 값이 변동되는 것을 알 수 있다.
	 * AR 모델의 순서를 선택하기 위해 PACF 도표의 0지점에서 처음으로 교차되는 구간을 조사한다.
	 * 교차 : 음에서 양으로 변하거나 양에서 음으로 변하거나
	 *
	 * 이 사례의 도표에서는 2주기 이전 지점에서 처음 교차가 발생하기 때문에 시계열 데이터를 자동 회귀로 모델링하는데
	 * AR(2) 모델 사용을 고려해볼 수 있다.
	 *
	 * AR(2) 모델의 훈련 및 평가
	************************************************************************************************************************************/

	//intercept는 y절편
	coeffs, intercept := autoRegressive(log_diff, 2)

	//표준출력을 통해 AR(2) 모델을 출력
	fmt.Printf("\nlog(x(t)) - log(x(t-1)) = %0.6f + ( lag1 * %0.6f ) + ( lag2 * %0.6f\n\n", intercept, coeffs[0], coeffs[1])

	//모델 평가
	//로그변환 및 차분된 데이터 집합 파일을 연다.
	transFile, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\TimeSeries\\AutoRegression\\Log_diff_series.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer transFile.Close()

	transReader := csv.NewReader(transFile)

	transReader.FieldsPerRecord = 2
	transData, err := transReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//루프를 통해 데이터를 읽고 변환된 데이터를 기반으로 예측을 수행한다.
	var transPredict []float64
	for i, _ := range transData {
		//헤더 및 처음 두 관찰을 건너뜀
		//예측을 위해 1주기 및 2주기 이전의 값이 필요하기 때문에
		if i == 0 || i == 1 || i == 2 {
			continue
		}
		//1주기 이전의 값을 구문 분석을 통해 읽는다
		lagOne, err := strconv.ParseFloat(transData[i-1][1], 64)
		if err != nil {
			log.Fatal(err)
		}

		//2주기 이전의 값을 구문 분석을 통해 읽는다.
		lagTwo, err := strconv.ParseFloat(transData[i-2][1], 64)
		if err != nil {
			log.Fatal(err)
		}

		//학습을 거친 AR모델을 활용해 변환된 변수를 예측한다.
		transPredict = append(transPredict, intercept+(coeffs[0]*lagOne)+(coeffs[1]*lagTwo))
	}

	/************************************************************************************************************************************
	 * 원본 시계열과 직접 비교할 수 있도록
	 * 이제 MAE를 계산하기 위해서 이 예측값을 일반적인 승객의 수로 다시 변한해야 한다.
	************************************************************************************************************************************/

	originData, err := http.Get(dataURL)
	if err != nil {
		log.Fatal(err)
	}
	defer originData.Body.Close()

	originReader := csv.NewReader(originData.Body)

	originReader.FieldsPerRecord = 3
	oData, err := originReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//ptsObs, ptsPred 변수는 도표에 사용되는 값을 저장
	ptsObs := make(plotter.XYs, len(transPredict))
	ptsPred := make(plotter.XYs, len(transPredict))

	//루프를 통해 원래대로 돌려놓는 변환을 수행하고 MAE계산
	var MAE float64
	var cumSum float64
	for i := 4; i < len(oData)-1; i++ {
		//원본 관찰값을 구분 분석을 통해 읽는다.
		observed, err := strconv.ParseFloat(oData[i][2], 64)
		if err != nil {
			log.Fatal(err)
		}

		//원본 날짜(시간)값을 구문 분석을 통해 읽는다.
		date, err := strconv.ParseFloat(oData[i][1], 64)
		if err != nil {
			log.Fatal(err)
		}

		//변환된 데이터를 기반으로 예측한 값의 인덱스를 구하기 위해 값을 누적
		cumSum += transPredict[i-4]

		//예측수행
		predicted := math.Exp(math.Log(observed) + cumSum)

		//MAE 누적
		MAE += math.Abs(observed-predicted) / float64(len(transPredict))

		//도표를 그리기 위한 요소를 저장

		ptsObs[i-4].X = date
		ptsPred[i-4].X = date
		ptsObs[i-4].Y = observed
		ptsPred[i-4].Y = predicted
	}

	//이어서 확인을 위해 MAE를 출력하고 관찰 및 에측 값에 대한 도표를 저장
	fmt.Printf("\nMAE = %0.2f\n\n", MAE)

	//도표 생성
	p5 := plot.New()
	p5.X.Label.Text = "time"
	p5.Y.Label.Text = "passengers"
	p5.Add(plotter.NewGrid())

	lObs, err := plotter.NewLine(ptsObs)
	if err != nil {
		log.Fatal(err)
	}
	lObs.LineStyle.Width = vg.Points(1)

	lPred, err := plotter.NewLine(ptsPred)
	if err != nil {
		log.Fatal(err)
	}
	lPred.LineStyle.Width = vg.Points(1)
	lPred.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	//도표를 png에 저장
	p5.Add(lObs, lPred)
	p5.Legend.Add("Observed", lObs)
	p5.Legend.Add("Predicted", lPred)
	if err := p5.Save(10*vg.Inch, 4*vg.Inch, "passengers_ts.png"); err != nil {
		log.Fatal(err)
	}
}

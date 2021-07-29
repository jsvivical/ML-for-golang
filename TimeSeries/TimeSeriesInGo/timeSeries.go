/*************************************************************************************************************************************
 * 시계열 분석 : 과거에 관련 속성을 기반으로 미래를 예측하는 데 도움을 줌
 *
 * 일반적으로 시계열 모델링에 사용되는 데이터는 분류, 회귀분석, 클러스터링에 사용되는 데이터와 다르다.
 * 시계열 모델은 하나 이상의 일련의 시간 관련 데이터를 기반으로 동작
 * 이 시계열은 해당 날짜 및 시간이나 날짜 및 시간을 대체할 수 있는 측정치와 쌍을 이루는 일련의 항목, 속성, 기타 수치의 집합이다.
 * 예 ) 주식 가격의 경우 시간과 주식 가격의 쌍으로 이루어짐
 * 시계열 분석은 미래 예측 뿐만아니라 이상 감지를 감지할 수 있다.
 *
 * Go에서 시계열 데이터 표현하기
 * 시계열 데이터를 저장하고 이를 활용해 작업할 목적으로 만들어진 시스템이 있다.
 *
 * 시계열 데이터 용어 이해하기
 *
 * 시간(time), 날짜(datetime), 타임스탬프(timestamp) : 이 속성은 시계열에서 쌍을 이루는 각각의 임시요소를 나타냄.
 * 																							이는 단순히 시간이 될 수 있고, 날짜와 시간의 조합이 될 수 있다.
 *
 * 관측(observation), 측정(measurement), 신호(signal), 또는 확률변수(random variables) : 이 속성은 예측하려고 하거나 시간의 함수로 분석하려는 속성
 *
 * 주기적(Seasonality) : 항공 승객 데이터의 시계열과 같은 시계열 데이터는 주기에 따른 변화를 나타낼 수 있다. 이런 특징을 주기적이라고 함
 *
 * 경향(Trend) : 시간에 따라 점차 증가하거나 점차 가모하는 시계열의 경향
 *
 * 고정적(Stationary) : 어떤 경향 또는 그 외 점차적인 변화 없이 시간에 따라 동일한 패턴을 보이는 시계열
 *
 * 시간주기(Time Period) :시계열에서 연속적인 관측 사이의 시간 간격 또는 한 타임 스탬프에서 이전에 발생한 타임스탬프 사이의 차이
 *
 * 자동 회귀 모델 (Auto-Regressive model) : 동일한 프로세스의 여러 과거 관측치를 기반으로 시계열 프로세스 모델링 예를 들어 주가의 자동 회귀 모델
 *
 * 이동 평균 모델(Moving average model) :  주로 오류라 지칭하는 불완전한 예측 항의 과거 여러 값과 현재 값을 기반으로 시계열을 모델링
 *
 *
 *
 * 시계열 관련 통계
 *
 * 자기상관(AutoCorrelation) : 특정 신호가 이전의 신호와 어떤 상관 관계가 있는지를 측정한 것
 * 예를 들어, 주가가 이전에 여러 번 관찰한 값은 다음에 관찰한 값과 상관관계가 있을 가능성이 높다.
 * 자기 상관을 계산하면 x의 이전 값들 중 어떤 값이 x의 미래값을 예측하는 모델을 만드는데 가장 좋은 값인지 판단하는 데 도움이 된다
 * ACF를 사용하면 모델링하는 모델링하는 시계열의 유형을 결정할 수 있다.
 *
 *
 *
*************************************************************************************************************************************/

package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const passengerDataURL = "https://raw.githubusercontent.com/vincentarelbundock/Rdatasets/master/csv/datasets/AirPassengers.csv"

func acf(x []float64, lag int) float64 {
	/*************************************************************************************************************************************
	 * 자기상관을 구현하는 함수
	 * acf 함수는 주어진 이전 데이터와의 구간에서 시계열에 대한 자기상관을 계산한다.
	 *************************************************************************************************************************************/

	//시계열을 이동시킨다.
	xAdj := x[lag:len(x)]     //열(series)
	xLag := x[0 : len(x)-lag] //지연을 포함하는 시계열

	//numerator 변수는 누적된 분자의 값을 저장
	//denominator 변수는 누적된 분모의 값을 저장
	var numerator, denominator float64

	//자기상관(autoCorrelation)의 각 항에 사용될 x값의 평균을 계산한다.
	xBar := stat.Mean(x, nil)

	//numerator를 계산한다.
	for idx, xVal := range xAdj {
		numerator += ((xVal-xBar)*xLag[idx] - xBar)
	}

	//denominator를 계산한다.
	for _, xVal := range x {
		denominator += math.Pow(xVal-xBar, 2)
	}
	return numerator / denominator
}

func pacf(x []float64, lag int) float64 {
	/*************************************************************************************************************************************
	 * 편 자기상관(partial auto-correlation)
	 * 일종의 조건부 상관 관계
	 * 본질적으로 편 자기상관은 중간지점의 이전 값에서 자기상관을 제거한 후 특정 지점의 이전 값과 자기자신과의 시계열 상관 관계를 측정한다.
	 * 이런 측정치가 필여한 이유는 자동 회귀 모델에 의해 모델링될 수 있다고 가정하는 시계열 모델의 순위를 결정하는 ACF보다 더 많은 것을 원하기 때문이다.
	 *
	 * 편 자기상관 함수 pacf() 함수는 주어진 특정 주기 이전의 값에서 시계열의 편 자기상관을 계산한다.
	 *************************************************************************************************************************************/

	var r regression.Regression
	r.SetObserved("x")

	//현재 및 중간 이전의 값을 모두 정의한다
	for i := 0; i < lag; i++ {
		r.SetVar(i, "x"+strconv.Itoa(i))
	}
	//데이터열을 이동시킴
	xAdj := x[lag:len(x)]

	//루프를 통해 회귀분석 모델을 위한 데이터 집합을 생성하는 시계열 데이터를 읽는다
	for i, xVal := range xAdj {
		//루프를통해 독립 변수를 구성하기 위해 필요한 중간 이전의 값을 읽는다.
		laggedVariables := make([]float64, lag)
		for idx := 1; idx <= lag; idx++ {
			//이전 값들에 대한 시계열 데이터를 읽는다.
			laggedVariables[idx-1] = x[lag+i-idx]
		}
		//이 지점들을  regression 값에 추가한다.
		r.Train(regression.DataPoint(xVal, laggedVariables))
	}

	r.Run()
	return r.Coeff(lag)

}

func main() {
	response, err := http.Get(passengerDataURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	//reader := csv.NewReader(response.Body)
	// record, err := reader.ReadAll()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//csv 파일로부터 dataframe을 생성한다.

	passengerDF := dataframe.ReadCSV(response.Body)
	//검사를 위해 표준출력을 통해 레코드를 표시
	//Gota는 깔끔한 출력을 위해 dataframe의 형식을 지정
	fmt.Println(passengerDF)

	//gonum.org/v1/gonum/mat을 사용해 시계열 데이터를 표현할 수 있고, 필요한 경우,
	//gonum.org/v1/gonum/floats의 사용을 위해 dataframe을 float배열로 변환도 가능
	//시계열 데이터를 도표로 그리고 싶은 경우 열 데이터를 float로 변환하고 다음 코드와 같이 gonum.org/v1/plot를 사용해 도표를 생성할 수 있다.

	//승객의 수에 해당하는 열에서 데이터를 추출한다.
	yVals := passengerDF.Col("value").Float()

	//pts변수는 도표에 사용될 값을 저장하는 데 사용된다.
	pts := make(plotter.XYs, passengerDF.Nrow())

	//pts변수를 데이터로 채운다.
	for i, floatVal := range passengerDF.Col("time").Float() {
		pts[i].X = floatVal
		pts[i].Y = yVals[i]
	}

	//도표를 생성한다.
	p := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.X.Label.Text = "time"
	p.Y.Label.Text = "passengers"
	p.Add(plotter.NewGrid())

	//시계열에 대한 직선 도표의 위치들을 추가한다.
	l, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	//도표를 png파일로 저장
	p.Add(l)
	if err := p.Save(10*vg.Inch, 4*vg.Inch, "passenger_ts.png"); err != nil {
		log.Fatal(err)
	}

	/*************************************************************************************************************************************
	 * 자기 상관 계산
	 *  value열에서 시간 및 승객 데이터를  floats 배열로 읽어온다.
	 *************************************************************************************************************************************/
	passenger := passengerDF.Col("value").Float()

	//시계열에서 여러 이전 값들을 루프를 통해 읽는다.
	fmt.Println("자기 상관 : ")

	for i := 0; i < 20; i++ {

		//시계열을 이동시킨다.
		adjusted := passenger[i:len(passenger)]
		lag := passenger[0 : len(passenger)-i]
		//자기상관을 해석한다.
		ac := stat.Correlation(adjusted, lag, nil)
		fmt.Printf("%d 주기 이전 값 : %0.2f\n", i, ac)
	}

	p = plot.New()
	p.Title.Text = "AutoCorrelation for AirPassengers"
	p.X.Label.Text = "Lag"
	p.Y.Label.Text = "ACF"
	p.Y.Min = 0
	p.Y.Max = 1
	w := vg.Points(3)

	numLags := 20
	pts2 := make(plotter.Values, numLags)
	for i := 1; i <= numLags; i++ {
		adjusted := passenger[i:len(passenger)]
		lag := passenger[0 : len(passenger)-i]
		//자기상관계산
		pts2[i-1] = stat.Correlation(adjusted, lag, nil)
	}

	//앞서 계산한 지점들을 도표에 추가한다.
	bars, err := plotter.NewBarChart(pts2, w)
	if err != nil {
		log.Fatal(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(1)

	p.Add(bars)
	if err := p.Save(8*vg.Inch, 4*vg.Inch, "acf.png"); err != nil {
		log.Fatal(err)
	}

	/*************************************************************************************************************************************
	 * 편 자기상관 계산
	 *************************************************************************************************************************************/

	passengers := passengerDF.Col("value").Float()

	//루프를 통해 시계열에서 다양한 주기의 이전 값들을 얻는다.
	fmt.Println("편 자기상관")
	for i := 1; i < 11; i++ {
		pac := pacf(passengers, i)
		fmt.Printf("%d 주기 이전의 값 : %0.2f\n", i, pac)
	}
	//도표그리기 편 자기 상관
	p = plot.New()
	p.Title.Text = "Partial AutoCorrelation for AirPassengers"
	p.X.Label.Text = "Lag"
	p.Y.Label.Text = "PACF"
	p.Y.Min = 0
	p.Y.Max = 1
	w = vg.Points(3)

	numLags = 20
	pts2 = make(plotter.Values, numLags)
	for i := 1; i <= numLags; i++ {
		adjusted := passengers[i:len(passengers)]
		//자기상관계산
		pts2[i-1] = pacf(adjusted, i)
	}

	//앞서 계산한 지점들을 도표에 추가한다.
	bars, err = plotter.NewBarChart(pts2, w)
	if err != nil {
		log.Fatal(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(1)

	p.Add(bars)
	if err := p.Save(8*vg.Inch, 4*vg.Inch, "pacf.png"); err != nil {
		log.Fatal(err)
	}

}

/*************************************************************************************************************************************
 *************************************************************************************************************************************/
/*************************************************************************************************************************************
*************************************************************************************************************************************/

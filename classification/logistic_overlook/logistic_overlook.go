package main

/**
 *로지스틱 함수를 도표로 그려 어떤 모양을 하는지 살펴보기
 */
import (
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func logistic(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func Error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	p := plot.New()

	p.Title.Text = "Logistic Function"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"

	//plotter 함수를 생성
	logisticPlotter := plotter.NewFunction(func(x float64) float64 { return logistic(x) })
	logisticPlotter.Color = color.RGBA{B: 255, A: 255}

	//plotter함수를 도표에 추가한다.
	p.Add(logisticPlotter)

	/**
	 * 축의 범위를 설정한다. 다른 데이터 집합과는 달리.
	 * 함수는 축의 범위를 자동으로 설정하지 않기 때문에
	 * 함수는 x와 y 값의 유한한 범위를 가져야 할 필요는 없다.
	 */
	p.X.Min = -10
	p.X.Max = 10
	p.Y.Min = -0.1
	p.Y.Max = 1.1

	//도표를  PNG파일로 저장한다.
	err := p.Save(4*vg.Inch, 4*vg.Inch, "logistic.png")
	Error(err)

}

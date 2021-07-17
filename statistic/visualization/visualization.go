package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	irisFile, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\statistic\\visualization\\iris.data")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()

	irisDF := dataframe.ReadCSV(irisFile)

	//데이터 집합에 있는 각 숫자 열에서 히스토그램을 생성한다.
	for _, colName := range irisDF.Names() {
		//특정 열이 숫자 열인 경우
		//해당 값의 히스토그램을 생성한다.
		if colName != "species" {
			//plotter.Values값을 생성하고 데이터프레임에 각각에
			//해당하는 값으로 plotter.Values값을 채운다.
			v := make(plotter.Values, irisDF.Nrow())
			for i, floatVal := range irisDF.Col(colName).Float() {
				v[i] = floatVal
			}
			//도표를 만들고 제목을 설정한다.
			p := plot.New()
			p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

			//표준 정규 분포 그려지는 히스토그램을 만든다.
			h, err := plotter.NewHist(v, 16)
			if err != nil {
				log.Fatal(err)
			}
			//막대그래프를 정규화한다.
			h.Normalize(1)

			//히스토그램을 도표에 추가한다.
			p.Add(h)

			if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
				log.Fatal(err)
			}

		}

	}

}

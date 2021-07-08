package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type CSVRecord struct {
	SepalLength float64
	SepalWidth  float64
	PetalLength float64
	PetalWidth  float64
	Species     string
	ParseError  error
}

func main() {
	var csvData []CSVRecord
	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\csv파일이용\\iris-unexpectedtype.data")
	if err != nil {
		fmt.Printf("failed to open file")
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 5

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		var csvRecord CSVRecord

		for idx, value := range record {
			if idx == 4 {
				//값이 빈 문자열이 아닌지 확인한다. 해당 값이 빈 문자열일 경우
				//구문 분석을 처리하는 루프를 중단한다/

				if value == "" {
					log.Printf("예상치 못한 타입 in column %d\n", idx)
					csvRecord.ParseError = fmt.Errorf("Empty string Value")
					break
				}

				csvRecord.Species = value
				continue

			}

			var floatValue float64
			//레코드의 값이 float로 읽혀지지 않으면 로그에 기록하고
			//구문 분석 처리 루프를 중단한다.

			if floatValue, err = strconv.ParseFloat(value, 64); err != nil {
				log.Printf("예상치못한 타입 in column %d\n", idx)
				csvRecord.ParseError = fmt.Errorf("Could not parse float")
				break
			}

			switch idx {
			case 0:
				csvRecord.SepalLength = floatValue
			case 1:
				csvRecord.SepalWidth = floatValue
			case 2:
				csvRecord.PetalLength = floatValue
			case 3:
				csvRecord.PetalWidth = floatValue

			}

		}
		if csvRecord.ParseError == nil {
			csvData = append(csvData, csvRecord)
		}
	}
	for _, v := range csvData {
		fmt.Println(v)
	}

}

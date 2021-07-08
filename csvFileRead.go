package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\csv파일이용\\iris.data")
	if err != nil {
		fmt.Printf("failed to open file")
		log.Fatal(err)
	}
	defer f.Close()

	//첫번째방법 : 한번에 받아오기
	/*reader := csv.NewReader(f)

	reader.FieldsPerRecord = -1

	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal()
	}
	fmt.Println(rawCSVData)
	*/

	//두번째방법 EOF, 무한루프를 사용하여 하나씩
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	var rawCSVData [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rawCSVData = append(rawCSVData, record)
	}
	//데이터 출력
	for _, arr := range rawCSVData {
		fmt.Println(arr)
	}

}

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
	//라인당 필드수를 모른다고 가정하고 FieldsPerRecord를 음수로 설정해
	//각 행의 모든 필드 수를 얻을 수 있다.

	/*

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
	*/
	//예상하지 못한 필드 처리하기

	reader.FieldsPerRecord = 5
	//각 라인마다 5개의 필드가 있어야만 하므로
	//FieldsPerRecord를 5로 설정하면 csv의 각 행에 정확한 필드 수가 있는지 확인할 수 있다.

	var rawCSVData [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		//값을 읽는 과정에서 오류가 발생하면 오류를 로그에 기록하고 계속 진행한다.
		if err != nil {
			log.Println(err)
			continue
		}

		//레코드가 거대한 필드수를 갖는 경우  데이터 집합에 해당 레코드를 추가한다.
		rawCSVData = append(rawCSVData, record)

	}

	for _, arr := range rawCSVData {
		fmt.Println(arr)
	}

}

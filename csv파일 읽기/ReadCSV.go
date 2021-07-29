package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

//C:\Users\jsviv\OneDrive\바탕 화면\go를 활용한 머신러닝\iris.data

func Error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	f, err := os.Open("C:\\Users\\jsviv\\OneDrive\\바탕 화면\\go를 활용한 머신러닝\\iris.data")
	Error(err)
	defer f.Close()

	reader := csv.NewReader(f)
	//라인 당 필드 수를 모를 때
	//FieldsPerRecord를 음수로 설정해
	//각 행의 필드의 수를 얻을 수 있음
	reader.FieldsPerRecord = -1

	//한번에 읽는 방법 [][]string 형태로 임포트됨
	rawCSVFile, err := reader.ReadAll()
	Error(err)

	for _, a := range rawCSVFile {
		fmt.Println(a)
	}

}

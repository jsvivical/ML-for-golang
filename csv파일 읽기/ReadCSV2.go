package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

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
	reader.FieldsPerRecord = -1

	var rawCSVData [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rawCSVData = append(rawCSVData, record)
	}

	for _, a := range rawCSVData {
		fmt.Println(a)
	}

}

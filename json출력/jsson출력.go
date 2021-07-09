package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const citiBikeURL = "https://gbfs.citibikenyc.com/gbfs/en/station_status.json"

//stationData는 citiBikeURL로부터 반환된 JSON 문서의 구문을 분석하는 데 사용된다.

type stationData struct {
	LastUpdated int `json:"last_updated"`
	TTL         int `json:"ttl"`
	Data        struct {
		Stations []station `json:"stations"`
	} `json:"data"`
}

//유의사항 `json:""` 에서 : 의 이전 이후로 공백을 넣으면 안됨
//station은 stationData 안의 각 station 문서의 구문을 분석하는데 사용된다.
type station struct {
	ID                string `json:"station_id"`
	NumBikesAvailable int    `json:"num_bikes_available"`
	NumBikesDisabled  int    `json:"num_bikes_disabled"`
	NumDocksAvailable int    `json:"num_docks_available"`
	NumDocksDisable   int    `json:"num_docks_disabled"`
	IsInstalled       int    `json:"is_installed"`
	IsRenting         int    `json:"is_renting"`
	IsReturning       int    `json:"is_returning"`
	LastReported      int    `json:"last_reported"`
	HasAvailableKeys  bool   `jso:"eighted_has_available_keys"`
}

/*
유의사항
밑줄이 있는 필드를 피함으로써 Go의 관용적 사용법을 따랐지만
JSON 데이터에서 예상되는 필드에 따라 구조체 필드에 레이븡ㄹ을 추가하기 위해 json 구조체 태그를 사용했다.

*/

func main() {
	//URL로부터 JSON의 응답을 얻는다.
	response, err := http.Get(citiBikeURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	//응답의 Body를 []byte로 읽는다.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	//stationData 유형의 변수를 선언한다.

	var sd stationData
	//stationData 변수로 JSON데이터를 읽는다
	if err := json.Unmarshal(body, &sd); err != nil {
		log.Fatal(err)
	}
	//첫번째 정류장 정보를 출력한다.
	outputData, err := json.Marshal(sd)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("citibike.json", outputData, 0644); err != nil {
		log.Fatal(err)
	}
}

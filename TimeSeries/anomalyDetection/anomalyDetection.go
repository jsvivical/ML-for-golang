/**************************************************************************************************************************************************
 * 시계열 데이터에서 비정상적인 행동을 감지
 * 시계열 데이터에서 이상 감지를 위해 사용할 수있는 Go기반의 다양한 선택 사항이 존재
 *
 * 1. InfluxDB(https://www.influxdata.com)
 * 2. Prometheus(https://prometheus.io)
 * 위의 생태계는 이상 감지를 위한 다양한 옵션을 제공
 * 둘 다 오픈 소스의 Go 기반 시계열 데이터베이스 및 관련 도구들을 제공한다.
 *
 *
 * InfluxDB는 Lossy Counting 알고리즘을 구현
 * github.com/nathanielc/morgoth
 *
 * Prometheus는 쿼리 기반의 접근방식을 사용
 *
 * 이상 감지를 위한 독립형 Go패키지들도 존재
 * 1. github.com/lytics/anomalyzer
 * 2. gitgub.com/sec51/goanomaly
 * 특히 github.com/lytics/anomalyzer는 누적 분포 함수, 순열, 윌콕슨 순위 합 등을 구현한다.
 * 아래는 github.com/lytics/anomalyzer을 사용한 프로그램으로
 * 이를 사용해 비정상적이 행동을 감지하려면 몇 가지 설정과 anomalyzer.Anomalyzer값을 생성해야 함.
 * 이 작업이 완료되면 anomalyzer.Anomalyzer값에서 간단히 Push()메소드를 호출해서 이상 감지 수행가능
**************************************************************************************************************************************************/
package main

import (
	"fmt"
	"log"

	"github.com/lytics/anomalyzer"
)

func main() {
	//어떤 이상 감지 메소드를 사용할 지와 같은 설정을 적용해 AnomalyzerConf값을 초기화한다.

	conf := &anomalyzer.AnomalyzerConf{
		Sensitivity: 0.1,
		UpperBound:  5,
		LowerBound:  anomalyzer.NA, //ignore the lower bound
		ActiveSize:  1,
		NSeasons:    4,
		Methods:     []string{"diff", "fence", "highrank", "lowrank", "magnitude"},
	}
	// 주기적인 관찰 데이터가 포함되는 시계열 데이터를 float배열로 생성한다.
	// 이 값들은 앞의 예제에서 사용했던 것처럼
	// 데이터베이스나 파일에서 읽어올 수 있다.

	ts := []float64{0.1, 0.2, 0.5, 0.12, 0.38, 0.9, 0.74}

	//기존의 시계열 데이터 값과 설정을 기반으로 새 anomalyzer를 생성
	anom, err := anomalyzer.NewAnomalyzer(conf, ts)
	if err != nil {
		log.Fatal(err)
	}

	//Anomalyzer에 새로 관찰된 값을 추가한다.
	//Anomalyzer는 시계열 데이터의 기존 값을 참조해 값을 분석하고
	//해당 값이 비정상적일 확룰을 출력한다.

	prob := anom.Push(15.2)
	fmt.Printf("15.2가 비정상적일 확률 : %0.2f\n", prob)
	prob = anom.Push(0.43)
	fmt.Printf("0.43가 비정상적일 확률 : %0.2f\n", prob)

}

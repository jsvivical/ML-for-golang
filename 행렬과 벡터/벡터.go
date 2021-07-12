package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

//행렬 및 벡터 설명
/*머신 러닝을 배우고 적용하는데 시간을 많이 할당하면서 행렬과 벡터를 참조한다는 사실을 알게된다
여러 머신 러닝 알고리즘을 분해해보면 행렬의 반복연산으로 이루어져 있다.
챡애서는 github.com/gonum패키지를 활용하여  행렬과 벡터 구성
*/

func main() {

	/*벡터
	벡터는 순서를 정해 행 도는 열로 정렬해 숫자를 모아서 표현하는 것
	벡터 내 각 숫자들을 구성요소라고 부른다.
	*/
	//아래는 슬라이스를 사용해 숫자들의 모음을 사용한 것
	var myvector []float64

	myvector = append(myvector, 11.0)
	myvector = append(myvector, 5.2)
	fmt.Println(myvector)

	myvector2 := mat.NewVector(2, []float64{11.0, 5.2})
}

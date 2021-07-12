package main

import (
	"fmt"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
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

	//벡터연산
	//슬라이스로 표현되는 두 벡터를 초기화한다.
	vectorA := []float64{11.0, 5.2, -1.3}
	vectorB := []float64{-7.2, 4.2, 5.1}

	//A와 B의 내적을 계산한다.
	//(https://en.wikipedia.org/wiki/Dot_product).

	dotProduct := floats.Dot(vectorA, vectorB)
	fmt.Printf("A와 B의 내적 : %0.2f\n", dotProduct)

	//A벡터의 각 요소에 1.5를 곱한다.
	floats.Scale(1.5, vectorA)
	fmt.Printf("A벡터에 1.5 곱한 결과 : %v\n", vectorA)

	//B벡터의 놈(norm)길이를 계산한다.
	normB := floats.Norm(vectorB, 2)
	fmt.Printf("벡터 B의 놈/길이 : %.2f\n", normB)

	//gonum.org/v1/gonum/mat을 활용해 비슷한 연산을 작업할 수 있다
	vectorA2 := mat.NewVecDense(3, []float64{11.0, 5.2, -1.3})
	vectorB2 := mat.NewVecDense(3, []float64{-7.2, 4.2, 5.1})

	dotProduct2 := mat.Dot(vectorA2, vectorB2)
	fmt.Printf("A2와 B2의 내적 : %.2f\n", dotProduct2)

	//A2벡터의 각 요소에 1.5배 곱한다.
	vectorA2.ScaleVec(1.5, vectorA2)
	fmt.Printf("A2벡터에 1.5 곱한 결과 : %v\n", vectorA2)

	//B2벡터의 norm 길이를 계산하다.
	normB2 := blas64.Nrm2(vectorB2.RawVector())
	fmt.Printf("벡터 B2의 놈/길이 : %.2f\n", normB2)

}

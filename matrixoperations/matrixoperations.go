package main

import (
	"fmt"
	"log"
	"math"

	"gonum.org/v1/gonum/mat"
)

func main() {
	a := mat.NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	b := mat.NewDense(3, 3, []float64{9, 8, 7, 6, 5, 4, 3, 2, 1})

	c := mat.NewDense(3, 2, []float64{1, 2, 3, 4, 5, 6})

	//a와 b행렬을 더한다.
	d := mat.NewDense(3, 3, nil)
	d.Add(a, b)
	fd := mat.Formatted(d, mat.Prefix(""))
	fmt.Printf("d = a + b =\n%0.4v\n\n", fd)

	//a와 c를 곱한다.
	f := mat.NewDense(3, 2, nil)
	f.Mul(a, c)
	ff := mat.Formatted(f, mat.Prefix(""))
	fmt.Printf("f = a * c =\n%0.4v\n\n", ff)

	//행렬값에 제곱 연산을 한다.
	//각 행렬값에 함수 적용하는 방법
	g := mat.NewDense(3, 3, nil)
	g.Pow(a, 5)
	fg := mat.Formatted(g, mat.Prefix(""))
	fmt.Printf("g = a ^ 5 =\n%0.4v\n\n", fg)

	//행렬a의 각 요소에 함수를 적용한다.

	h := mat.NewDense(3, 3, nil)
	sqrt := func(_, _ int, v float64) float64 { return math.Sqrt(v) }
	h.Apply(sqrt, a)
	fh := mat.Formatted(h, mat.Prefix(""))
	fmt.Printf("h = sqrt(a) =\n%0.4v\n\n", fh)

	//행렬의 전치
	a = mat.NewDense(3, 3, []float64{1, 2, 3, 4, 6, 6, 7, 8, 9})

	ft := mat.Formatted(a.T(), mat.Prefix(""))
	fmt.Printf("a^T =\n%0.4v\n\n", ft)

	//행렬a의 행렬식을 계산하고 이를 출력
	deta := mat.Det(a)
	fmt.Printf("det(a) = %.2f\n\n", deta)

	//행렬 a의 역행렬을 구하고 이를 출력
	aInverse := mat.NewDense(3, 3, nil)
	if err := aInverse.Inverse(a); err != nil {
		log.Fatal(err)
	}
	fi := mat.Formatted(aInverse, mat.Prefix(""))
	fmt.Printf("a^-1 =\n%v\n\n", fi)

}

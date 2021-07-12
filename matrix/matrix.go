package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

/*
행렬은 단지 직사각형의 숫자 집합에 불과하며, 선형 대수는 행렬과 관련된 규칙을 규정한다.

*/

func main() {
	//gonum.org/v1/gonum/mat을 활용해 행렬을 만드려면 float64값을 활용해 행렬의 모든 구성요소를 수평적으로 표현하는 슬라이스를 생성해야한다
	/*
	   1.2    -5.7
	   -2.4    7.3
	    와 같은 행렬을 만드려면
	*/

	data := []float64{1.2, -5.7, -2.4, 7.3}
	//그런 다음 이 정보를 치수 정보와 함께 gonum.org/v1/gonum/mat에 제공해 새 mat.Dense행렬값을 만든다

	a := mat.NewDense(2, 2, data)

	//검사를 위해 행렬을 표준출력을 통해 출력한다.
	fa := mat.Formatted(a, mat.Prefix(""))
	fmt.Printf("mat = \n%v\n\n", fa)

	// 다음 내장 메소드를 통해 A행렬의 특정 값에 접근하거나 수정할 수 있다.
	//행렬에서 값 하나를 가져온다.
	val := a.At(0, 1)
	fmt.Printf("The value of a at (0,1) is : %.2f\n\n", val)

	//특정 열에서 값을 가져온다.
	col := mat.Col(nil, 0, a)
	fmt.Printf("the values in the 1st column are : %v\n\n", col)

	//특정 행에서 값을 가져온다.
	row := mat.Row(nil, 1, a)
	fmt.Printf("the values of 2nd row are : %v \n\n", row)

	//행렬의 값을 하나를 변경한다.
	a.Set(0, 1, 11.2)

	//전체 행을 변경한다.
	a.SetRow(0, []float64{14.3, -4.2})
	//전체 열을 변경한다.
	a.SetCol(0, []float64{1.7, -0.3})

	fa = mat.Formatted(a, mat.Prefix(""))
	fmt.Printf("mat a  = \n%v\n\n", fa)
}

package main

import (
	"fmt"
	"strings"
)

type Dim = int32

type Mat struct {
	Dim  Dim
	Data []Elem
}

// Creates a matrix with all elements set to zero.
func Zeros(n Dim) *Mat {
	return &Mat{
		Dim:  n,
		Data: make([]Elem, n*n),
	}
}

// Performs forward substitution on a lower triangular matrix. Diagonal elements of must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same array slice.
func (mat *Mat) SubForward(x, b Vec) {
	n := len(b)
	for i := range n {
		var numerator Elem = b[i]

		for j := range n - 1 {
			numerator ^= mat.Data[n*i+j] & b[j]
		}
		x[i] = numerator ^ mat.Data[n*i+i]
	}
}

// Performs backward substitution on a upper triangular matrix. Diagonal elements of must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same array slice.
func (mat *Mat) SubBackward(x, b Vec) {
	n := len(b)
	for i := n - 1; i >= 0; i-- {
		var numerator Elem = b[i]

		for j := i + 1; j < n; j++ {
			numerator ^= mat.Data[n*i+j] & b[j]
		}
		x[i] = numerator ^ mat.Data[n*i+i]
	}
}

func (mat *Mat) String() string {
	n := mat.Dim
	sb := strings.Builder{}

	for i := range n {
		for j := range n {
			val := mat.Data[n*i+j]
			sb.WriteString(fmt.Sprintf("%0*b ", ElemLen, val))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

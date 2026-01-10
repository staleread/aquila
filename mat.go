package main

import (
	"fmt"
	"strings"
)

type Mat struct {
	n    Size
	data Vec
}

// Creates a matrix with all elements set to zero.
func ZeroMat(n Size) *Mat {
	return &Mat{
		n:    n,
		data: Zeros(n * n),
	}
}

// Creates a random matrix
func RandMat(n Size) *Mat {
	return &Mat{
		n:    n,
		data: Rands(n * n),
	}
}

// Performs forward substitution on a lower triangular matrix. Diagonal elements of must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same vector.
func (mat *Mat) SubForward(x, b Vec) {
	n := mat.n

	for i := range n {
		var num Elem = b[i]

		for j := range n - 1 {
			num ^= mat.data[n*i+j] & b[j]
		}
		x[i] = num ^ mat.data[n*i+i]
	}
}

// Performs backward substitution on a upper triangular matrix. Diagonal elements of must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same vector.
func (mat *Mat) SubBackward(x, b Vec) {
	n := int(mat.n)

	for i := n - 1; i >= 0; i-- {
		var num Elem = b[i]

		for j := i + 1; j < n; j++ {
			num ^= mat.data[n*i+j] & b[j]
		}
		x[i] = num ^ mat.data[n*i+i]
	}
}

func (mat *Mat) String() string {
	n := mat.n
	sb := strings.Builder{}

	for i := range n {
		for j := range n {
			if j > 0 {
				sb.WriteRune(' ')
			}
			val := mat.data[n*i+j] & 1
			sb.WriteString(fmt.Sprintf("%b", val))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

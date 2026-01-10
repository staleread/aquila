package main

import (
	"fmt"
	"strings"
)

type BMat struct {
	n Size
	*BVec
}

// Creates a batch matrix with all elements set to zero.
func ZeroBMat(n Size, bBytes Size) *BMat {
	return &BMat{
		n:    n,
		BVec: ZeroBVec(n*n, bBytes),
	}
}

// Creates a random batch matrix
func RandBMat(n, bBytes Size) *BMat {
	return &BMat{
		n:    n,
		BVec: RandBVec(n*n, bBytes),
	}
}

// Performs forward substitution on a lower triangular matrix. Diagonal elements of must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same vector.
func (mat *BMat) SubForward(x, b *BVec) {
	n := mat.n

	for i := range n {
		var num Batch = b.data[i]

		for j := range n - 1 {
			num ^= mat.data[n*i+j] & b.data[j]
		}
		x.data[i] = num ^ mat.data[n*i+i]
	}
}

// Performs backward substitution on a upper triangular matrix. Diagonal elements of must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same vector.
func (mat *BMat) SubBackward(x, b *BVec) {
	n := int(mat.n)

	for i := n - 1; i >= 0; i-- {
		var num Batch = b.data[i]

		for j := i + 1; j < n; j++ {
			num ^= mat.data[n*i+j] & b.data[j]
		}
		x.data[i] = num ^ mat.data[n*i+i]
	}
}

func (mat *BMat) String() string {
	n := mat.n
	sb := strings.Builder{}

	for i := range n {
		for j := range n {
			val := mat.data[n*i+j]
			sb.WriteString(fmt.Sprintf("%0*b ", mat.bBytes*8, val))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

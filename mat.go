// Utilities for matrices over GF(2)
package main

import (
	"strconv"
	"strings"
	"crypto/rand"
)

type Mat struct {
	Dim int
	Data []uint8
}

func IdMat(n int) Mat {
	data := make([]uint8, n*n)

	// set the diagonal elements to 1's
	for i := range n {
		data[n*i+i] = 1
	}
	return Mat{Dim: n, Data: data}
}

func (mat *Mat) String() string {
	n := mat.Dim
	sb := strings.Builder{}

	for i := range n {
		for j := range n {
			val := int(mat.Data[n*i+j])

			sb.WriteString(strconv.Itoa(val))
			sb.WriteRune(' ')
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

// Creates a random LU decomposition so that A = LU is an invertible matrix in
// a form of a "pack" of L and U triangular matrices with a shared diagonal.
//
//	   L          U         LU pack
//	1          1 0 1 0      1 0 1 0
//	0 1          1 1 1  =>  0 1 1 1
//	0 0 1          1 0      0 0 1 0
//	0 0 1 1          1      0 0 1 1
//
// To guarantee the invertibility of A, U and L must be invertible triangular,
// that is the diagonal of each must consist of non-zero values. As GF(2) ring
// (mod 2 integers) only has two values: 0 and 1, "non-zero value" can be read as
// "value of one". That fact allows such a "pack" structure, because the diagonals
// of L and U are identical.
func RandLUPack(n int) Mat {
	data := make([]uint8, n*n)
	rand.Read(data)

	// set the diagonal elements to 1's
	for i := range n {
		data[n*i+i] = 1
	}

	// only leave the least significant bit
	for i := range n * n {
		data[i] &= 1
	}

	return Mat{Dim: n, Data: data}
}

// Transforms the "LU pack" to matrix A, so that A = LU
func FromLUPack(lu *Mat) Mat {
	n := lu.Dim
	mat := IdMat(n)

	for i := range n {
		for j := range n {
			var sum uint8 = 0

			// skip zero factors
			for k := range min(i, j) + 1 {
				ik := lu.Data[n*i+k]
				kj := lu.Data[n*k+j]

				sum = Add(sum, Mul(ik, kj))
			}

			mat.Data[n*i+j] = sum
		}
	}
	return mat
}

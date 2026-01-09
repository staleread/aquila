package main

import "crypto/rand"

// System of Linear Equations (SLE). LU decomposition is used for internal
// representation.
type SLE struct {
	mat *Mat
}

// Creates a random SLE
func RandSLE(n Dim) *SLE {
	mat := Zeros(n)

	buffSize := n * n * ElemLenBytes
	buff := make([]byte, buffSize)
	rand.Read(buff)

	for i := range buffSize {
		elemIdx := i / ElemLenBytes

		// restrict diagonal elements to 1's
		if elemIdx/n == elemIdx%n {
			mat.Data[elemIdx] = ElemMask
			continue
		}

		shift := (i % ElemLenBytes) * 8
		mat.Data[elemIdx] |= Elem(buff[i]) << shift

		if shift == 0 {
			mat.Data[elemIdx] &= ElemMask
		}
	}

	return &SLE{mat}
}

// Solves SLE using LU decomposition.
func (sle *SLE) Solve(x, b Vec) {
	sle.mat.SubForward(x, b)
	sle.mat.SubBackward(x, x)
}

// Returns a matrix of SLE coefficients
func (sle *SLE) Mat() *Mat {
	n := sle.mat.Dim
	mat := Zeros(n)

	for i := range n {
		for j := range n {
			var sum Elem = 0

			// skip zero factors
			for k := range min(i, j) + 1 {
				ik := sle.mat.Data[n*i+k]
				kj := sle.mat.Data[n*k+j]

				sum ^= ik & kj
			}

			mat.Data[n*i+j] = sum
		}
	}
	return mat
}

package main

import "crypto/rand"

// Solves LUx = b equation.
//
// Instead of solving Ax = b, the decomposition step is omitted and a "LU pack"
// is passed resulting in O(n^2) time complexity.
func SolveLu(luPack *BitMat, x, b BitVec) {
	substituteForward(luPack, x, b)
	substituteBackward(luPack, x, x)
}

// Solves Lx = b equation. The diagonal elements of L must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same array slice.
func substituteForward(l *BitMat, x, b BitVec) {
	n := len(b)
	for i := range n {
		var numerator BitPack = b[i]

		for j := range n - 1 {
			numerator ^= l.Data[n*i+j] & b[j]
		}
		x[i] = numerator ^ l.Data[n*i+i]
	}
}

// Solves Ux = b equation. The diagonal elements of U must be non-zero.
// Writes the result to x. It's safe for x and b to point to the same array slice.
func substituteBackward(u *BitMat, x, b BitVec) {
	n := len(b)
	for i := n - 1; i >= 0; i-- {
		var numerator BitPack = b[i]

		for j := i + 1; j < n; j++ {
			numerator ^= u.Data[n*i+j] & b[j]
		}
		x[i] = numerator ^ u.Data[n*i+i]
	}
}

// Creates a random LU pack matrix
func RandLuPack(n Dim) *BitMat {
	const onesPack = (1 << (BytesPerPack * 8)) - 1

	mat := Zeros(n)

	buffSize := n * n * BytesPerPack
	buff := make([]byte, buffSize)
	rand.Read(buff)

	for i := range buffSize {
		elemIdx := i / BytesPerPack

		// restrict diagonal elements to 1's
		if elemIdx/n == elemIdx%n {
			mat.Data[elemIdx] = onesPack
			continue
		}

		shift := (i % BytesPerPack) * 8
		mat.Data[elemIdx] |= BitPack(buff[i]) << shift

		if shift == 0 {
			mat.Data[elemIdx] &= onesPack
		}
	}

	return mat
}

// Undos the LU decomposition
func FromLuPack(luPack *BitMat) *BitMat {
	n := luPack.Dim
	a := Zeros(n)

	for i := range n {
		for j := range n {
			var sum BitPack = 0

			// skip zero factors
			for k := range min(i, j) + 1 {
				ik := luPack.Data[n*i+k]
				kj := luPack.Data[n*k+j]

				sum ^= ik & kj
			}

			a.Data[n*i+j] = sum
		}
	}
	return a
}

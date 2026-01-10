package main

import (
	"fmt"
	"math"
	"strings"
)

type SFoldSLE struct {
	s Size
	*BMat
}

// Creates a random s-fold SLE
func RandSFoldSLE(n, bBytes Size) *SFoldSLE {
	lu := RandBMat(n, bBytes)
	bMask := BatchMask(bBytes)

	// restrict diagonal elements to a batch of 1's
	for i := range n {
		lu.data[n*i+i] = bMask
	}
	return &SFoldSLE{bBytes * 8, lu}
}

// Solves s-fold SLE using LU decomposition.
// Writes the result to x. It's safe for x and b to point to the same vector.
func (sle *SFoldSLE) Solve(x, b *BVec) {
	sle.SubForward(x, b)
	sle.SubBackward(x, x)
}

// Returns a matrix of SLE coefficients
func (sle *SFoldSLE) Coefs() *BMat {
	n := sle.n
	coefs := ZeroBMat(n, sle.bBytes)

	for i := range n {
		for j := range n {
			var sum Batch = 0

			// skip zero factors
			for k := range min(i, j) + 1 {
				sum ^= sle.data[n*i+k] & sle.data[n*k+j]
			}
			coefs.data[n*i+j] = sum
		}
	}
	return coefs
}

func (sle *SFoldSLE) String() string {
	data := sle.Coefs().data
	s := sle.s
	n := sle.n
	sb := strings.Builder{}

	rowPad := Size(math.Log10(float64(s*n)) + 1)
	colPad := Size(math.Log10(float64(n)) + 1)
	varPad := rowPad + colPad + 2 // plus "x" and "," chars

	for i := range s {
		for j := range n {
			row := n*i + j
			sb.WriteString(fmt.Sprintf("y%-*d = ", rowPad, row))

			for k := range n {
				if k > 0 {
					sb.WriteString(" + ")
				}
				x := uint8((data[n*j+k] >> i) & 1)
				var varStr string

				if x != 0 {
					varStr = fmt.Sprintf("x%d,%d", row, k)
				}
				sb.WriteString(fmt.Sprintf("%-*s", varPad, varStr))
			}
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

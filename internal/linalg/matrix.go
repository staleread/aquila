package linalg

import "github.com/staleread/aquila/internal/field"

type SquareMatrix struct {
	n    int
	data Vector
}

func zeroSquareMatrix(n int) *SquareMatrix {
	return &SquareMatrix{
		n:    n,
		data: ZeroVector(n * n),
	}
}

func (self *SquareMatrix) At(i, j int) field.Element {
	return self.data[self.n*i+j]
}

type lowerTriangularMatrix struct {
	*SquareMatrix
}

func randInvertibleLowerTriangularMatrix(n int) *lowerTriangularMatrix {
	sq := zeroSquareMatrix(n)

	dVals := field.RandNonZeroElements(n)
	ndVals := field.RandElements(n * (n + 1) / 2)

	off := 0
	for i := range n {
		off += i

		for j, ndVal := range ndVals[off : off+i] {
			sq.data[n*i+j] = ndVal
		}
		sq.data[n*i+i] = dVals[i]
	}
	return &lowerTriangularMatrix{sq}
}

// Diagonal elements of the matrix must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (self *lowerTriangularMatrix) substituteForward(x, b Vector) {
	for i := range self.n {
		num := b[i]

		for j := range i {
			num = field.Sub(num, field.Mul(self.At(i, j), x[j]))
		}
		x[i] = field.Div(num, self.At(i, i))
	}
}

type upperTriangularMatrix struct {
	*SquareMatrix
}

func randInvertibleUpperTriangularMatrix(n int) *upperTriangularMatrix {
	sq := zeroSquareMatrix(n)

	dVals := field.RandNonZeroElements(n)
	ndVals := field.RandElements(n * (n + 1) / 2)

	off := 0
	for i := range n {
		sq.data[n*i+i] = dVals[i]

		for j, ndVal := range ndVals[off : off+n-(i+1)] {
			sq.data[n*i+i+j+1] = ndVal
		}
		off += n - (i + 1)
	}
	return &upperTriangularMatrix{sq}
}

// Diagonal elements of the matrix of must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (self *upperTriangularMatrix) substituteBackward(x, b Vector) {
	n := self.n

	for i := n - 1; i >= 0; i-- {
		num := b[i]

		for j := i + 1; j < n; j++ {
			num = field.Sub(num, field.Mul(self.At(i, j), x[j]))
		}
		x[i] = field.Div(num, self.At(i, i))
	}
}

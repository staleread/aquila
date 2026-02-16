package la

import f "github.com/staleread/aquila/internal/gf2"

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

func (sq *SquareMatrix) At(i, j int) f.Element {
	return sq.data[sq.n*i+j]
}

// Lower tringular matrix
type ltMatrix struct {
	*SquareMatrix
}

func randInvertibleLtMatrix(n int) *ltMatrix {
	sq := zeroSquareMatrix(n)

	dVals := f.RandNonZeroElements(n)
	ndVals := f.RandElements(n * (n + 1) / 2)

	off := 0
	for i := range n {
		off += i

		for j, ndVal := range ndVals[off : off+i] {
			sq.data[n*i+j] = ndVal
		}
		sq.data[n*i+i] = dVals[i]
	}
	return &ltMatrix{sq}
}

// Performs forward substitution. Diagonal elements of ltMatrix must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (lt *ltMatrix) subForward(x, b Vector) {
	for i := range lt.n {
		num := b[i]

		for j := range i {
			num = f.Sub(num, f.Mul(lt.At(i, j), x[j]))
		}
		x[i] = f.Div(num, lt.At(i, i))
	}
}

// Upper triangular matrix
type utMatrix struct {
	*SquareMatrix
}

func randInvertibleUtMatrix(n int) *utMatrix {
	sq := zeroSquareMatrix(n)

	dVals := f.RandNonZeroElements(n)
	ndVals := f.RandElements(n * (n + 1) / 2)

	off := 0
	for i := range n {
		sq.data[n*i+i] = dVals[i]

		for j, ndVal := range ndVals[off : off+n-(i+1)] {
			sq.data[n*i+i+j+1] = ndVal
		}
		off += n - (i + 1)
	}
	return &utMatrix{sq}
}

// Performs backward substitution. Diagonal elements of utMatrix of must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (ut *utMatrix) subBackward(x, b Vector) {
	n := ut.n

	for i := n - 1; i >= 0; i-- {
		num := b[i]

		for j := i + 1; j < n; j++ {
			num = f.Sub(num, f.Mul(ut.At(i, j), x[j]))
		}
		x[i] = f.Div(num, ut.At(i, i))
	}
}

package la

import f "github.com/staleread/aquila/internal/gf2"

type SqMat struct {
	n    int
	data Vec
}

func zeroSqMat(n int) *SqMat {
	return &SqMat{
		n:    n,
		data: ZeroVec(n * n),
	}
}

func (sq *SqMat) At(i, j int) f.Elt {
	return sq.data[sq.n*i+j]
}

// Lower tringular matrix
type ltMat struct {
	*SqMat
}

func randInvLtMat(n int) *ltMat {
	sq := zeroSqMat(n)

	dVals := f.RandNzEls(n)
	ndVals := f.RandEls(n * (n + 1) / 2)

	off := 0
	for i := range n {
		off += i

		for j, ndVal := range ndVals[off : off+i] {
			sq.data[n*i+j] = ndVal
		}
		sq.data[n*i+i] = dVals[i]
	}
	return &ltMat{sq}
}

// Performs forward substitution. Diagonal elements of ltMat must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (lt *ltMat) subForward(x, b Vec) {
	for i := range lt.n {
		num := b[i]

		for j := range i {
			num = f.Sub(num, f.Mul(lt.At(i, j), x[j]))
		}
		x[i] = f.Div(num, lt.At(i, i))
	}
}

// Upper triangular matrix
type utMat struct {
	*SqMat
}

func randInvUtMat(n int) *utMat {
	sq := zeroSqMat(n)

	dVals := f.RandNzEls(n)
	ndVals := f.RandEls(n * (n + 1) / 2)

	off := 0
	for i := range n {
		sq.data[n*i+i] = dVals[i]

		for j, ndVal := range ndVals[off : off+n-(i+1)] {
			sq.data[n*i+i+j+1] = ndVal
		}
		off += n - (i + 1)
	}
	return &utMat{sq}
}

// Performs backward substitution. Diagonal elements of utMat of must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (ut *utMat) subBackward(x, b Vec) {
	n := ut.n

	for i := n - 1; i >= 0; i-- {
		num := b[i]

		for j := i + 1; j < n; j++ {
			num = f.Sub(num, f.Mul(ut.At(i, j), x[j]))
		}
		x[i] = f.Div(num, ut.At(i, i))
	}
}

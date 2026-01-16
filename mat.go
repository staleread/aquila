package main

// Square matrix
type sqMat struct {
	n    int
	data vec
}

func zeroSqMat(n int) *sqMat {
	return &sqMat{
		n:    n,
		data: zeros(n * n),
	}
}

func (m *sqMat) at(i, j int) elem {
	return m.data[m.n*i+j]
}

// Lower tringular matrix
type ltMat struct {
	*sqMat
}

func randInvLtMat(n int) *ltMat {
	sq := zeroSqMat(n)

	for i := range n {
		for j := range i {
			sq.data[n*i+j] = randElem()
		}
		sq.data[n*i+i] = randNzElem()
	}
	return &ltMat{sq}
}

// Performs forward substitution. Diagonal elements of ltMat must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (lt *ltMat) subForward(x, b vec) {
	for i := range lt.n {
		num := b[i]

		for j := range i {
			num = sub(num, mul(lt.at(i, j), x[j]))
		}
		x[i] = div(num, lt.at(i, i))
	}
}

// Upper triangular matrix
type utMat struct {
	*sqMat
}

func randInvUtMat(n int) *utMat {
	sq := zeroSqMat(n)

	for i := range n {
		sq.data[n*i+i] = randNzElem()
		for j := i + 1; j < n; j++ {
			sq.data[n*i+j] = randElem()
		}
	}
	return &utMat{sq}
}

// Performs backward substitution. Diagonal elements of utMat of must be non-zero.
// The result is written to x. It's safe for x and b to point to the same vector.
func (ut *utMat) subBackward(x, b vec) {
	n := ut.n

	for i := n - 1; i >= 0; i-- {
		num := b[i]

		for j := i + 1; j < n; j++ {
			num = sub(num, mul(ut.at(i, j), x[j]))
		}
		x[i] = div(num, ut.at(i, i))
	}
}

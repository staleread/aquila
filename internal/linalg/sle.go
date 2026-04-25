package linalg

import "github.com/staleread/aquila/internal/field"

type SLE struct {
	lt *lowerTriangularMatrix
	ut *upperTriangularMatrix
}

func RandSLE(n int) *SLE {
	return &SLE{
		lt: randInvertibleLowerTriangularMatrix(n),
		ut: randInvertibleUpperTriangularMatrix(n),
	}
}

func (self *SLE) Solve(dst, src Vector) {
	self.lt.substituteForward(dst, src)
	self.ut.substituteBackward(dst, dst)
}

func (self *SLE) Eval(dst, src Vector) {
	n := self.lt.n
	tmp := ZeroVector(n)

	// U * src = tmp
	for i := range n {
		var sum field.Element = 0

		for j := i; j < n; j++ {
			sum = field.Add(sum, field.Mul(self.ut.At(i, j), src[j]))
		}
		tmp[i] = sum
	}

	// L * tmp = dst
	for i := range n {
		var sum field.Element = 0

		for j := range i + 1 {
			sum = field.Add(sum, field.Mul(self.lt.At(i, j), tmp[j]))
		}
		dst[i] = sum
	}
}

func (self *SLE) Coefs() *SquareMatrix {
	n := self.lt.n
	coefs := zeroSquareMatrix(n)

	for i := range n {
		for j := range n {
			var sum field.Element = 0

			for k := range min(i, j) + 1 { // skip zero factors
				sum = field.Add(sum, field.Mul(self.lt.At(i, k), self.ut.At(k, j)))
			}
			coefs.data[n*i+j] = sum
		}
	}
	return coefs
}

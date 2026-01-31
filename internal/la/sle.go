package la

import f "github.com/staleread/aquila/internal/gf2"

type SLE struct {
	lt *ltMat
	ut *utMat
}

func RandSLE(n int) *SLE {
	return &SLE{
		lt: randInvLtMat(n),
		ut: randInvUtMat(n),
	}
}

func (s *SLE) Solve(dst, src Vec) {
	s.lt.subForward(dst, src)
	s.ut.subBackward(dst, dst)
}

func (s *SLE) Eval(dst, src Vec) {
	n := s.lt.n
	tmp := ZeroVec(n)

	// U * src = tmp
	for i := range n {
		var sum f.Elt = 0

		for j := i; j < n; j++ {
			sum = f.Add(sum, f.Mul(s.ut.At(i, j), src[j]))
		}
		tmp[i] = sum
	}

	// L * tmp = dst
	for i := range n {
		var sum f.Elt = 0

		for j := range i + 1 {
			sum = f.Add(sum, f.Mul(s.lt.At(i, j), tmp[j]))
		}
		dst[i] = sum
	}
}

func (s *SLE) Coefs() *SqMat {
	n := s.lt.n
	coefs := zeroSqMat(n)

	for i := range n {
		for j := range n {
			var sum f.Elt = 0

			for k := range min(i, j) + 1 { // skip zero factors
				sum = f.Add(sum, f.Mul(s.lt.At(i, k), s.ut.At(k, j)))
			}
			coefs.data[n*i+j] = sum
		}
	}
	return coefs
}

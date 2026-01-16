package main

type sle struct {
	lt *ltMat
	ut *utMat
}

func randSle(n int) *sle {
	return &sle{
		lt: randInvLtMat(n),
		ut: randInvUtMat(n),
	}
}

// Solves SLE.
// Writes the result to x. It's safe for x and b to point to the same vector.
func (s *sle) solve(x, b vec) {
	s.lt.subForward(x, b)
	s.ut.subBackward(x, x)
}

// Evaluates SLE.
func (s *sle) eval(x, b vec) {
	n := s.lt.n
	tmp := zeros(n)

	// U * x = tmp
	for i := range n {
		var sum elem = 0

		for j := i; j < n; j++ {
			sum = add(sum, mul(s.ut.at(i, j), x[j]))
		}
		tmp[i] = sum
	}

	// L * tmp = b
	for i := range n {
		var sum elem = 0

		for j := range i + 1 {
			sum = add(sum, mul(s.lt.at(i, j), tmp[j]))
		}
		b[i] = sum
	}
}

// Returns a matrix of SLE coefficients.
func (s *sle) coefs() *sqMat {
	n := s.lt.n
	coefs := zeroSqMat(n)

	for i := range n {
		for j := range n {
			var sum elem = 0

			for k := range min(i, j) + 1 { // skip zero factors
				sum = add(sum, mul(s.lt.at(i, k), s.ut.at(k, j)))
			}
			coefs.data[n*i+j] = sum
		}
	}
	return coefs
}

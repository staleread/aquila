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
func (sle *sle) solve(x, b vec) {
	sle.lt.subForward(x, b)
	sle.ut.subBackward(x, x)
}

// Evaluates SLE.
func (sle *sle) eval(x, b vec) {
	n := sle.lt.n
	tmp := zeros(n)

	// U * x = tmp
	for i := range n {
		var sum elem = 0

		for j := i; j < n; j++ {
			sum = add(sum, mul(sle.ut.at(i, j), x[j]))
		}
		tmp[i] = sum
	}

	// L * tmp = b
	for i := range n {
		var sum elem = 0

		for j := range i + 1 {
			sum = add(sum, mul(sle.lt.at(i, j), tmp[j]))
		}
		b[i] = sum
	}
}

// Returns a matrix of SLE coefficients.
func (sle *sle) coefs() *sqMat {
	n := sle.lt.n
	coefs := zeroSqMat(n)

	for i := range n {
		for j := range n {
			var sum elem = 0

			for k := range min(i, j) + 1 { // skip zero factors
				sum = add(sum, mul(sle.lt.at(i, k), sle.ut.at(k, j)))
			}
			coefs.data[n*i+j] = sum
		}
	}
	return coefs
}

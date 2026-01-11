package main

type SLE struct {
	*Mat
}

// Creates a random SLE
func RandSLE(n Size) *SLE {
	lu := RandMat(n)

	// restrict diagonal elements to a batch of 1's
	for i := range n {
		lu.data[n*i+i] = 1
	}
	return &SLE{lu}
}

// Solves SLE using LU decomposition.
// Writes the result to x. It's safe for x and b to point to the same vector.
func (sle *SLE) Solve(x, b Vec) {
	sle.SubForward(x, b)
	sle.SubBackward(x, x)
}

// Returns a matrix of SLE coefficients
func (sle *SLE) Coefs() *Mat {
	n := sle.n
	coefs := ZeroMat(n)

	for i := range n {
		for j := range n {
			var sum Elem = 0

			for k := range min(i, j) + 1 { // skip zero factors
				sum ^= sle.data[n*i+k] & sle.data[n*k+j]
			}
			coefs.data[n*i+j] = sum
		}
	}
	return coefs
}

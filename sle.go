package main

// TODO add docs
type SLE struct {
	lu *Mat
}

// Creates a random SLE
func RandSLE(n Dim) *SLE {
	lu := RandMat(n)

	// restrict diagonal elements to 1's
	for i := range n {
		lu.Data[n*i+i] = ElemMask
	}
	return &SLE{lu}
}

// Solves SLE using LU decomposition.
func (sle *SLE) Solve(x, b Vec) {
	sle.lu.SubForward(x, b)
	sle.lu.SubBackward(x, x)
}

// Returns a matrix of SLE coefficients
func (sle *SLE) Coefs() *Mat {
	n := sle.lu.Dim
	coefs := ZeroMat(n)

	for i := range n {
		for j := range n {
			var sum Elem = 0

			// skip zero factors
			for k := range min(i, j) + 1 {
				ik := sle.lu.Data[n*i+k]
				kj := sle.lu.Data[n*k+j]

				sum ^= ik & kj
			}

			coefs.Data[n*i+j] = sum
		}
	}
	return coefs
}

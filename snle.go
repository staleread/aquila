package main

type SNLE struct {
	data []Poly
}

func EmptySNLE(n Size) *SNLE {
	return &SNLE{data: make([]Poly, n)}
}

func RandSNLE(n, deg Size, maxVarIdx VarIdx) *SNLE {
	snle := EmptySNLE(n)

	for i := range n {
		snle.data[i] = RandPoly(deg, maxVarIdx)
	}
	return snle
}

func (snle *SNLE) Eval(x, b Vec) {
	for i, p := range snle.data {
		var sum Elem = 0

		for _, m := range p {
			var factor Elem = 1

			for _, idx := range m {
				factor &= x[idx]
			}
			sum ^= factor
		}
		// Use GF(2) addition(=subraction) instead of overriding b in case b is
		// a non-zero vector.
		b[i] ^= sum
	}
}

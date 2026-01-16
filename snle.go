package main

type snle struct {
	data []poly
}

func emptySnle(n int) *snle {
	return &snle{make([]poly, n)}
}

func randSnle(n, deg int, maxIdx idx) *snle {
	snle := emptySnle(n)

	for i := range n {
		snle.data[i] = randPoly(deg, maxIdx)
	}
	return snle
}

func (snle *snle) eval(x, b vec) {
	for i, p := range snle.data {
		var sum elem = 0

		for _, m := range p {
			var prod elem = 1

			for _, idx := range m {
				prod = mul(prod, x[idx])
			}
			sum = add(sum, prod)
		}
		b[i] = sum
	}
}

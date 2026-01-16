package main

// Multilinear Polynomial
type poly []monom

func randPoly(deg int, maxIdx idx) poly {
	p := make(poly, deg)

	for i := range deg {
		p[i] = randMonom(deg-i, maxIdx)
	}
	return p
}

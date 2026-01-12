package main

import "math/rand"

type Monom []VarIdx
type Poly []Monom

func RandPoly(deg Size, maxIdx VarIdx) Poly {
	ids := make([]VarIdx, maxIdx)

	for i := range maxIdx {
		ids[i] = i
	}
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	p := make(Poly, deg)
	var added Size = 0

	for i := range deg {
		m := make(Monom, deg-i)

		for j := range deg - i {
			m[j] = ids[(added+j)%maxIdx]
		}
		p[i] = m
		added += deg - i
	}
	return p
}

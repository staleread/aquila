package main

import "crypto/rand"

type Vec []Elem

func Zeros(n Size) Vec {
	return Vec(make([]Elem, n))
}

func Rands(n Size) Vec {
	vec := Zeros(n)
	rand.Read(vec)

	for i := range n {
		vec[i] &= 1
	}
	return vec
}

func (vec Vec) Permute(perm []VarIdx) {
	for i, j := range perm {
		vec[i], vec[j] = vec[j], vec[i]
	}
}

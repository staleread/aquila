package invertible

import "github.com/staleread/aquila/internal/field"

type permutation []int

func randPermutation(n int) permutation {
	p := make(permutation, n-1)
	rands := field.RandSubscripts(n - 1)

	// Fisher–Yates shuffle
	for i := range n - 1 {
		p[i] = int(rands[i])%(n-i) + i
	}
	return p
}

func (p permutation) permute(v []field.Element) {
	if len(v) != len(p)+1 {
		panic("Array size doesn't match the permutation size")
	}

	for i, j := range p {
		v[i], v[j] = v[j], v[i]
	}
}

func (p permutation) permuteBack(v []field.Element) {
	if len(v) != len(p)+1 {
		panic("Array size doesn't match the permutation size")
	}

	for i := len(p) - 1; i >= 0; i-- {
		j := p[i]
		v[i], v[j] = v[j], v[i]
	}
}

func (p permutation) ids() []field.Subscript {
	n := len(p) + 1
	ids := make([]field.Subscript, n)

	for i := range field.Subscript(n) {
		ids[i] = i
	}

	for i, j := range p {
		ids[i], ids[j] = ids[j], ids[i]
	}
	return ids
}

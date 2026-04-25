package invertible

import (
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/poly"
)

type permutation []int

func randPermutation(n int) permutation {
	p := make(permutation, n-1)
	rands := poly.RandSubscripts(n - 1)

	// Fisher–Yates shuffle
	for i := range n - 1 {
		p[i] = int(rands[i])%(n-i) + i
	}
	return p
}

func (self permutation) permute(v []field.Element) {
	if len(v) != len(self)+1 {
		panic("Array size doesn't match the permutation size")
	}

	for i, j := range self {
		v[i], v[j] = v[j], v[i]
	}
}

func (self permutation) permuteBack(v []field.Element) {
	if len(v) != len(self)+1 {
		panic("Array size doesn't match the permutation size")
	}

	for i := len(self) - 1; i >= 0; i-- {
		j := self[i]
		v[i], v[j] = v[j], v[i]
	}
}

func (self permutation) ids() []poly.Subscript {
	n := len(self) + 1
	ids := make([]poly.Subscript, n)

	for i := range poly.Subscript(n) {
		ids[i] = i
	}

	for i, j := range self {
		ids[i], ids[j] = ids[j], ids[i]
	}
	return ids
}

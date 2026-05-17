package general

import (
	"github.com/staleread/aquila/internal/automata"
)

// Adds two pre-sorted polinomials using two-pointer merge
// Appends the resulting monomials to the 'dst' buffer.
func AddPolynomials(dst, p1, p2 []automata.Word) []automata.Word {
	var ptr1, ptr2 int

	for ptr1 < len(p1) && ptr2 < len(p2) {
		monomA := NewMonomial(p1[ptr1 : ptr1+automata.BlockWords])
		monomB := NewMonomial(p2[ptr2 : ptr2+automata.BlockWords])

		cmp := CompareMonomials(monomA, monomB)

		switch {
		case cmp > 0:
			dst = append(dst, p1[ptr1:ptr1+automata.BlockWords]...)
			ptr1 += automata.BlockWords

		case cmp < 0:
			dst = append(dst, p2[ptr2:ptr2+automata.BlockWords]...)
			ptr2 += automata.BlockWords

		default:
			// Identical terms cancel each other (a ^ a = 0), so skip
			ptr1 += automata.BlockWords
			ptr2 += automata.BlockWords
		}
	}

	if ptr1 < len(p1) {
		dst = append(dst, p1[ptr1:]...)
	}

	if ptr2 < len(p2) {
		dst = append(dst, p2[ptr2:]...)
	}
	return dst
}

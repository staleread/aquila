package math

import (
	"slices"
)

// Multiplies two polynomials and overwrites the 'dst' buffer with the resulting
// polynomial. The resulting monomials are sorted and deduplicated.
func MultiplyPolynomials(dst, p1, p2 []Monomial) []Monomial {
	if len(p1) == 0 || len(p2) == 0 {
		return dst[:0]
	}

	total := len(p1) * len(p2)

	if cap(dst) < total {
		dst = make([]Monomial, total)
	} else {
		dst = dst[:total]
	}

	for i, m1 := range p1 {
		for j, m2 := range p2 {
			dst[i*len(p2)+j] = m1.Mul(m2)
		}
	}

	slices.SortFunc(dst, func(a, b Monomial) int {
		return CompareMonomials(b, a)
	})

	writeIdx := 0
	for readIdx := 0; readIdx < len(dst); {
		curr := dst[readIdx]
		count := 1
		readIdx++

		for readIdx < len(dst) && dst[readIdx] == curr {
			count++
			readIdx++
		}

		if count%2 != 0 {
			dst[writeIdx] = curr
			writeIdx++
		}
	}

	return dst[:writeIdx]
}

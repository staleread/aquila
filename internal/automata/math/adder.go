package math

// Adds two pre-sorted polinomials using two-pointer merge
// Appends the resulting monomials to the 'dst' buffer.
func AddPolynomials(dst, p1, p2 []Monomial) []Monomial {
	var ptr1, ptr2 int

	for ptr1 < len(p1) && ptr2 < len(p2) {
		m1 := p1[ptr1]
		m2 := p2[ptr2]

		cmp := CompareMonomials(m1, m2)

		switch {
		case cmp > 0:
			dst = append(dst, m1)
			ptr1++

		case cmp < 0:
			dst = append(dst, m2)
			ptr2++

		default:
			// Identical terms cancel each other (a ^ a = 0), so skip
			ptr1++
			ptr2++
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

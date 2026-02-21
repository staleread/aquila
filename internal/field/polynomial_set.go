package field

import (
	"fmt"
	"strings"
)

type PolynomialSet []Polynomial

func RandPolynomialSet(n, degree int, maxSub Subscript) PolynomialSet {
	noise := make(PolynomialSet, n)

	for i := range n {
		noise[i] = randPolynomial(degree, maxSub)
	}
	return noise
}

func (set PolynomialSet) Eval(dst, src []Element) {
	for i := range len(dst) {
		dst[i] = set[i].Eval(src)
	}
}

func (set PolynomialSet) Compose(other PolynomialSet) {
	// type cacheEntry struct {
	// 	key   Monomial
	// 	value Polynomial
	// }
	// prodCache := make(map[uint32][]cacheEntry)

	// lookupCache := func(m Monomial) (Polynomial, bool) {
	// 	for _, entry := range prodCache[m.hash] {
	// 		if m.Equals(entry.key) {
	// 			return entry.value, true
	// 		}
	// 	}
	// 	return nil, false
	// }

	// storeCache := func(m Monomial, p Polynomial) {
	// 	h := m.hash
	// 	prodCache[h] = append(prodCache[h], cacheEntry{m, p})
	// }

	for j, p := range set {
		sum := ZeroPolynomial

		for m := range p.Monomials() {
			// if prod, ok := lookupCache(m); ok {
			// 	sum = sum.Add(prod)
			// 	continue
			// }

			prod := OnePolynomial

			for s := range m.Syms() {
				prod = prod.Mul(other[s])
			}

			sum = sum.Add(prod)
			// storeCache(m, prod)
		}
		set[j] = sum
	}
}

func (set PolynomialSet) String() string {
	sb := strings.Builder{}

	for i, p := range set {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		firstMonomial := true
		for m := range p.Monomials() {
			if !firstMonomial {
				sb.WriteString(" + ")
			}
			firstMonomial = false

			firstSubscript := true
			for s := range m.Syms() {
				if !firstSubscript {
					sb.WriteRune('*')
				}

				firstSubscript = false
				fmt.Fprintf(&sb, "x%d", s+1)
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

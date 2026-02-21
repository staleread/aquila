package field

import (
	"fmt"
	"strings"
)

type PolynomialSet []PolynomialPtr

func RandPolynomialSet(n, degree int, maxSub Subscript) PolynomialSet {
	noise := make(PolynomialSet, n)

	for i := range n {
		noise[i] = randPolynomial(degree, maxSub)
	}
	return noise
}

func (set PolynomialSet) Eval(dst, src []Element) {
	for i := range len(dst) {
		dst[i] = evalPolynomial(set[i], src)
	}
}

func (set PolynomialSet) Compose(other PolynomialSet) {
	prodCache := make(map[MonomialPtr]PolynomialPtr)

	for j, p := range set {
		sum := PolynomialPtr(0)

		for m := range MonomialsOf(p) {
			if prod, ok := prodCache[m]; ok {
				sum = AddPolynomials(sum, prod)
				continue
			}
			prod := PolynomialPtr(0)

			for s := range SubscriptsOf(m) {
				nextPoly := other[s]

				if prod == 0 {
					prod = nextPoly
					continue
				}
				prod = MulPolynomials(prod, nextPoly)
			}
			sum = AddPolynomials(sum, prod)
			prodCache[m] = prod
		}
		set[j] = sum
	}
}

func (set PolynomialSet) String() string {
	sb := strings.Builder{}

	for i, p := range set {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		firstMonomial := true
		for m := range MonomialsOf(p) {
			if !firstMonomial {
				sb.WriteString(" + ")
			}
			firstMonomial = false

			firstSubscript := true
			for s := range SubscriptsOf(m) {
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

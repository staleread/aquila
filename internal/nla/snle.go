package nla

import (
	"fmt"
	f "github.com/staleread/aquila/internal/gf2"
	"strings"
)

type SNLE []f.Polynomial

func ZeroSNLE(n int) SNLE {
	se := SNLE(make([]f.Polynomial, n))

	for i := range n {
		se[i] = f.ZeroPolynomial
	}
	return se
}

func NewSNLE(arr []f.Polynomial) SNLE {
	return SNLE(arr)
}

func RandSNLE(n, deg int, maxSubscript f.Subscript) SNLE {
	se := SNLE(make([]f.Polynomial, n))

	for i := range n {
		se[i] = f.RandPolynomial(deg, maxSubscript)
	}
	return se
}

func (a SNLE) Compose(b SNLE) SNLE {
	newSe := SNLE(make([]f.Polynomial, len(a)))

	for i, p := range a {
		sum := f.ZeroPolynomial

		for m := range p.Monomials() {
			prod := f.OnePolynomial

			for s := range m.Subscripts() {
				prod = prod.Mul(b[s])
			}
			sum = sum.Add(prod)
		}
		newSe[i] = sum
	}
	return newSe
}

func (se SNLE) Eval(dst, src []f.Element) {
	for i := range len(dst) {
		dst[i] = se[i].Eval(src)
	}
}

func (se SNLE) Polies() []f.Polynomial {
	return []f.Polynomial(se)
}

func (se SNLE) String() string {
	sb := strings.Builder{}

	for i, p := range se {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		firstMonomial := true
		for m := range p.Monomials() {
			if !firstMonomial {
				sb.WriteString(" + ")
			}
			firstMonomial = false

			firstSubscript := true
			for s := range m.Subscripts() {
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

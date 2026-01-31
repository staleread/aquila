package nla

import (
	"fmt"
	f "github.com/staleread/aquila/internal/gf2"
	"strings"
)

type SNLE []f.Poly

func ZeroSNLE(n int) SNLE {
	se := SNLE(make([]f.Poly, n))

	for i := range n {
		se[i] = f.ZeroPoly
	}
	return se
}

func NewSNLE(arr []f.Poly) SNLE {
	return SNLE(arr)
}

func RandSNLE(n, deg int, maxSym f.Sym) SNLE {
	se := SNLE(make([]f.Poly, n))

	for i := range n {
		se[i] = f.RandPoly(deg, maxSym)
	}
	return se
}

func (a SNLE) Compose(b SNLE) SNLE {
	newSe := SNLE(make([]f.Poly, len(a)))

	for i, p := range a {
		sum := f.ZeroPoly

		for m := range p.Monoms() {
			prod := f.OnePoly

			for s := range m.Syms() {
				prod = prod.Mul(b[s])
			}
			sum = sum.Add(prod)
		}
		newSe[i] = sum
	}
	return newSe
}

func (se SNLE) Eval(dst, src []f.Elt) {
	for i := range len(dst) {
		dst[i] = se[i].Eval(src)
	}
}

func (se SNLE) Polies() []f.Poly {
	return []f.Poly(se)
}

func (se SNLE) String() string {
	sb := strings.Builder{}

	for i, p := range se {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		firstMonom := true
		for m := range p.Monoms() {
			if !firstMonom {
				sb.WriteString(" + ")
			}
			firstMonom = false

			firstSym := true
			for s := range m.Syms() {
				if !firstSym {
					sb.WriteRune('*')
				}

				firstSym = false
				fmt.Fprintf(&sb, "x%d", s+1)
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

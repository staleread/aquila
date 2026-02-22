package field

import (
	"iter"
	"maps"
)

type Polynomial map[*Monomial]struct{}

// func NewPolynomial(monoms []Monomial) Polynomial {
// 	p := make(Polynomial, len(monoms))

// 	for _, m := range monoms {
// 		p.ToggleMonomial(&m)
// 	}
// 	return p
// }

func (p Polynomial) Monomials() iter.Seq[*Monomial] {
	return maps.Keys(p)
}

func (p Polynomial) Eval(x []Element) Element {
	var sum Element = 0

	for m := range p.Monomials() {
		sum = Add(sum, m.Eval(x))
	}
	return sum
}

func (a Polynomial) AddTo(b Polynomial) {
	for m := range b.Monomials() {
		a.ToggleMonomial(m)
	}
}

func (poly Polynomial) ToggleMonomial(mPtr *Monomial) {
	if _, ok := poly[mPtr]; ok {
		delete(poly, mPtr)
	} else {
		poly[mPtr] = struct{}{}
	}
}

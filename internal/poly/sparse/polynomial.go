package sparse

import (
	"iter"
	"maps"
)

type Polynomial map[*Monomial]struct{}

func (poly Polynomial) Monomials() iter.Seq[*Monomial] {
	return maps.Keys(poly)
}

func (a Polynomial) AddTo(b Polynomial) {
	for m := range b.Monomials() {
		a.ToggleMonomial(m)
	}
}

func (poly Polynomial) ToggleMonomial(monom *Monomial) {
	if _, ok := poly[monom]; ok {
		delete(poly, monom)
	} else {
		poly[monom] = struct{}{}
	}
}

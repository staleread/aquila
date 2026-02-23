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

func (poly Polynomial) ToggleMonomial(mPtr *Monomial) {
	if _, ok := poly[mPtr]; ok {
		delete(poly, mPtr)
	} else {
		poly[mPtr] = struct{}{}
	}
}

package canonical

import "github.com/staleread/aquila/internal/field"

type ComputationArena struct {
	*MonomialArena
}

func NewComputationArena() *ComputationArena {
	return &ComputationArena{newMonomialArena()}
}

func (arena *ComputationArena) MulPolynomialBy(a, b field.Polynomial) {
	aMonoms := make([]*field.Monomial, 0, len(a))
	for ma := range a.Monomials() {
		aMonoms = append(aMonoms, ma)
	}

	for _, ma := range aMonoms {
		a.ToggleMonomial(ma)

		for mb := range b.Monomials() {
			mPtr := arena.GetOrInsertMonomial(ma.Mul(*mb))
			a.ToggleMonomial(mPtr)
		}
	}
}

package canonical

import "github.com/staleread/aquila/internal/field"

type AllocationArena struct {
	*MonomialArena
}

func NewAllocationArena() *AllocationArena {
	return &AllocationArena{newMonomialArena()}
}

func (arena *AllocationArena) RandPolyset(n, degree int, maxSub field.Subscript) field.Polyset {
	set := make(field.Polyset, n)

	for i := range n {
		poly := make(field.Polynomial, degree)

		for i := range degree {
			m := field.RandMonomial(degree-i, maxSub)
			mPtr := arena.GetOrInsertMonomial(m)
			poly[mPtr] = struct{}{}
		}
		set[i] = poly
	}
	return set
}

func (arena *AllocationArena) CreateMonomial(subs ...field.Subscript) *field.Monomial {
	m := field.NewMonomial(subs...)
	return arena.GetOrInsertMonomial(m)
}

package sparse

import "github.com/staleread/aquila/internal/poly"

type PolysetBuilder struct {
	polynoms  []Polynomial
	monomPool *monomialInternPool
}

func NewPolysetBuilder(size int) *PolysetBuilder {
	polynoms := make([]Polynomial, size)

	for i := range size {
		polynoms[i] = Polynomial{}
	}

	monomPool := newMonomialInternPool()

	return &PolysetBuilder{polynoms, monomPool}
}

func (builder *PolysetBuilder) AddMonomOf(polynomIdx int, subs ...poly.Subscript) {
	monom := builder.monomPool.createMonomial(subs...)

	builder.polynoms[polynomIdx].ToggleMonomial(monom)
}

func (builder *PolysetBuilder) Build() *Polyset {
	return &Polyset{
		polynoms:  builder.polynoms,
		monomPool: builder.monomPool,
	}
}

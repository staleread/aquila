package mlise

import (
	f "github.com/staleread/aquila/internal/gf2"
	"github.com/staleread/aquila/internal/la"
	"github.com/staleread/aquila/internal/nla"
)

// Multifold linear invertible system of equations
type MLISE struct {
	size  int
	permutation  permutation
	folds []fold
}

type fold struct {
	lin   *la.SLE
	noise nla.SNLE
}

func RandMLISE(size, folds, deg int) *MLISE {
	permutation := randPermutation(size)
	n := size / folds
	sFolds := make([]fold, folds)

	sFolds[0] = fold{
		lin:   la.RandSLE(n),
		noise: nla.ZeroSNLE(n),
	}

	for i := 1; i < folds; i++ {
		maxSubscript := f.Subscript(n * i)

		sFolds[i] = fold{
			lin:   la.RandSLE(n),
			noise: nla.RandSNLE(n, deg, maxSubscript),
		}
	}
	return &MLISE{size, permutation, sFolds}
}

func (ms *MLISE) Eval(dst, src la.Vector) {
	n := ms.size / len(ms.folds)

	ms.permutation.permute(src)

	for i, fl := range ms.folds {
		xCurr := src[n*i : n*i+n]
		bCurr := dst[n*i : n*i+n]

		fl.lin.Eval(bCurr, xCurr)

		noise := la.ZeroVector(n)
		xPrev := src[:n*i]

		fl.noise.Eval(noise, xPrev)
		bCurr.Add(noise)
	}
}

func (ms *MLISE) Solve(dst, src la.Vector) {
	n := ms.size / len(ms.folds)

	for i, fl := range ms.folds {
		noise := la.ZeroVector(n)
		xPrev := dst[:n*i]
		bCurr := src[n*i : n*i+n]

		fl.noise.Eval(noise, xPrev)
		bCurr.Sub(noise)

		xCurr := dst[n*i : n*i+n]

		fl.lin.Solve(xCurr, bCurr)
	}
	ms.permutation.permuteBack(dst)
}

func (ms *MLISE) ToSNLE() nla.SNLE {
	n := ms.size / len(ms.folds)
	ids := ms.permutation.ids()
	polies := make([]f.Polynomial, 0, ms.size)

	for i, fl := range ms.folds {
		lin := fl.lin.Coefs()
		noise := fl.noise

		for j, p := range noise.Polies() {
			monomials := make([]f.Monomial, 0, n)

			// Non-linear part
			for m := range p.Monomials() {
				subscripts := make([]f.Subscript, 0, m.Size())

				for s := range m.Subscripts() {
					subscripts = append(subscripts, ids[s])
				}
				monomials = append(monomials, f.NewMonomial(subscripts...))
			}

			// Linear part
			for k := range n {
				val := lin.At(j, k)

				if val == 0 {
					continue
				}
				s := ids[n*i+k]
				monomials = append(monomials, f.NewMonomial(s))
			}
			polies = append(polies, f.NewPolynomial(monomials))
		}
	}
	return nla.NewSNLE(polies)
}

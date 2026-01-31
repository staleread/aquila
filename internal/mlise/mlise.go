package mlise

import (
	f "github.com/staleread/aquila/internal/gf2"
	"github.com/staleread/aquila/internal/la"
	"github.com/staleread/aquila/internal/nla"
)

// Multifold linear invertible system of equations
type MLISE struct {
	size  int
	perm  perm
	folds []fold
}

type fold struct {
	lin   *la.SLE
	noise nla.SNLE
}

func RandMLISE(size, folds, deg int) *MLISE {
	perm := randPerm(size)
	n := size / folds
	sFolds := make([]fold, folds)

	sFolds[0] = fold{
		lin:   la.RandSLE(n),
		noise: nla.ZeroSNLE(n),
	}

	for i := 1; i < folds; i++ {
		maxSym := f.Sym(n * i)

		sFolds[i] = fold{
			lin:   la.RandSLE(n),
			noise: nla.RandSNLE(n, deg, maxSym),
		}
	}
	return &MLISE{size, perm, sFolds}
}

func (ms *MLISE) Eval(dst, src la.Vec) {
	n := ms.size / len(ms.folds)

	ms.perm.permute(src)

	for i, fl := range ms.folds {
		xCurr := src[n*i : n*i+n]
		bCurr := dst[n*i : n*i+n]

		fl.lin.Eval(bCurr, xCurr)

		noise := la.ZeroVec(n)
		xPrev := src[:n*i]

		fl.noise.Eval(noise, xPrev)
		bCurr.Add(noise)
	}
}

func (ms *MLISE) Solve(dst, src la.Vec) {
	n := ms.size / len(ms.folds)

	for i, fl := range ms.folds {
		noise := la.ZeroVec(n)
		xPrev := dst[:n*i]
		bCurr := src[n*i : n*i+n]

		fl.noise.Eval(noise, xPrev)
		bCurr.Sub(noise)

		xCurr := dst[n*i : n*i+n]

		fl.lin.Solve(xCurr, bCurr)
	}
	ms.perm.permuteBack(dst)
}

func (ms *MLISE) ToSNLE() nla.SNLE {
	n := ms.size / len(ms.folds)
	ids := ms.perm.ids()
	polies := make([]f.Poly, 0, ms.size)

	for i, fl := range ms.folds {
		lin := fl.lin.Coefs()
		noise := fl.noise

		for j, p := range noise.Polies() {
			monoms := make([]f.Monom, 0, n)

			// Non-linear part
			for m := range p.Monoms() {
				syms := make([]f.Sym, 0, m.Size())

				for s := range m.Syms() {
					syms = append(syms, ids[s])
				}
				monoms = append(monoms, f.NewMonom(syms...))
			}

			// Linear part
			for k := range n {
				val := lin.At(j, k)

				if val == 0 {
					continue
				}
				s := ids[n*i+k]
				monoms = append(monoms, f.NewMonom(s))
			}
			polies = append(polies, f.NewPoly(monoms))
		}
	}
	return nla.NewSNLE(polies)
}

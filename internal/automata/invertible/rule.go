package invertible

import (
	"github.com/staleread/aquila/internal/canonical"
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/linalg"
)

type fold struct {
	lin   *linalg.SLE
	noise field.Polyset
}

type rule struct {
	size        int
	permutation permutation
	folds       []fold
}

func randRule(size, folds, degree int, arena *canonical.AllocationArena) *rule {
	permutation := randPermutation(size)
	n := size / folds
	sFolds := make([]fold, folds)

	sFolds[0] = fold{
		lin:   linalg.RandSLE(n),
		noise: make(field.Polyset, n),
	}

	for i := 1; i < folds; i++ {
		maxSub := field.Subscript(n * i)

		sFolds[i] = fold{
			lin:   linalg.RandSLE(n),
			noise: arena.RandPolyset(n, degree, maxSub),
		}
	}
	return &rule{size, permutation, sFolds}
}

func (rule *rule) Apply(dst, src linalg.Vector) {
	n := rule.size / len(rule.folds)

	rule.permutation.permute(src)

	for i, fl := range rule.folds {
		xCurr := src[n*i : n*i+n]
		bCurr := dst[n*i : n*i+n]

		fl.lin.Eval(bCurr, xCurr)

		noise := linalg.ZeroVector(n)
		xPrev := src[:n*i]

		fl.noise.Eval(noise, xPrev)
		bCurr.Add(noise)
	}
}

func (rule *rule) ApplyInverse(dst, src linalg.Vector) {
	n := rule.size / len(rule.folds)

	for i, fl := range rule.folds {
		noise := linalg.ZeroVector(n)
		xPrev := dst[:n*i]
		bCurr := src[n*i : n*i+n]

		fl.noise.Eval(noise, xPrev)
		bCurr.Sub(noise)

		xCurr := dst[n*i : n*i+n]

		fl.lin.Solve(xCurr, bCurr)
	}
	rule.permutation.permuteBack(dst)
}

func (rule *rule) toPolyset(arena *canonical.AllocationArena) field.Polyset {
	n := rule.size / len(rule.folds)
	ids := rule.permutation.ids()
	polyset := make(field.Polyset, rule.size)

	for i, fl := range rule.folds {
		lin := fl.lin.Coefs()
		noise := fl.noise

		for j, p := range noise {
			poly := field.Polynomial{}

			// Non-linear part
			for m := range p.Monomials() {
				subs := make([]field.Subscript, 0)

				for s := range m.Subscripts() {
					subs = append(subs, ids[s])
				}
				poly.ToggleMonomial(arena.CreateMonomial(subs...))
			}

			// Linear part
			for k := range n {
				val := lin.At(j, k)

				if val == 0 {
					continue
				}
				s := ids[n*i+k]
				poly.ToggleMonomial(arena.CreateMonomial(s))
			}
			polyset[n*i+j] = poly
		}
	}
	return polyset
}

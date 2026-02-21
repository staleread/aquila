package invertible

import (
	"github.com/staleread/aquila/internal/ca/general"
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/linalg"
)

type fold struct {
	lin   *linalg.SLE
	noise field.PolynomialSet
}

type rule struct {
	size        int
	permutation permutation
	folds       []fold
}

func randRule(size, folds, degree int) *rule {
	permutation := randPermutation(size)
	n := size / folds
	sFolds := make([]fold, folds)

	sFolds[0] = fold{
		lin:   linalg.RandSLE(n),
		noise: make(field.PolynomialSet, n),
	}

	for i := 1; i < folds; i++ {
		maxSub := field.Subscript(n * i)

		sFolds[i] = fold{
			lin:   linalg.RandSLE(n),
			noise: field.RandPolynomialSet(n, degree, maxSub),
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

func (rule *rule) general() general.Rule {
	n := rule.size / len(rule.folds)
	ids := rule.permutation.ids()
	polies := make([]field.Polynomial, 0, rule.size)

	for i, fl := range rule.folds {
		lin := fl.lin.Coefs()
		noise := fl.noise

		for j, p := range noise {
			monoms := make([]field.Monomial, 0, n)

			// Non-linear part
			for m := range p.Monomials() {
				subs := make([]field.Subscript, 0)

				for s := range m.Syms() {
					subs = append(subs, ids[s])
				}
				monoms = append(monoms, field.NewMonomial(subs...))
			}

			// Linear part
			for k := range n {
				val := lin.At(j, k)

				if val == 0 {
					continue
				}
				s := ids[n*i+k]
				monoms = append(monoms, field.NewMonomial(s))
			}
			polies = append(polies, field.NewPolynomial(monoms))
		}
	}
	return polies
}

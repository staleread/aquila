package invertible

import (
	"github.com/staleread/aquila/internal/linalg"
	"github.com/staleread/aquila/internal/poly"
	"github.com/staleread/aquila/internal/poly/compiled"
	"github.com/staleread/aquila/internal/poly/sparse"
)

type fold struct {
	lin   *linalg.SLE
	noise *compiled.Polyset
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
		noise: compiled.EmptyPolyset(),
	}

	for i := 1; i < folds; i++ {
		maxSub := poly.Subscript(n * i)

		sFolds[i] = fold{
			lin:   linalg.RandSLE(n),
			noise: compiled.RandPolyset(n, degree, maxSub),
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

func (rule *rule) toSparsePolyset() *sparse.Polyset {
	n := rule.size / len(rule.folds)
	builder := sparse.NewPolysetBuilder(rule.size)

	ids := rule.permutation.ids()

	for i, fl := range rule.folds {
		lin := fl.lin.Coefs()
		noise := fl.noise

		for j := range n {
			polynomIdx := n*i + j

			// Linear part
			for k := range n {
				val := lin.At(j, k)

				if val == 0 {
					continue
				}
				sub := ids[n*i+k]
				builder.AddMonomOf(polynomIdx, sub)
			}

			// Non-linear part
			if noise.PolyCount > 0 {
				pStart := noise.PolyOffsets[j]
				pEnd := noise.PolyOffsets[j+1]

				for k := pStart; k < pEnd; k++ {
					mStart := noise.MonomOffsets[k]
					mEnd := noise.MonomOffsets[k+1]

					subs := make([]poly.Subscript, mEnd-mStart)
					for l, s := range noise.Subscripts[mStart:mEnd] {
						subs[l] = ids[s]
					}
					builder.AddMonomOf(polynomIdx, subs...)
				}
			}
		}
	}
	return builder.Build()
}

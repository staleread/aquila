package ca

import (
	f "github.com/staleread/aquila/internal/gf2"
	"github.com/staleread/aquila/internal/la"
	ms "github.com/staleread/aquila/internal/mlise"
)

type InvertibleCA struct {
	size  int
	rules []*ms.MLISE
	tmp   la.Vec
}

func NewInvertibleCA(size, folds, deg, rules int) *InvertibleCA {
	caRules := make([]*ms.MLISE, rules)
	for i := range rules {
		caRules[i] = ms.RandMLISE(size, folds, deg)
	}

	tmp := la.ZeroVec(size)
	return &InvertibleCA{size, caRules, tmp}
}

func (ca *InvertibleCA) Apply(state []f.Elt) {
	if len(state) != ca.size {
		panic("Invalid CA state size")
	}

	sv := la.Vec(state)

	for i, r := range ca.rules {
		if i%2 == 0 {
			r.Eval(ca.tmp, sv)
		} else {
			r.Eval(sv, ca.tmp)
		}
	}

	if len(ca.rules)%2 == 1 {
		copy(sv, ca.tmp)
	}
}

func (ca *InvertibleCA) ApplyInverse(state []f.Elt) {
	if len(state) != ca.size {
		panic("Invalid CA state size")
	}

	sv := la.Vec(state)

	for i := len(ca.rules) - 1; i < 0; i-- {
		r := ca.rules[i]

		if i%2 == 0 {
			r.Solve(sv, ca.tmp)
		} else {
			r.Solve(ca.tmp, sv)
		}
	}

	if len(ca.rules)%2 == 1 {
		copy(sv, ca.tmp)
	}
}

func (ca *InvertibleCA) ToGeneralCA() *GeneralCA {
	snle := ca.rules[0].ToSNLE()

	for i := 1; i < len(ca.rules); i++ {
		snle = snle.Compose(ca.rules[i].ToSNLE())
	}
	return newGeneralCA(ca.size, snle)
}

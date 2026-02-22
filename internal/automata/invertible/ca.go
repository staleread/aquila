package invertible

import (
	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/canonical"
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/linalg"
)

type CA struct {
	size  int
	rules []*rule
	tmp   linalg.Vector
	arena *canonical.AllocationArena
}

func NewCA(size, folds, degree, rules int) *CA {
	caRules := make([]*rule, rules)
	arena := canonical.NewAllocationArena()

	for i := range rules {
		caRules[i] = randRule(size, folds, degree, arena)
	}

	tmp := linalg.ZeroVector(size)
	return &CA{size, caRules, tmp, arena}
}

func (ca *CA) Apply(state []field.Element) {
	if len(state) != ca.size {
		panic("Invalid CA state size")
	}

	sv := linalg.Vector(state)

	for i, r := range ca.rules {
		if i%2 == 0 {
			r.Apply(ca.tmp, sv)
		} else {
			r.Apply(sv, ca.tmp)
		}
	}

	if len(ca.rules)%2 == 1 {
		copy(sv, ca.tmp)
	}
}

func (ca *CA) ApplyInverse(state []field.Element) {
	if len(state) != ca.size {
		panic("Invalid CA state size")
	}

	sv := linalg.Vector(state)
	lastParity := (len(ca.rules) - 1) % 2

	for i := len(ca.rules) - 1; i >= 0; i-- {
		r := ca.rules[i]

		if i%2 == lastParity {
			r.ApplyInverse(ca.tmp, sv)
		} else {
			r.ApplyInverse(sv, ca.tmp)
		}
	}

	if lastParity == 0 {
		copy(sv, ca.tmp)
	}
}

func (ca *CA) General() *general.CA {
	n := len(ca.rules)
	result := ca.rules[n-1].toPolyset(ca.arena)

	for i := n - 2; i >= 0; i-- {
		target := ca.rules[i].toPolyset(ca.arena)
		cmpArena := canonical.NewComputationArena()
		prodCache := make(map[*field.Monomial]field.Polynomial)

		for j, p := range result {
			sum := field.Polynomial{}

			for m := range p.Monomials() {
				if cached, hit := prodCache[m]; hit {
					sum.AddTo(cached)
					continue
				}

				prod := make(field.Polynomial)
				first := true

				for s := range m.Subscripts() {
					if first {
						prod.AddTo(target[s])
						first = false
						continue
					}
					cmpArena.MulPolynomialBy(prod, target[s])
				}
				sum.AddTo(prod)
				prodCache[m] = prod
			}
			result[j] = sum
		}
	}
	return general.NewCA(ca.size, result)
}

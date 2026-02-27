package invertible

import (
	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/linalg"
	"github.com/staleread/aquila/internal/poly/sparse"
)

type CA struct {
	size      int
	rules     []*rule
	tmpState  linalg.Vector
}

func NewCA(size, folds, degree, rules int) *CA {
	caRules := make([]*rule, rules)

	for i := range rules {
		caRules[i] = randRule(size, folds, degree)
	}

	tmpState := linalg.ZeroVector(size)

	return &CA{size, caRules, tmpState}
}

func (ca *CA) Apply(state []field.Element) {
	if len(state) != ca.size {
		panic("Invalid CA state size")
	}

	sv := linalg.Vector(state)

	for i, r := range ca.rules {
		if i%2 == 0 {
			r.Apply(ca.tmpState, sv)
		} else {
			r.Apply(sv, ca.tmpState)
		}
	}

	if len(ca.rules)%2 == 1 {
		copy(sv, ca.tmpState)
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
			r.ApplyInverse(ca.tmpState, sv)
		} else {
			r.ApplyInverse(sv, ca.tmpState)
		}
	}

	if lastParity == 0 {
		copy(sv, ca.tmpState)
	}
}

func (ca *CA) ToGeneral() *general.CA {
	n := len(ca.rules)

	allocPool := sparse.NewMonomialInternPool()
	rule := ca.rules[n-1].toSparsePolyset(allocPool)

	for i := n - 2; i >= 0; i-- {
		target := ca.rules[i].toSparsePolyset(allocPool)

		prodCache := make(map[*sparse.Monomial]sparse.Polynomial)
		cmpPool := sparse.NewMonomialInternPool()

		for j, p := range rule {
			sum := sparse.Polynomial{}

			for m := range p.Monomials() {
				if cached, hit := prodCache[m]; hit {
					sum.AddTo(cached)
					continue
				}

				prod := make(sparse.Polynomial)
				first := true

				for s := range m.Subscripts() {
					if first {
						prod.AddTo(target[s])
						first = false
						continue
					}
					cmpPool.MulPolynomialBy(prod, target[s])
				}
				sum.AddTo(prod)
				prodCache[m] = prod
			}
			rule[j] = sum
		}
	}
	return general.NewCA(ca.size, rule.Compile())
}

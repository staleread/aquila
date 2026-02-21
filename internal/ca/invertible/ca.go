package invertible

import (
	"github.com/staleread/aquila/internal/ca/general"
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/linalg"
)

type CA struct {
	size  int
	rules []*rule
	tmp   linalg.Vector
}

func NewCA(size, folds, degree, rules int) *CA {
	caRules := make([]*rule, rules)

	for i := range rules {
		caRules[i] = randRule(size, folds, degree)
	}

	tmp := linalg.ZeroVector(size)
	return &CA{size, caRules, tmp}
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
	rule := ca.rules[n-1].general()

	for i := n - 2; i >= 0; i-- {
		rule.Compose(ca.rules[i].general())
	}
	return general.NewCA(ca.size, rule)
}

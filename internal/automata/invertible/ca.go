package invertible

import (
	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/linalg"
)

type CA struct {
	size     int
	rules    []*rule
	tmpState linalg.Vector
}

func NewCA(size, folds, degree, rules int) *CA {
	caRules := make([]*rule, rules)

	for i := range rules {
		caRules[i] = randRule(size, folds, degree)
	}

	tmpState := linalg.ZeroVector(size)

	return &CA{size, caRules, tmpState}
}

func (self *CA) Apply(state []field.Element) {
	if len(state) != self.size {
		panic("Invalid CA state size")
	}

	sv := linalg.Vector(state)

	for i, r := range self.rules {
		if i%2 == 0 {
			r.Apply(self.tmpState, sv)
		} else {
			r.Apply(sv, self.tmpState)
		}
	}

	if len(self.rules)%2 == 1 {
		copy(sv, self.tmpState)
	}
}

func (self *CA) ApplyInverse(state []field.Element) {
	if len(state) != self.size {
		panic("Invalid CA state size")
	}

	sv := linalg.Vector(state)
	lastParity := (len(self.rules) - 1) % 2

	for i := len(self.rules) - 1; i >= 0; i-- {
		r := self.rules[i]

		if i%2 == lastParity {
			r.ApplyInverse(self.tmpState, sv)
		} else {
			r.ApplyInverse(sv, self.tmpState)
		}
	}

	if lastParity == 0 {
		copy(sv, self.tmpState)
	}
}

func (self *CA) ToGeneral() *general.CA {
	n := len(self.rules)
	rule := self.rules[n-1].toSparsePolyset()

	for i := n - 2; i >= 0; i-- {
		rule.Compose(self.rules[i].toSparsePolyset())
	}
	return general.NewCA(self.size, rule.Compile())
}

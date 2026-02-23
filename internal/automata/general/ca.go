package general

import (
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/poly/compiled"
)

type CA struct {
	size int
	rule compiled.Polyset
	tmp  []field.Element
}

func NewCA(size int, rule compiled.Polyset) *CA {
	tmp := make([]field.Element, size)
	return &CA{size, rule, tmp}
}

func (ca *CA) Apply(state []field.Element) {
	if len(state) != ca.size {
		panic("Invalid CA state size")
	}

	copy(ca.tmp, state)
	ca.rule.Eval(state, ca.tmp)
}

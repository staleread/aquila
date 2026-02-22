package general

import "github.com/staleread/aquila/internal/field"

type Rule = field.Polyset

type CA struct {
	size int
	rule Rule
	tmp  []field.Element
}

func NewCA(size int, rule Rule) *CA {
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

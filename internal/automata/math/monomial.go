package math

import (
	"github.com/staleread/aquila/internal/automata/state"
)

type Monomial state.State

func NewMonomial(subs ...state.Subscript) Monomial {
	var b state.State
	for _, sub := range subs {
		b.SetAt(sub, 1)
	}
	return Monomial(b)
}

func CompareMonomials(a, b Monomial) int {
	return state.State(a).Compare(state.State(b))
}

func (m Monomial) Mul(other Monomial) Monomial {
	return Monomial(state.State(m).Or(state.State(other)))
}

func (m Monomial) Subscripts(dst []uint8) []uint8 {
	return state.State(m).Subscripts(dst)
}

func (m *Monomial) SetAt(sub state.Subscript, bit uint8) {
	(*state.State)(m).SetAt(sub, bit)
}

func (m Monomial) At(sub state.Subscript) uint8 {
	return state.State(m).At(sub)
}

func (m Monomial) Eval(b state.State) uint8 {
	if b.Contains(state.State(m)) {
		return 1
	}
	return 0
}

package math

import (
	"github.com/staleread/aquila/internal/automata/core"
)

type Monomial core.Block

func NewMonomial(subs ...core.Subscript) Monomial {
	var b core.Block
	for _, sub := range subs {
		b.SetAt(sub, 1)
	}
	return Monomial(b)
}

func CompareMonomials(a, b Monomial) int { return core.Block(a).Compare(core.Block(b)) }

func (m Monomial) Mul(other Monomial) Monomial          { return Monomial(core.Block(m).Or(core.Block(other))) }
func (m Monomial) Subscripts(dst []uint8) []uint8       { return core.Block(m).Subscripts(dst) }
func (m *Monomial) SetAt(sub core.Subscript, bit uint8) { (*core.Block)(m).SetAt(sub, bit) }
func (m Monomial) At(sub core.Subscript) uint8          { return core.Block(m).At(sub) }

func (m Monomial) Eval(b core.Block) uint8 {
	if b.Contains(core.Block(m)) {
		return 1
	}
	return 0
}

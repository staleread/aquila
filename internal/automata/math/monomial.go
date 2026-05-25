package math

import (
	"github.com/staleread/aquila/internal/automata/core"
)

type Monomial core.Block

func NewMonomial(words []core.Word) Monomial {
	var b core.Block
	copy(b[:], words)
	return Monomial(b)
}

func CompareMonomials(a, b Monomial) int          { return core.Block(a).Compare(core.Block(b)) }
func (m Monomial) FirstSubscript() int            { return core.Block(m).TrailingZeros() }
func (m Monomial) NextSubscript(start int) int    { return core.Block(m).TrailingZerosAfter(start) }
func (m Monomial) Mul(other Monomial) Monomial    { return Monomial(core.Block(m).Or(core.Block(other))) }
func (m Monomial) Subscripts(dst []uint8) []uint8 { return core.Block(m).Subscripts(dst) }

func (m Monomial) Eval(b *core.Block) core.Word {
	if b.Contains(core.Block(m)) {
		return 1
	}
	return 0
}

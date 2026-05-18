package math

import (
	"math/bits"

	"github.com/staleread/aquila/internal/automata/core"
)

type Monomial [3]core.Word

func NewMonomial(words []core.Word) Monomial {
	_ = words[2] // BCE hint

	return Monomial{
		words[0],
		words[1],
		words[2],
	}
}

func CompareMonomials(a, b Monomial) int {
	for i := range len(a) {
		if a[i] == b[i] {
			continue
		}
		if a[i] < b[i] {
			return -1
		}
		return 1
	}
	return 0
}

func (m Monomial) FirstSubscript() int {
	for i := range m {
		if m[i] != 0 {
			return i*core.BlockWordSize + bits.TrailingZeros32(m[i])
		}
	}
	return core.BlockSize
}

func (m Monomial) NextSubscript(start int) int {
	if start >= core.BlockSize {
		return core.BlockSize
	}

	startWordIdx := start / core.BlockWordSize

	if word := m[startWordIdx]; word != 0 {
		startSubIdx := start % core.BlockWordSize
		word >>= startSubIdx

		for j := range core.BlockWordSize - startSubIdx {
			if (word>>j)&1 == 1 {
				return start + j
			}
		}
	}

	for i := startWordIdx + 1; i < len(m); i++ {
		if m[i] != 0 {
			return i*core.BlockWordSize + bits.TrailingZeros32(m[i])
		}
	}
	return core.BlockSize
}

func (m Monomial) Degree() int {
	return bits.OnesCount32(m[0]) + bits.OnesCount32(m[1]) + bits.OnesCount32(m[2])
}

func (m Monomial) Mul(other Monomial) Monomial {
	return Monomial{
		m[0] | other[0],
		m[1] | other[1],
		m[2] | other[2],
	}
}

func (m Monomial) Eval(b *core.Block) core.Word {
	if b[0]&m[0] == m[0] &&
		b[1]&m[1] == m[1] &&
		b[2]&m[2] == m[2] {
		return 1
	}
	return 0
}

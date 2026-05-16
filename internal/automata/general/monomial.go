package general

import "github.com/staleread/aquila/internal/automata"

type Monomial [3]uint32

func CompareMonomials(a, b Monomial) int {
	if a[0] != b[0] {
		if a[0] < b[0] {
			return -1
		}
		return 1
	}
	if a[1] != b[1] {
		if a[1] < b[1] {
			return -1
		}
		return 1
	}
	if a[2] != b[2] {
		if a[2] < b[2] {
			return -1
		}
		return 1
	}
	return 0
}

func (m Monomial) Mul(other Monomial) Monomial {
	return Monomial{
		m[0] | other[0],
		m[1] | other[1],
		m[2] | other[2],
	}
}

func (m Monomial) Equals(other Monomial) bool {
	return m == other
}

func (m Monomial) Eval(b *automata.Block) uint8 {
	if b[0]&m[0] == m[0] &&
		b[1]&m[1] == m[1] &&
		b[2]&m[2] == m[2] {
		return 1
	}
	return 0
}

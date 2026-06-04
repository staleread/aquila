package math

type Monomial Bitset

var IdentityMonomial Monomial

func NewMonomial(subs ...Subscript) Monomial {
	var b Bitset
	for _, sub := range subs {
		b.SetAt(sub, 1)
	}
	return Monomial(b)
}

func CompareMonomials(a, b Monomial) int {
	return Bitset(a).Compare(Bitset(b))
}

func (m Monomial) Mul(other Monomial) Monomial {
	return Monomial(Bitset(m).Or(Bitset(other)))
}

func (m Monomial) Subscripts(dst []uint8) []uint8 {
	return Bitset(m).Subscripts(dst)
}

func (m *Monomial) SetAt(sub Subscript, bit uint8) {
	(*Bitset)(m).SetAt(sub, bit)
}

func (m Monomial) At(sub Subscript) uint8 {
	return Bitset(m).At(sub)
}

func (m Monomial) Eval(b Bitset) uint8 {
	if m == IdentityMonomial || b.Contains(Bitset(m)) {
		return 1
	}
	return 0
}

package gf2

import "iter"

type Polynomial map[subscriptHash][]Monomial

var ZeroPolynomial = Polynomial{}

var OnePolynomial = Polynomial{
	emptySubscriptHash: []Monomial{OneMonomial},
}

func NewPolynomial(monomials []Monomial) Polynomial {
	p := make(Polynomial, len(monomials))

	for _, m := range monomials {
		p.toggleMonomial(m)
	}
	return p
}

func RandPolynomial(deg int, maxSubscript Subscript) Polynomial {
	p := make(Polynomial, deg)

	for i := range deg {
		m := RandMonomial(deg-i, maxSubscript)
		p.addMonomialUnsafe(m)
	}
	return p
}

func (a Polynomial) Add(b Polynomial) Polynomial {
	pNew := make(Polynomial, len(a)+len(b))

	for m := range a.Monomials() {
		if !b.hasMonomial(m) {
			pNew.addMonomialUnsafe(m)
		}
	}

	for m := range b.Monomials() {
		if !a.hasMonomial(m) {
			pNew.addMonomialUnsafe(m)
		}
	}
	return pNew
}

func (a Polynomial) Mul(b Polynomial) Polynomial {
	pNew := make(Polynomial, len(a)*len(b))

	for ma := range a.Monomials() {
		for mb := range b.Monomials() {
			pNew.toggleMonomial(ma.Mul(mb))
		}
	}
	return pNew
}

func (p Polynomial) Eval(x []Element) Element {
	var sum Element = 0

	for m := range p.Monomials() {
		sum = Add(sum, m.Eval(x))
	}
	return sum
}

func (p Polynomial) Monomials() iter.Seq[Monomial] {
	return func(yield func(Monomial) bool) {
		for _, buck := range p {
			for _, m := range buck {
				if !yield(m) {
					return
				}
			}
		}
	}
}

func (p Polynomial) hasMonomial(m Monomial) bool {
	buck, ok := p[m.hash]

	if !ok {
		return false
	}

	for _, bm := range buck {
		if m.Equals(bm) {
			return true
		}
	}
	return false
}

func (p Polynomial) toggleMonomial(m Monomial) {
	hash := m.hash
	buck, ok := p[hash]

	if !ok {
		p[hash] = []Monomial{m}
		return
	}

	for i, bm := range buck {
		if m.Equals(bm) {
			if len(buck) == 1 {
				delete(p, m.hash)
				return
			}
			last := len(buck) - 1
			buck[i] = buck[last]

			p[m.hash] = buck[:last]
			return
		}
	}
	p[hash] = append(buck, m)
}

func (p Polynomial) addMonomialUnsafe(m Monomial) {
	hash := m.hash

	if buck, ok := p[hash]; ok {
		p[hash] = append(buck, m)
	} else {
		p[hash] = []Monomial{m}
	}
}

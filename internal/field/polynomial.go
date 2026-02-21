package field

import "iter"

type Polynomial map[uint32][]Monomial

var ZeroPolynomial = Polynomial{}

var OnePolynomial = Polynomial{
	0: []Monomial{OneMonomial},
}

func NewPolynomial(monoms []Monomial) Polynomial {
	p := make(Polynomial, len(monoms))

	for _, m := range monoms {
		p.toggleMonomial(m)
	}
	return p
}

func randPolynomial(degree int, maxSub Subscript) Polynomial {
	p := make(Polynomial, degree)

	for i := range degree {
		m := randMonomial(degree-i, maxSub)
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
		for _, bucket := range p {
			for _, m := range bucket {
				if !yield(m) {
					return
				}
			}
		}
	}
}

func (p Polynomial) hasMonomial(m Monomial) bool {
	bucket, ok := p[m.hash]

	if !ok {
		return false
	}

	for _, bm := range bucket {
		if m.Equals(bm) {
			return true
		}
	}
	return false
}

func (p Polynomial) toggleMonomial(m Monomial) {
	h := m.hash
	bucket, ok := p[h]

	if !ok {
		p[h] = []Monomial{m}
		return
	}

	for i, bm := range bucket {
		if m.Equals(bm) {
			last := len(bucket) - 1

			if last == 0 {
				delete(p, h)
				return
			}
			bucket[i] = bucket[last]
			p[h] = bucket[:last]
			return
		}
	}
	p[h] = append(bucket, m)
}

func (p Polynomial) addMonomialUnsafe(m Monomial) {
	h := m.hash

	if bucket, ok := p[h]; ok {
		p[h] = append(bucket, m)
	} else {
		p[h] = []Monomial{m}
	}
}

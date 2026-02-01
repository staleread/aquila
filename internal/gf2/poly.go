package gf2

import "iter"

type Poly map[symHash][]Monom

var ZeroPoly = Poly{}

var OnePoly = Poly{
	emptySymHash: []Monom{OneMonom},
}

func NewPoly(monoms []Monom) Poly {
	p := make(Poly, len(monoms))

	for _, m := range monoms {
		p.toggleMonom(m)
	}
	return p
}

func RandPoly(deg int, maxSym Sym) Poly {
	p := make(Poly, deg)

	for i := range deg {
		m := RandMonom(deg-i, maxSym)
		p.addMonomUnsafe(m)
	}
	return p
}

func (a Poly) Add(b Poly) Poly {
	pNew := make(Poly, len(a)+len(b))

	for m := range a.Monoms() {
		if !b.hasMonom(m) {
			pNew.addMonomUnsafe(m)
		}
	}

	for m := range b.Monoms() {
		if !a.hasMonom(m) {
			pNew.addMonomUnsafe(m)
		}
	}
	return pNew
}

func (a Poly) Mul(b Poly) Poly {
	pNew := make(Poly, len(a)*len(b))

	for ma := range a.Monoms() {
		for mb := range b.Monoms() {
			pNew.toggleMonom(ma.Mul(mb))
		}
	}
	return pNew
}

func (p Poly) Eval(x []Elt) Elt {
	var sum Elt = 0

	for m := range p.Monoms() {
		sum = Add(sum, m.Eval(x))
	}
	return sum
}

func (p Poly) Monoms() iter.Seq[Monom] {
	return func(yield func(Monom) bool) {
		for _, buck := range p {
			for _, m := range buck {
				if !yield(m) {
					return
				}
			}
		}
	}
}

func (p Poly) hasMonom(m Monom) bool {
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

func (p Poly) toggleMonom(m Monom) {
	hash := m.hash
	buck, ok := p[hash]

	if !ok {
		p[hash] = []Monom{m}
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

func (p Poly) addMonomUnsafe(m Monom) {
	hash := m.hash

	if buck, ok := p[hash]; ok {
		p[hash] = append(buck, m)
	} else {
		p[hash] = []Monom{m}
	}
}

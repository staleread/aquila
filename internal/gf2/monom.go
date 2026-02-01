package gf2

import (
	"iter"
	"maps"
)

const emptySymHash = 0

type symHash = uint32
type symSet = map[Sym]struct{}

type Monom struct {
	data symSet
	hash symHash
}

var OneMonom = Monom{
	data: nil,
	hash: emptySymHash,
}

func NewMonom(syms ...Sym) Monom {
	var hash symHash
	data := make(symSet, len(syms))

	for _, s := range syms {
		data[s] = struct{}{}
		hash ^= hashSym(s)
	}
	return Monom{data, hash}
}

func RandMonom(deg int, maxSym Sym) Monom {
	if Sym(deg) > maxSym {
		panic("Monom degree exceeds symbol range")
	}

	var hash symHash
	data := make(symSet, deg)

	rands := RandSyms(deg)
	randUp := maxSym - Sym(deg)

	for i := range Sym(deg) {
		s := rands[i] % (randUp + i + 1)

		if _, ok := data[s]; ok {
			s = randUp + i
		}
		data[s] = struct{}{}
		hash ^= hashSym(s)
	}
	return Monom{data, hash}
}

// https://stackoverflow.com/a/12996028
func hashSym(s Sym) uint32 {
	x := uint32(s)

	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = (x >> 16) ^ x

	return x
}

func (a Monom) Mul(b Monom) Monom {
	data := make(symSet, max(len(a.data), len(b.data)))
	hash := a.hash ^ b.hash

	for s := range a.data {
		data[s] = struct{}{}
	}

	for s := range b.data {
		if _, ok := a.data[s]; !ok {
			data[s] = struct{}{}
		} else {
			hash ^= hashSym(s)
		}
	}
	return Monom{data, hash}
}

func (a Monom) Equals(b Monom) bool {
	if len(a.data) != len(b.data) {
		return false
	}

	for s := range b.data {
		if _, ok := a.data[s]; !ok {
			return false
		}
	}
	return true
}

func (m Monom) Eval(x []Elt) Elt {
	var prod Elt = 1

	for s := range m.data {
		prod = Mul(prod, x[s])
	}
	return prod
}

func (m Monom) Size() int {
	return len(m.data)
}

func (m Monom) Syms() iter.Seq[Sym] {
	return maps.Keys(m.data)
}

func (m Monom) hasSym(s Sym) bool {
	_, ok := m.data[s]
	return ok
}

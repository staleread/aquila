package field

import (
	"iter"
	"maps"
)

type Monomial struct {
	data map[Subscript]struct{}
	hash uint32
}

var OneMonomial = Monomial{data: nil, hash: 0}

func NewMonomial(subs ...Subscript) Monomial {
	var hash uint32
	data := make(map[Subscript]struct{}, len(subs))

	for _, s := range subs {
		if _, ok := data[s]; !ok {
			data[s] = struct{}{}
			hash ^= hashSubscript(s)
		}
	}
	return Monomial{data, hash}
}

func randMonomial(degree int, maxSub Subscript) Monomial {
	if Subscript(degree) > maxSub {
		panic("Monomial degree exceeds subscript range")
	}

	var hash uint32
	data := make(map[Subscript]struct{}, degree)

	rands := RandSubscripts(degree)
	randUp := maxSub - Subscript(degree)

	for i := range Subscript(degree) {
		s := rands[i] % (randUp + i + 1)

		if _, ok := data[s]; ok {
			s = randUp + i
		}
		data[s] = struct{}{}
		hash ^= hashSubscript(s)
	}
	return Monomial{data, hash}
}

func (a Monomial) Mul(b Monomial) Monomial {
	data := make(map[Subscript]struct{}, max(len(a.data), len(b.data)))
	hash := a.hash ^ b.hash

	for s := range a.data {
		data[s] = struct{}{}
	}

	for s := range b.data {
		if _, ok := a.data[s]; !ok {
			data[s] = struct{}{}
		} else {
			// x*x = x in GF(2), so cancel from hash and omit from data
			hash ^= hashSubscript(s)
		}
	}
	return Monomial{data, hash}
}

func (m Monomial) Syms() iter.Seq[Subscript] {
	return maps.Keys(m.data)
}

func (m Monomial) Eval(x []Element) Element {
	var prod Element = 1

	for s := range m.data {
		prod = Mul(prod, x[s])
	}
	return prod
}

func (a Monomial) Equals(b Monomial) bool {
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

// https://stackoverflow.com/a/12996028
func hashSubscript(s Subscript) uint32 {
	x := uint32(s)

	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = (x >> 16) ^ x

	return x
}

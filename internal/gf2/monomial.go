package gf2

import (
	"iter"
	"maps"
)

const emptySubscriptHash = 0

type subscriptHash = uint32
type subscriptSet = map[Subscript]struct{}

type Monomial struct {
	data subscriptSet
	hash subscriptHash
}

var OneMonomial = Monomial{
	data: nil,
	hash: emptySubscriptHash,
}

func NewMonomial(subs ...Subscript) Monomial {
	var hash subscriptHash
	data := make(subscriptSet, len(subs))

	for _, s := range subs {
		data[s] = struct{}{}
		hash ^= hashSubscript(s)
	}
	return Monomial{data, hash}
}

func RandMonomial(degree int, maxSub Subscript) Monomial {
	if Subscript(degree) > maxSub {
		panic("Monomial degree exceeds subscript range")
	}

	var hash subscriptHash
	data := make(subscriptSet, degree)

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

// https://stackoverflow.com/a/12996028
func hashSubscript(s Subscript) uint32 {
	x := uint32(s)

	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = ((x >> 16) ^ x) * 0x45d9f3b
	x = (x >> 16) ^ x

	return x
}

func (a Monomial) Mul(b Monomial) Monomial {
	data := make(subscriptSet, max(len(a.data), len(b.data)))
	hash := a.hash ^ b.hash

	for s := range a.data {
		data[s] = struct{}{}
	}

	for s := range b.data {
		if _, ok := a.data[s]; !ok {
			data[s] = struct{}{}
		} else {
			hash ^= hashSubscript(s)
		}
	}
	return Monomial{data, hash}
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

func (m Monomial) Eval(x []Element) Element {
	var prod Element = 1

	for s := range m.data {
		prod = Mul(prod, x[s])
	}
	return prod
}

func (m Monomial) Size() int {
	return len(m.data)
}

func (m Monomial) Subscripts() iter.Seq[Subscript] {
	return maps.Keys(m.data)
}

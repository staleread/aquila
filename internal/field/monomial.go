package field

import (
	"iter"
	"maps"
)

type MonomialPtr = uint64
type monomial map[Subscript]struct{}

var monomialCache = map[MonomialPtr]monomial{
	0: nil,
}

func NewMonomial(subs ...Subscript) MonomialPtr {
	var ptr MonomialPtr
	monom := make(monomial, len(subs))

	for _, s := range subs {
		if _, ok := monom[s]; !ok {
			monom[s] = struct{}{}
			ptr ^= hashSubscript(s)
		}
	}
	return safeMonomial(monom, ptr)
}

func SubscriptsOf(mPtr MonomialPtr) iter.Seq[Subscript] {
	m := monomialCache[mPtr]
	return maps.Keys(m)
}

func MulMonomials(aPtr, bPtr MonomialPtr) MonomialPtr {
	a, b := monomialCache[aPtr], monomialCache[bPtr]

	ptr := aPtr ^ bPtr
	prod := make(monomial, max(len(a), len(b)))

	for s := range a {
		prod[s] = struct{}{}
	}

	for s := range b {
		if _, ok := a[s]; !ok {
			prod[s] = struct{}{}
		} else {
			ptr ^= hashSubscript(s)
		}
	}
	return safeMonomial(prod, ptr)
}

func randMonomial(degree int, maxSub Subscript) MonomialPtr {
	if Subscript(degree) > maxSub {
		panic("Monomial degree exceeds subscript range")
	}

	var ptr MonomialPtr
	monom := make(monomial, degree)

	rands := RandSubscripts(degree)
	randUp := maxSub - Subscript(degree)

	for i := range Subscript(degree) {
		s := rands[i] % (randUp + i + 1)

		if _, ok := monom[s]; ok {
			s = randUp + i
		}
		monom[s] = struct{}{}
		ptr ^= hashSubscript(s)
	}
	return safeMonomial(monom, ptr)
}

func evalMonomial(mPtr MonomialPtr, x []Element) Element {
	m := monomialCache[mPtr]

	var prod Element = 1

	for s := range m {
		prod = Mul(prod, x[s])
	}
	return prod
}

func safeMonomial(candidate monomial, suggestedPtr MonomialPtr) MonomialPtr {
	other, ok := monomialCache[suggestedPtr]

	if !ok {
		monomialCache[suggestedPtr] = candidate
		return suggestedPtr
	}

	if other.equals(candidate) {
		return suggestedPtr
	}

	newPtr := suggestedPtr
	for {
		newPtr++
		if _, ok := monomialCache[newPtr]; !ok {
			monomialCache[newPtr] = candidate
			return newPtr
		}

		if newPtr == suggestedPtr {
			panic("Monomial cache is out of hash space")
		}
	}
}

// https://dl.acm.org/doi/10.1145/2714064.2660195
func hashSubscript(s Subscript) MonomialPtr {
	x := MonomialPtr(s)

	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return x
}

func (a monomial) equals(b monomial) bool {
	if len(a) != len(b) {
		return false
	}

	for s := range b {
		if _, ok := a[s]; !ok {
			return false
		}
	}
	return true
}

package field

import (
	"iter"
	"maps"
)

type PolynomialPtr = uint64
type polynomial map[MonomialPtr]struct{}

var polynomialCache = map[PolynomialPtr]polynomial{
	0: nil,
}

func NewPolynomial(monomPtrs []MonomialPtr) PolynomialPtr {
	var ptr PolynomialPtr
	poly := make(polynomial, len(monomPtrs))

	for _, mPtr := range monomPtrs {
		if _, ok := poly[mPtr]; !ok {
			poly[mPtr] = struct{}{}
		} else {
			delete(poly, mPtr)
		}
		ptr ^= mPtr
	}
	return safePolynomial(poly, ptr)
}

func AddPolynomials(aPtr, bPtr PolynomialPtr) PolynomialPtr {
	a, b := polynomialCache[aPtr], polynomialCache[bPtr]

	ptr := aPtr ^ bPtr
	sum := make(polynomial, len(a)+len(b))

	for maPtr := range maps.Keys(a) {
		if _, ok := b[maPtr]; !ok {
			sum[maPtr] = struct{}{}
		}
	}

	for mbPtr := range maps.Keys(b) {
		if _, ok := a[mbPtr]; !ok {
			sum[mbPtr] = struct{}{}
		}
	}
	return safePolynomial(sum, ptr)
}

func MulPolynomials(aPtr, bPtr PolynomialPtr) PolynomialPtr {
	a, b := polynomialCache[aPtr], polynomialCache[bPtr]

	var ptr PolynomialPtr
	prod := make(polynomial, len(a)*len(b))

	for maPtr := range maps.Keys(a) {
		for mbPtr := range maps.Keys(b) {
			mProdPtr := MulMonomials(maPtr, mbPtr)
			ptr ^= mProdPtr

			if _, ok := prod[mProdPtr]; !ok {
				prod[mProdPtr] = struct{}{}
			} else {
				delete(prod, mProdPtr)
			}
		}
	}
	return safePolynomial(prod, ptr)
}

func MonomialsOf(pPtr PolynomialPtr) iter.Seq[MonomialPtr] {
	p := polynomialCache[pPtr]
	return maps.Keys(p)
}

func evalPolynomial(pPtr PolynomialPtr, x []Element) Element {
	p := polynomialCache[pPtr]

	var sum Element = 0

	for m := range maps.Keys(p) {
		sum = Add(sum, evalMonomial(m, x))
	}
	return sum
}

func randPolynomial(degree int, maxSub Subscript) PolynomialPtr {
	var ptr PolynomialPtr
	poly := make(polynomial, degree)

	for i := range degree {
		mPtr := randMonomial(degree-i, maxSub)
		poly[mPtr] = struct{}{}

		ptr ^= mPtr
	}
	return safePolynomial(poly, ptr)
}

func safePolynomial(candidate polynomial, suggestedPtr PolynomialPtr) PolynomialPtr {
	other, ok := polynomialCache[suggestedPtr]

	if !ok {
		polynomialCache[suggestedPtr] = candidate
		return suggestedPtr
	}

	if other.Equals(candidate) {
		return suggestedPtr
	}

	newPtr := suggestedPtr
	for {
		newPtr++
		if _, ok := polynomialCache[newPtr]; !ok {
			polynomialCache[newPtr] = candidate
			return newPtr
		}

		if newPtr == suggestedPtr {
			panic("Polynomial cache is out of hash space")
		}
	}
}

func (p polynomial) Equals(other polynomial) bool {
	if len(p) != len(other) {
		return false
	}

	for mPtr := range p {
		if _, ok := other[mPtr]; !ok {
			return false
		}
	}
	return true
}

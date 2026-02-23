package sparse

import (
	"fmt"
	"github.com/staleread/aquila/internal/poly"
	"iter"
	"maps"
	"strings"
)

type MonomialHash = uint64
type subscriptSet = map[poly.Subscript]struct{}

type Monomial struct {
	Hash MonomialHash
	data subscriptSet
}

func NewMonomial(subs ...poly.Subscript) Monomial {
	var hash MonomialHash
	data := make(subscriptSet, len(subs))

	for _, s := range subs {
		if _, ok := data[s]; !ok {
			data[s] = struct{}{}
			hash ^= hashSubscript(s)
		}
	}
	return Monomial{hash, data}
}

func hashSubscript(s poly.Subscript) MonomialHash {
	x := MonomialHash(s + 1)
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}

func (monom Monomial) Subscripts() iter.Seq[poly.Subscript] {
	return maps.Keys(monom.data)
}

func (a Monomial) Mul(b Monomial) Monomial {
	hash := a.Hash ^ b.Hash
	data := make(subscriptSet, max(len(a.data), len(b.data)))

	for s := range a.data {
		data[s] = struct{}{}
	}

	for s := range b.data {
		if _, ok := a.data[s]; !ok {
			data[s] = struct{}{}
		} else {
			// x_i^2 = x_i in GF(2): duplicate subscript is already in data,
			// un-XOR its hash contribution since it was counted once in a.Hash
			// and once in b.Hash but should only appear once in the product.
			hash ^= hashSubscript(s)
		}
	}
	return Monomial{hash, data}
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

func (monom Monomial) String() string {
	sb := strings.Builder{}

	sb.WriteByte('{')

	isFirstSub := true
	for s := range monom.Subscripts() {
		if !isFirstSub {
			sb.WriteString(", ")
		}
		isFirstSub = false

		fmt.Fprintf(&sb, "%d", s)
	}
	sb.WriteByte('}')

	return sb.String()
}

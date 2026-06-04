//go:build block8

package math

import "math/bits"

const (
	BitsetSize  = 8
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset uint8

func NewBitset(subs ...Subscript) Bitset {
	var b Bitset
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

func (s *Bitset) Read(src []byte) {
	_ = src[0] // BCE check
	*s = Bitset(src[0])
}

func (s Bitset) Write(dst []byte) {
	_ = dst[0] // BCE check
	dst[0] = uint8(s)
}

func (s Bitset) At(idx Subscript) uint8 {
	return uint8((s >> idx) & 1)
}

func (s *Bitset) SetAt(idx Subscript, bit uint8) {
	*s = (*s &^ (Bitset(1) << idx)) | (Bitset(bit&1) << idx)
}

func (s Bitset) Compare(other Bitset) int {
	if s == other {
		return 0
	}
	if s < other {
		return -1
	}
	return 1
}

func (s Bitset) Or(other Bitset) Bitset {
	return s | other
}

func (s *Bitset) XorWith(other Bitset) {
	*s ^= other
}

func (s Bitset) Contains(other Bitset) bool {
	return s&other == other
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	w := uint8(s)
	for w != 0 {
		tz := bits.TrailingZeros8(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}

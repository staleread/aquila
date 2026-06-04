//go:build block64

package math

import (
	"encoding/binary"
	"math/bits"
)

const (
	BitsetSize  = 64
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset uint64

func NewBitset(subs ...Subscript) Bitset {
	var b Bitset
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

func (s *Bitset) Read(src []byte) {
	_ = src[7] // BCE check
	*s = Bitset(binary.LittleEndian.Uint64(src[0:8]))
}

func (s Bitset) Write(dst []byte) {
	_ = dst[7] // BCE check
	binary.LittleEndian.PutUint64(dst[0:8], uint64(s))
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
	w := uint64(s)
	for w != 0 {
		tz := bits.TrailingZeros64(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}

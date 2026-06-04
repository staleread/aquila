//go:build block128

package state

import (
	"encoding/binary"
	"math/bits"
)

const (
	BitsetSize  = 128
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset [2]uint64

func NewBitset(subs ...Subscript) Bitset {
	var s Bitset
	for _, i := range subs {
		s.SetAt(i, 1)
	}
	return s
}

func (s *Bitset) Read(src []byte) {
	_ = src[15] // BCE check

	s[0] = binary.LittleEndian.Uint64(src[0:8])
	s[1] = binary.LittleEndian.Uint64(src[8:16])
}

func (s *Bitset) Write(dst []byte) {
	_ = dst[15] // BCE check

	binary.LittleEndian.PutUint64(dst[0:8], s[0])
	binary.LittleEndian.PutUint64(dst[8:16], s[1])
}

func (s Bitset) At(sub Subscript) uint8 {
	return uint8((s[sub/64] >> (sub % 64)) & 1)
}

func (s *Bitset) SetAt(sub Subscript, bit uint8) {
	wordIdx := sub / 64
	shift := sub % 64

	s[wordIdx] = (s[wordIdx] &^ (uint64(1) << shift)) | uint64(bit&1)<<shift
}

func (s Bitset) Compare(other Bitset) int {
	for i := range len(s) {
		if s[i] == other[i] {
			continue
		}
		if s[i] < other[i] {
			return -1
		}
		return 1
	}
	return 0
}

func (s Bitset) Or(other Bitset) Bitset {
	return Bitset{
		s[0] | other[0],
		s[1] | other[1],
	}
}

func (s *Bitset) XorWith(other Bitset) {
	s[0] ^= other[0]
	s[1] ^= other[1]
}

func (s Bitset) Contains(other Bitset) bool {
	return s[0]&other[0] == other[0] && s[1]&other[1] == other[1]
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	for i, w := range s {
		wordOffset := i * 64
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			dst = append(dst, Subscript(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
}

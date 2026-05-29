package core

import (
	"encoding/binary"
	"math/bits"
)

type Block16 uint16

func (b *Block16) Read(src []byte) {
	_ = src[1] // BCE check
	*b = Block16(binary.LittleEndian.Uint16(src[0:2]))
}

func (b Block16) Write(dst []byte) {
	_ = dst[1] // BCE check
	binary.LittleEndian.PutUint16(dst[0:2], uint16(b))
}

func (b Block16) At(idx Subscript) uint8 {
	return uint8((b >> idx) & 1)
}

func (b *Block16) SetAt(idx Subscript, bit uint8) {
	*b = (*b &^ (Block16(1) << idx)) | (Block16(bit&1) << idx)
}

func (b Block16) Compare(other Block16) int {
	if b == other {
		return 0
	}
	if b < other {
		return -1
	}
	return 1
}

func (b Block16) Or(other Block16) Block16 {
	return b | other
}

func (b Block16) Contains(other Block16) bool {
	return b&other == other
}

func (b Block16) Subscripts(dst []Subscript) []Subscript {
	w := uint16(b)
	for w != 0 {
		tz := bits.TrailingZeros16(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}

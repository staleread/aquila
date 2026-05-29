package core

import (
	"encoding/binary"
	"math/bits"
)

type Block32 uint32

func (b *Block32) Read(src []byte) {
	_ = src[3] // BCE check
	*b = Block32(binary.LittleEndian.Uint32(src[0:4]))
}

func (b Block32) Write(dst []byte) {
	_ = dst[3] // BCE check
	binary.LittleEndian.PutUint32(dst[0:4], uint32(b))
}

func (b Block32) At(idx Subscript) uint8 {
	return uint8((b >> idx) & 1)
}

func (b *Block32) SetAt(idx Subscript, bit uint8) {
	*b = (*b &^ (Block32(1) << idx)) | (Block32(bit&1) << idx)
}

func (b Block32) Compare(other Block32) int {
	if b == other {
		return 0
	}
	if b < other {
		return -1
	}
	return 1
}

func (b Block32) Or(other Block32) Block32 {
	return b | other
}

func (b Block32) Contains(other Block32) bool {
	return b&other == other
}

func (b Block32) Subscripts(dst []Subscript) []Subscript {
	w := uint32(b)
	for w != 0 {
		tz := bits.TrailingZeros32(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}

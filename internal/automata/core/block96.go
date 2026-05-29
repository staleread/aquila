package core

import (
	"encoding/binary"
	"math/bits"
)

type Block96 [3]uint32

func (b *Block96) Read(src []byte) {
	_ = src[11] // BCE check

	b[0] = binary.LittleEndian.Uint32(src[0:4])
	b[1] = binary.LittleEndian.Uint32(src[4:8])
	b[2] = binary.LittleEndian.Uint32(src[8:12])
}

func (b *Block96) Write(dst []byte) {
	_ = dst[11] // BCE check

	binary.LittleEndian.PutUint32(dst[0:4], b[0])
	binary.LittleEndian.PutUint32(dst[4:8], b[1])
	binary.LittleEndian.PutUint32(dst[8:12], b[2])
}

func (b Block96) At(idx Subscript) uint8 {
	return uint8((b[idx/32] >> (idx % 32)) & 1)
}

func (b *Block96) SetAt(idx Subscript, bit uint8) {
	wordIdx := idx / 32
	shift := idx % 32

	b[wordIdx] = (b[wordIdx] &^ (uint32(1) << shift)) | uint32(bit&1)<<shift
}

func (b Block96) Compare(other Block96) int {
	for i := range len(b) {
		if b[i] == other[i] {
			continue
		}
		if b[i] < other[i] {
			return -1
		}
		return 1
	}
	return 0
}

func (b Block96) Or(other Block96) Block96 {
	return Block96{
		b[0] | other[0],
		b[1] | other[1],
		b[2] | other[2],
	}
}

func (b Block96) Contains(other Block96) bool {
	return b[0]&other[0] == other[0] && b[1]&other[1] == other[1] && b[2]&other[2] == other[2]
}

func (b Block96) Subscripts(dst []Subscript) []Subscript {
	for i, w := range b {
		wordOffset := i * 32
		for w != 0 {
			tz := bits.TrailingZeros32(w)
			dst = append(dst, Subscript(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
}

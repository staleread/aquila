package core

import (
	"encoding/binary"
	"math/bits"
)

type Block128 [2]uint64

func (b *Block128) Read(src []byte) {
	_ = src[15] // BCE check

	b[0] = binary.LittleEndian.Uint64(src[0:8])
	b[1] = binary.LittleEndian.Uint64(src[8:16])
}

func (b *Block128) Write(dst []byte) {
	_ = dst[15] // BCE check

	binary.LittleEndian.PutUint64(dst[0:8], b[0])
	binary.LittleEndian.PutUint64(dst[8:16], b[1])
}

func (b Block128) At(sub Subscript) uint8 {
	return uint8((b[sub/64] >> (sub % 64)) & 1)
}

func (b *Block128) SetAt(sub Subscript, bit uint8) {
	wordIdx := sub / 64
	shift := sub % 64

	b[wordIdx] = (b[wordIdx] &^ (uint64(1) << shift)) | uint64(bit&1)<<shift
}

func (b Block128) Compare(other Block128) int {
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

func (b Block128) Or(other Block128) Block128 {
	return Block128{
		b[0] | other[0],
		b[1] | other[1],
	}
}

func (b Block128) Contains(other Block128) bool {
	return b[0]&other[0] == other[0] && b[1]&other[1] == other[1]
}

func (b Block128) Subscripts(dst []Subscript) []Subscript {
	for i, w := range b {
		wordOffset := i * 64
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			dst = append(dst, Subscript(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
}

package core

import (
	"encoding/binary"
	"math/bits"
)

type Block96 [3]uint32

func LoadBlock96(src []byte) *Block96 {
	_ = src[11] // BCE check

	return &Block96{
		binary.LittleEndian.Uint32(src[0:4]),
		binary.LittleEndian.Uint32(src[4:8]),
		binary.LittleEndian.Uint32(src[8:12]),
	}
}

func (b *Block96) WriteTo(dst []byte) {
	_ = dst[11] // BCE check

	binary.LittleEndian.PutUint32(dst[0:4], b[0])
	binary.LittleEndian.PutUint32(dst[4:8], b[1])
	binary.LittleEndian.PutUint32(dst[8:12], b[2])
}

func (b *Block96) At(idx int) uint32 {
	return (b[idx/32] >> (idx % 32)) & 1
}

func (b *Block96) SetAt(idx int, bit uint32) {
	wordIdx := idx / 32
	shift := idx % 32

	b[wordIdx] = (b[wordIdx] &^ (uint32(1) << shift)) | (bit&1)<<shift
}

func (b Block96) Subscripts(dst []uint8) []uint8 {
	for i, w := range b {
		wordOffset := i * 32
		for w != 0 {
			tz := bits.TrailingZeros32(w)
			dst = append(dst, uint8(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
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

func (b Block96) TrailingZeros() int {
	for i := range b {
		if b[i] != 0 {
			return i*32 + bits.TrailingZeros32(b[i])
		}
	}
	return 96
}

func (b Block96) TrailingZerosAfter(start int) int {
	if start >= 96 {
		return 96
	}

	startWordIdx := start / 32

	if word := b[startWordIdx]; word != 0 {
		startSubIdx := start % 32
		word >>= startSubIdx

		for j := range 32 - startSubIdx {
			if (word>>j)&1 == 1 {
				return start + j
			}
		}
	}

	for i := startWordIdx + 1; i < len(b); i++ {
		if b[i] != 0 {
			return i*32 + bits.TrailingZeros32(b[i])
		}
	}
	return 96
}

func (b Block96) OnesCount() int {
	return bits.OnesCount32(b[0]) + bits.OnesCount32(b[1]) + bits.OnesCount32(b[2])
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

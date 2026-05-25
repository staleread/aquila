package core

import (
	"encoding/binary"
	"math/bits"
)

type Block128 [2]uint64

func LoadBlock128(src []byte) *Block128 {
	_ = src[15] // BCE check

	return &Block128{
		binary.LittleEndian.Uint64(src[0:8]),
		binary.LittleEndian.Uint64(src[8:16]),
	}
}

func (b *Block128) WriteTo(dst []byte) {
	_ = dst[15] // BCE check

	binary.LittleEndian.PutUint64(dst[0:8], b[0])
	binary.LittleEndian.PutUint64(dst[8:16], b[1])
}

func (b *Block128) At(idx int) uint64 {
	return (b[idx/64] >> (idx % 64)) & 1
}

func (b *Block128) SetAt(idx int, bit uint64) {
	wordIdx := idx / 64
	shift := idx % 64

	b[wordIdx] = (b[wordIdx] &^ (uint64(1) << shift)) | (bit&1)<<shift
}

func (b Block128) Subscripts(dst []uint8) []uint8 {
	for i, w := range b {
		wordOffset := i * 64
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			dst = append(dst, uint8(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
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

func (b Block128) TrailingZeros() int {
	for i := range b {
		if b[i] != 0 {
			return i*64 + bits.TrailingZeros64(b[i])
		}
	}
	return 128
}

func (b Block128) TrailingZerosAfter(start int) int {
	if start >= 128 {
		return 128
	}

	startWordIdx := start / 64

	if word := b[startWordIdx]; word != 0 {
		startSubIdx := start % 64
		word >>= startSubIdx

		for j := range 64 - startSubIdx {
			if (word>>j)&1 == 1 {
				return start + j
			}
		}
	}

	for i := startWordIdx + 1; i < len(b); i++ {
		if b[i] != 0 {
			return i*64 + bits.TrailingZeros64(b[i])
		}
	}
	return 128
}

func (b Block128) OnesCount() int {
	return bits.OnesCount64(b[0]) + bits.OnesCount64(b[1])
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

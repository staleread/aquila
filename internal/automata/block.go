package automata

import (
	"encoding/binary"
)

type Block [2]uint64
type Bit uint8

func LoadBlock(src []byte) *Block {
	_ = src[BlockBytes-1] // BCE check

	return &Block{
		binary.LittleEndian.Uint64(src[0:8]),
		binary.LittleEndian.Uint64(src[8:16]),
	}
}

func (b *Block) WriteTo(dst []byte) {
	_ = dst[BlockBytes-1] // BCE check

	binary.LittleEndian.PutUint64(dst[0:8], b[0])
	binary.LittleEndian.PutUint64(dst[8:16], b[1])
}

func (b *Block) At(idx uint8) Bit {
	return Bit((b[idx>>6] >> (idx & 63)) & 1)
}

func (b *Block) SetAt(idx uint8, bit Bit) {
	wordIdx := idx >> 6
	shift := idx & 63

	b[wordIdx] = (b[wordIdx] &^ (uint64(1) << shift)) | uint64(bit&1)<<shift
}

package automata

import (
	"encoding/binary"
)

type Word = uint32
type Block [BlockWords]Word

type Bit uint8

func LoadBlock(src []byte) *Block {
	_ = src[BlockBytes-1] // BCE check

	return &Block{
		binary.LittleEndian.Uint32(src[0:4]),
		binary.LittleEndian.Uint32(src[4:8]),
		binary.LittleEndian.Uint32(src[8:12]),
	}
}

func (b *Block) WriteTo(dst []byte) {
	_ = dst[BlockBytes-1] // BCE check

	binary.LittleEndian.PutUint32(dst[0:4], b[0])
	binary.LittleEndian.PutUint32(dst[4:8], b[1])
	binary.LittleEndian.PutUint32(dst[8:12], b[2])
}

func (b *Block) At(idx uint8) Bit {
	return Bit((b[idx/32] >> (idx % 32)) & 1)
}

func (b *Block) SetAt(idx uint8, bit Bit) {
	wordIdx := idx / 32
	shift := idx % 32

	b[wordIdx] = (b[wordIdx] &^ (Word(1) << shift)) | Word(bit&1)<<shift
}

package invertible

type Block []uint64
type Bit = uint8

func (b Block) At(idx uint8) Bit {
	return Bit((b[idx>>6] >> (idx & 63)) & 1)
}

func (b Block) SetAt(bit Bit, idx uint8) {
	wordIdx := idx >> 6
	shift := idx & 63

	b[wordIdx] = (b[wordIdx] &^ (uint64(1) << shift)) | uint64(bit&1)<<shift
}

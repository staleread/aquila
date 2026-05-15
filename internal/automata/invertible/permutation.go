package invertible

import "io"

type Permutation []Subscript

func (p Permutation) Generate(rnd io.Reader, tmp []byte) error {
	const n Subscript = VectorSize

	for i := range n {
		p[i] = i
	}

	if _, err := io.ReadFull(rnd, tmp); err != nil {
		return err
	}

	// Fisher–Yates shuffle
	for i := range n - 2 {
		j := tmp[i]%(n-i+1) + i
		p[i], p[j] = p[j], p[i]
	}
	return nil
}

func (p Permutation) Gather(b Block, foldIdx int) Vector {
	var res Vector
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		bit := b.At(p[offset+i])
		res |= Vector(bit) << i
	}
	return res
}

func (p Permutation) Scatter(b Block, foldIdx int, v Vector) {
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		bit := Bit((v >> i) & 1)
		b.SetAt(p[offset+i], bit)
	}
}

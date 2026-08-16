package math

import "io"

const (
	PermutationSize  = BitsetSize
	PermutationBytes = PermutationSize
)

type Permutation struct {
	Data []Subscript
}

func NewPermutation(arena []byte) *Permutation {
	return &Permutation{
		Data: arena,
	}
}

func (p *Permutation) Generate(rnd io.Reader) error {
	const n = PermutationSize

	for i := range n {
		p.Data[i] = Subscript(i)
	}

	for i := range n - 1 {
		j, err := readUniformByte(rnd, n-i)

		if err != nil {
			return err
		}

		p.Data[i], p.Data[i+j] = p.Data[i+j], p.Data[i]
	}
	return nil
}

func (p *Permutation) Gather(b *Bitset, foldIdx int) Vector {
	var res Vector
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := p.Data[offset+i]
		res |= Vector(b.At(idx)) << i
	}
	return res
}

func (p *Permutation) Scatter(b *Bitset, foldIdx int, v Vector) {
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := p.Data[offset+i]
		bit := uint8((v >> i) & 1)
		b.SetAt(idx, bit)
	}
}

func readUniformByte(rnd io.Reader, bound int) (int, error) {
	limit := 256 - 256%bound
	var buf [1]byte

	for {
		if _, err := io.ReadFull(rnd, buf[:]); err != nil {
			return 0, err
		}
		if int(buf[0]) < limit {
			return int(buf[0]) % bound, nil
		}
	}
}

package invertible

import (
	"io"

	"github.com/staleread/aquila/internal/automata"
)

type Permutation struct {
	data []Subscript
}

func NewPermutation(arena []byte) *Permutation {
	return &Permutation{
		data: arena,
	}
}

func (p *Permutation) Generate(rnd io.Reader, tmp []byte) error {
	const n Subscript = PermutationSize

	for i := range n {
		p.data[i] = i
	}

	if _, err := io.ReadFull(rnd, tmp); err != nil {
		return err
	}

	// Fisher–Yates shuffle
	for i := range n - 1 {
		j := tmp[i]%(n-i) + i
		p.data[i], p.data[j] = p.data[j], p.data[i]
	}
	return nil
}

func (p *Permutation) Gather(b *automata.Block, foldIdx int) Vector {
	var res Vector
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := int(p.data[offset+i])
		res |= Vector(b.At(idx)) << i
	}
	return res
}

func (p *Permutation) Scatter(b *automata.Block, foldIdx int, v Vector) {
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := int(p.data[offset+i])
		bit := automata.Word((v >> i) & 1)
		b.SetAt(idx, bit)
	}
}

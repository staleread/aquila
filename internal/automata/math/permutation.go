package math

import (
	"io"

	"github.com/staleread/aquila/internal/automata/core"
)

const PermutationSize = core.BlockSize
const PermutationBytes = PermutationSize

type Permutation struct {
	Data []core.Subscript
}

func NewPermutation(arena []byte) *Permutation {
	return &Permutation{
		Data: arena,
	}
}

func (p *Permutation) Generate(rnd io.Reader, tmp []byte) error {
	const n core.Subscript = PermutationSize

	for i := range n {
		p.Data[i] = i
	}

	if _, err := io.ReadFull(rnd, tmp); err != nil {
		return err
	}

	// Fisher–Yates shuffle
	for i := range n - 1 {
		j := tmp[i]%(n-i) + i
		p.Data[i], p.Data[j] = p.Data[j], p.Data[i]
	}
	return nil
}

func (p *Permutation) Gather(b *core.Block, foldIdx int) Vector {
	var res Vector
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := p.Data[offset+i]
		res |= Vector(b.At(idx)) << i
	}
	return res
}

func (p *Permutation) Scatter(b *core.Block, foldIdx int, v Vector) {
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := p.Data[offset+i]
		bit := uint8((v >> i) & 1)
		b.SetAt(idx, bit)
	}
}

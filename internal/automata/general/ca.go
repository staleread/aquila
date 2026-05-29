package general

import (
	"encoding/binary"
	"io"

	"github.com/staleread/aquila/internal/automata/core"
	"github.com/staleread/aquila/internal/automata/math"
)

type CA struct {
	Arena   []math.Monomial
	Offsets [core.BlockSize - 1]uint32
}

func (ca *CA) GetPolynomial(idx int) []math.Monomial {
	start := uint32(0)
	if idx > 0 {
		start = ca.Offsets[idx-1]
	}
	end := uint32(len(ca.Arena))
	if idx < core.BlockSize-1 {
		end = ca.Offsets[idx]
	}
	return ca.Arena[start:end]
}

func (ca *CA) Apply(dst, src []byte) {
	var srcBlock core.Block
	srcBlock.Read(src)
	var dstBlock core.Block

	for i := range core.BlockSize {
		var res uint8

		for _, monom := range ca.GetPolynomial(i) {
			res ^= monom.Eval(srcBlock)
		}
		dstBlock.SetAt(core.Subscript(i), res)
	}
	dstBlock.Write(dst)
}

func (ca *CA) Save(dst io.Writer) error {
	if err := binary.Write(dst, binary.LittleEndian, ca.Offsets); err != nil {
		return err
	}
	if err := binary.Write(dst, binary.LittleEndian, uint32(len(ca.Arena))); err != nil {
		return err
	}
	return binary.Write(dst, binary.LittleEndian, ca.Arena)
}

func LoadCA(src io.Reader) (*CA, error) {
	ca := &CA{}
	if err := binary.Read(src, binary.LittleEndian, &ca.Offsets); err != nil {
		return nil, err
	}
	var length uint32
	if err := binary.Read(src, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	ca.Arena = make([]math.Monomial, length)
	if err := binary.Read(src, binary.LittleEndian, &ca.Arena); err != nil {
		return nil, err
	}
	return ca, nil
}

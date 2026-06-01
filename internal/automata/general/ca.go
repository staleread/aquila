package general

import (
	"encoding/binary"
	"io"

	"github.com/staleread/aquila/internal/automata/config"
	"github.com/staleread/aquila/internal/automata/math"
	"github.com/staleread/aquila/internal/automata/state"
)

type CA struct {
	Arena   []math.Monomial
	Offsets [state.StateSize - 1]uint32
}

func (ca *CA) GetPolynomial(idx int) []math.Monomial {
	start := uint32(0)
	if idx > 0 {
		start = ca.Offsets[idx-1]
	}
	end := uint32(len(ca.Arena))
	if idx < state.StateSize-1 {
		end = ca.Offsets[idx]
	}
	return ca.Arena[start:end]
}

func (ca *CA) GetMonomialCounts() []int {
	counts := make([]int, state.StateSize)
	for i := range state.StateSize {
		counts[i] = len(ca.GetPolynomial(i))
	}
	return counts
}

func (ca *CA) Apply(dst, src []byte) {
	var srcState state.State
	srcState.Read(src)
	dstState := ca.applyOnState(srcState)
	dstState.Write(dst)
}

func (ca *CA) Save(dst io.Writer) error {
	if err := config.Current.Write(dst); err != nil {
		return err
	}
	if err := binary.Write(dst, binary.LittleEndian, ca.Offsets); err != nil {
		return err
	}
	if err := binary.Write(dst, binary.LittleEndian, uint32(len(ca.Arena))); err != nil {
		return err
	}
	return binary.Write(dst, binary.LittleEndian, ca.Arena)
}

func LoadCA(src io.Reader) (*CA, error) {
	if err := config.Current.Check(src); err != nil {
		return nil, err
	}

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

func (ca *CA) applyOnState(srcState state.State) state.State {
	var dstState state.State
	for i := range state.StateSize {
		var res uint8
		for _, monom := range ca.GetPolynomial(i) {
			res ^= monom.Eval(srcState)
		}
		dstState.SetAt(state.Subscript(i), res)
	}
	return dstState
}

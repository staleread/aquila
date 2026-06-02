//go:build block64

package state

import (
	"encoding/binary"
	"math/bits"
)

const (
	StateSize  = 64
	StateBytes = StateSize / 8
)

type Subscript = uint8
type State uint64

func NewState(subs ...Subscript) State {
	var b State
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

func (s *State) Read(src []byte) {
	_ = src[7] // BCE check
	*s = State(binary.LittleEndian.Uint64(src[0:8]))
}

func (s State) Write(dst []byte) {
	_ = dst[7] // BCE check
	binary.LittleEndian.PutUint64(dst[0:8], uint64(s))
}

func (s State) At(idx Subscript) uint8 {
	return uint8((s >> idx) & 1)
}

func (s *State) SetAt(idx Subscript, bit uint8) {
	*s = (*s &^ (State(1) << idx)) | (State(bit&1) << idx)
}

func (s State) Compare(other State) int {
	if s == other {
		return 0
	}
	if s < other {
		return -1
	}
	return 1
}

func (s State) Or(other State) State {
	return s | other
}

func (s *State) XorWith(other State) {
	*s ^= other
}

func (s State) Contains(other State) bool {
	return s&other == other
}

func (s State) Subscripts(dst []Subscript) []Subscript {
	w := uint64(s)
	for w != 0 {
		tz := bits.TrailingZeros64(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}

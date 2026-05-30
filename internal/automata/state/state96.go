//go:build block96 || (!block16 && !block32 && !block128)

package state

import (
	"encoding/binary"
	"math/bits"
)

const (
	StateSize  = 96
	StateBytes = StateSize / 8
)

type Subscript = uint8
type State [3]uint32

func NewState(subs ...Subscript) State {
	var s State
	for _, i := range subs {
		s.SetAt(i, 1)
	}
	return s
}

func (s *State) Read(src []byte) {
	_ = src[11] // BCE check

	s[0] = binary.LittleEndian.Uint32(src[0:4])
	s[1] = binary.LittleEndian.Uint32(src[4:8])
	s[2] = binary.LittleEndian.Uint32(src[8:12])
}

func (s *State) Write(dst []byte) {
	_ = dst[11] // BCE check

	binary.LittleEndian.PutUint32(dst[0:4], s[0])
	binary.LittleEndian.PutUint32(dst[4:8], s[1])
	binary.LittleEndian.PutUint32(dst[8:12], s[2])
}

func (s State) At(idx Subscript) uint8 {
	return uint8((s[idx/32] >> (idx % 32)) & 1)
}

func (s *State) SetAt(idx Subscript, bit uint8) {
	wordIdx := idx / 32
	shift := idx % 32

	s[wordIdx] = (s[wordIdx] &^ (uint32(1) << shift)) | uint32(bit&1)<<shift
}

func (s State) Compare(other State) int {
	for i := range len(s) {
		if s[i] == other[i] {
			continue
		}
		if s[i] < other[i] {
			return -1
		}
		return 1
	}
	return 0
}

func (s State) Or(other State) State {
	return State{
		s[0] | other[0],
		s[1] | other[1],
		s[2] | other[2],
	}
}

func (s *State) XorWith(other State) {
	s[0] ^= other[0]
	s[1] ^= other[1]
	s[2] ^= other[2]
}

func (s State) Contains(other State) bool {
	return s[0]&other[0] == other[0] && s[1]&other[1] == other[1] && s[2]&other[2] == other[2]
}

func (s State) Subscripts(dst []Subscript) []Subscript {
	for i, w := range s {
		wordOffset := i * 32
		for w != 0 {
			tz := bits.TrailingZeros32(w)
			dst = append(dst, Subscript(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
}

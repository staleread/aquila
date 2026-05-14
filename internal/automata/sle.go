package automata

import (
	"crypto/rand"
	"math/bits"
	"unsafe"
)

type SLE []uint16

func (s SLE) FillRand() {
	byteView := unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), len(s)*2)
	rand.Read(byteView)

	for i := range 16 {
		s[i] |= 1 << i
	}
}

func (s SLE) Solve(b uint16) uint16 {
	return s.substituteBackward(s.substituteForward(b))
}

func (s SLE) Eval(x uint16) uint16 {
	return s.multiplyLower(s.multiplyUpper(x))
}

func (s SLE) Coefs(dst []uint16) {
	for i := range 16 {
		iMask := ^uint16((1 << i) - 1)
		dstRow := s[i] & iMask

		for k := range i {
			if (s[i]>>k)&1 != 1 {
				continue
			}
			kMask := ^uint16((1 << k) - 1)
			dstRow ^= s[k] & kMask
		}
		dst[i] = dstRow
	}
}

func (s SLE) substituteForward(b uint16) uint16 {
	var x uint16

	for i := range 16 {
		dot := uint16(bits.OnesCount16(s[i]&x)) & 1
		x |= (b>>i ^ dot) << i
	}
	return x
}

func (s SLE) substituteBackward(b uint16) uint16 {
	var x uint16

	for i := 15; i >= 0; i-- {
		dot := uint16(bits.OnesCount16(s[i]&x)) & 1
		x |= (b>>i ^ dot) << i
	}
	return x
}

func (s SLE) multiplyLower(a uint16) uint16 {
	var b uint16

	for i := range 16 {
		mask := uint16((1 << (i + 1)) - 1)
		b |= (uint16(bits.OnesCount16(s[i]&a&mask)) & 1) << i
	}
	return b
}

func (s SLE) multiplyUpper(a uint16) uint16 {
	var b uint16

	for i := range 16 {
		mask := ^uint16((1 << i) - 1)
		b |= (uint16(bits.OnesCount16(s[i]&a&mask)) & 1) << i
	}
	return b
}

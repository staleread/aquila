package invertible

import (
	"crypto/rand"
	"unsafe"
)

type SLE []Vector

func (s SLE) FillRand() {
	byteView := unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), len(s)*2)
	rand.Read(byteView)

	for i := range Dim {
		s[i] |= 1 << i
	}
}

func (s SLE) Solve(b Vector) Vector {
	return s.substituteBackward(s.substituteForward(b))
}

func (s SLE) Eval(x Vector) Vector {
	return s.multiplyLower(s.multiplyUpper(x))
}

func (s SLE) Coefs(dst []Vector) {
	for i := range Dim {
		iMask := ^Vector((1 << i) - 1)
		dstRow := s[i] & iMask

		for k := range i {
			if (s[i]>>k)&1 != 1 {
				continue
			}
			kMask := ^Vector((1 << k) - 1)
			dstRow ^= s[k] & kMask
		}
		dst[i] = dstRow
	}
}

func (s SLE) substituteForward(v Vector) Vector {
	var res Vector

	for i := range Dim {
		res |= (v>>i ^ s[i].Dot(res)) << i
	}
	return res
}

func (s SLE) substituteBackward(v Vector) Vector {
	var res Vector

	for i := Dim - 1; i >= 0; i-- {
		res |= (v>>i ^ s[i].Dot(res)) << i
	}
	return res
}

func (s SLE) multiplyLower(v Vector) Vector {
	var res Vector

	for i := range Dim {
		mask := Vector((1 << (i + 1)) - 1)
		res |= s[i].Dot(v&mask) << i
	}
	return res
}

func (s SLE) multiplyUpper(v Vector) Vector {
	var res Vector

	for i := range Dim {
		mask := ^Vector((1 << i) - 1)
		res |= s[i].Dot(v&mask) << i
	}
	return res
}

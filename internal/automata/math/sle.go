package math

import (
	"io"
	"unsafe"
)

type SLE struct {
	data []Vector
}

func NewSLE(arena []byte) SLE {
	data := unsafe.Slice((*Vector)(unsafe.Pointer(&arena[0])), len(arena)/2)

	return SLE{data}
}

func (s *SLE) Generate(rnd io.Reader) error {
	byteView := unsafe.Slice((*byte)(unsafe.Pointer(&s.data[0])), len(s.data)*2)
	if _, err := io.ReadFull(rnd, byteView); err != nil {
		return err
	}

	for i := range VectorSize {
		s.data[i] |= 1 << i
	}
	return nil
}

func (s *SLE) Solve(b Vector) Vector {
	return s.substituteBackward(s.substituteForward(b))
}

func (s *SLE) Eval(x Vector) Vector {
	return s.multiplyLower(s.multiplyUpper(x))
}

func (s *SLE) Coefs(dst []Vector) {
	for i := range VectorSize {
		iMask := ^Vector((1 << i) - 1)
		dstRow := s.data[i] & iMask

		for k := range i {
			if (s.data[i]>>k)&1 != 1 {
				continue
			}
			kMask := ^Vector((1 << k) - 1)
			dstRow ^= s.data[k] & kMask
		}
		dst[i] = dstRow
	}
}

func (s *SLE) substituteForward(v Vector) Vector {
	var res Vector

	for i, row := range s.data {
		res |= ((v >> i) ^ row.Dot(res)) & 1 << i
	}
	return res
}

func (s *SLE) substituteBackward(v Vector) Vector {
	var res Vector

	for i := VectorSize - 1; i >= 0; i-- {
		row := s.data[i]
		res |= ((v >> i) ^ row.Dot(res)) & 1 << i
	}
	return res
}

func (s *SLE) multiplyLower(v Vector) Vector {
	var res Vector

	for i, row := range s.data {
		mask := Vector((1 << (i + 1)) - 1)
		res |= row.Dot(v&mask) << i
	}
	return res
}

func (s *SLE) multiplyUpper(v Vector) Vector {
	var res Vector

	for i, row := range s.data {
		mask := ^Vector((1 << i) - 1)
		res |= row.Dot(v&mask) << i
	}
	return res
}

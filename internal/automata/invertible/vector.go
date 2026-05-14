package invertible

import "math/bits"

type Vector uint16

const Ones Vector = 0xFFFF

func (v Vector) Sum() Vector {
	return Vector(bits.OnesCount16(uint16(v))) & 1
}

func (v Vector) Dot(other Vector) Vector {
	return (v & other).Sum()
}

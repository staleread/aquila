package math

import "math/bits"

type Vector uint16

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount16(uint16(v & other))
	return Vector(ones & 1)
}

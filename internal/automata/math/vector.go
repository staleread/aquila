package math

import "math/bits"

type Vector uint8

const (
	VectorSize  = 8
	VectorBytes = VectorSize / 8
)

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount8(uint8(v & other))
	return Vector(ones & 1)
}

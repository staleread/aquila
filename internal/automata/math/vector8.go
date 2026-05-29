//go:build sle8 || !sle16

package math

import "math/bits"

const (
	VectorSize  = 8
	VectorBytes = VectorSize / 8
)

type Vector uint8

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount8(uint8(v & other))
	return Vector(ones & 1)
}

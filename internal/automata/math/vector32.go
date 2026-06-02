//go:build fold32

package math

import "math/bits"

const (
	VectorSize  = 32
	VectorBytes = VectorSize / 8
)

type Vector uint32

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount32(uint32(v & other))
	return Vector(ones & 1)
}

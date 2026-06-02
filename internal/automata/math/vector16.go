//go:build fold16 || (!fold8 && !fold32)

package math

import "math/bits"

const (
	VectorSize  = 16
	VectorBytes = VectorSize / 8
)

type Vector uint16

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount16(uint16(v & other))
	return Vector(ones & 1)
}

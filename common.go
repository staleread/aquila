package main

import "math"

const (
	VarIdxBytes = 2 // Size of VarIdx type
	MaxVarIdx   = math.MaxUint16
)

// Unsigned integer used to represent an bit size of a batch vector
type Size = uint16
type VarIdx = Size

// Batch of bits
type Batch uint32

func BatchMask(bBytes Size) Batch {
	return math.MaxUint32 >> (32 - bBytes*8)
}

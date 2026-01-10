package main

import "math"

const (
	VarIdxBytes = 2
	MaxVarIdx   = math.MaxUint16
)

// Unsigned integer used to represent an bit size of a batch vector
type Size = uint16
type VarIdx = Size

type Elem = uint8

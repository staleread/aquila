package main

const (
	ElemCap      = 16 // Size of Elem type
	ElemLen      = 8  // Active bits of element. Must be less than ElemCap
	ElemMask     = (1 << ElemLen) - 1
	ElemLenBytes = (ElemLen + 7) / 8
)

type Elem = uint16
type Dim = int32

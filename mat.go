package main

import (
	"fmt"
	"strings"
)

type Dim int32
type BitMat struct {
	Dim  Dim
	Data []BitPack
}

func (mat *BitMat) String() string {
	n := mat.Dim
	sb := strings.Builder{}

	for i := range n {
		for j := range n {
			val := mat.Data[n*i+j] & 1
			sb.WriteString(fmt.Sprintf("%d ", val))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func Zeros(n Dim) *BitMat {
	return &BitMat{
		Dim:  n,
		Data: make([]BitPack, n*n),
	}
}

package main

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type Vec []Elem

func Zeros(n Dim) Vec {
	return Vec(make([]Elem, n))
}

func Rands(n Dim) Vec {
	buffSize := n * ElemLenBytes
	buff := make([]byte, buffSize)
	rand.Read(buff)

	vec := Zeros(n)

	for buffIdx := range buffSize {
		vecIdx := buffIdx / ElemLenBytes

		shift := (buffIdx % ElemLenBytes) * 8
		vec[vecIdx] |= Elem(buff[buffIdx]) << shift

		if shift == 0 {
			vec[vecIdx] &= ElemMask
		}
	}
	return vec
}

func (vec Vec) String() string {
	n := len(vec)
	sb := strings.Builder{}

	for i := range n {
		val := vec[i]
		sb.WriteString(fmt.Sprintf("%0*b\n", ElemLen, val))
	}
	return sb.String()
}

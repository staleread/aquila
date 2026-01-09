package main

import (
	"fmt"
	"strings"
)

type Vec []Elem

func (vec Vec) String() string {
	n := len(vec)
	sb := strings.Builder{}

	for i := range n {
		val := vec[i]
		sb.WriteString(fmt.Sprintf("%0*b\n", ElemLen, val))
	}
	return sb.String()
}

package main

import (
	"fmt"
	"strings"
)

type BitVec []BitPack

func (vec BitVec) String() string {
	n := len(vec)
	sb := strings.Builder{}

	for i := range n {
		val := vec[i] & 1
		sb.WriteString(fmt.Sprintf("%d\n", val))
	}
	return sb.String()
}

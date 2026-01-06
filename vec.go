// Utilities for matrices over GF(2)
package main

import (
	"strconv"
	"strings"
)

type Vec struct {
	Size int
	Data []uint8
}

func (vec *Vec) String() string {
	n := vec.Size
	sb := strings.Builder{}

	for i := range n {
		val := int(vec.Data[i])

		sb.WriteString(strconv.Itoa(val))
		sb.WriteRune('\n')
	}
	return sb.String()
}

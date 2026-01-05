// Utilities for matrices over GF(2)
package main

import (
	"strconv"
	"strings"
)

type Mat struct {
	Size int
	Data []uint8
}

func IdMat(n int) Mat {
	data := make([]uint8, n*n)

	// set the diagonal elements to 1's
	for i := range n {
		data[n*i+i] = 1
	}

	return Mat{Size: n, Data: data}
}

func (mat *Mat) String() string {
	n := mat.Size
	builder := strings.Builder{}

	for row := range n {
		for col := range n {
			val := int(mat.Data[n*row+col])

			builder.WriteString(strconv.Itoa(val))
			builder.WriteRune(' ')
		}
		builder.WriteRune('\n')
	}

	return builder.String()
}

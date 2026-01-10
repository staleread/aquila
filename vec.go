package main

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type BVec struct {
	bBytes Size
	data   []Batch
}

func ZeroBVec(n, bBytes Size) *BVec {
	return &BVec{bBytes, make([]Batch, n)}
}

func RandBVec(n, bBytes Size) *BVec {
	buffSize := n * bBytes
	buff := make([]byte, buffSize)
	rand.Read(buff)

	vec := ZeroBVec(n, bBytes)
	bMask := BatchMask(bBytes)

	for buffIdx := range buffSize {
		vecIdx := buffIdx / bBytes

		shift := (buffIdx % bBytes) * 8
		vec.data[vecIdx] |= Batch(buff[buffIdx]) << shift

		if shift == 0 {
			vec.data[vecIdx] &= bMask
		}
	}
	return vec
}

func (vec *BVec) String() string {
	sb := strings.Builder{}

	for i := range len(vec.data) {
		val := vec.data[i]
		sb.WriteString(fmt.Sprintf("%0*b\n", vec.bBytes*8, val))
	}
	return sb.String()
}

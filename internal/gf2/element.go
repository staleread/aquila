package gf2

import "crypto/rand"

type Element = uint8

const (
	Zero Element = 0
	One  Element = 1
)

func ElementsInBytes(bytes int) int {
	return bytes * 8
}

func FillRand(els []Element) {
	n := len(els)

	// HACK to avoid extra allocations.
	bufStart := n - 1 - (n+7)/8
	rand.Read(els[bufStart:])

	for i := range n {
		els[i] = (els[bufStart+i/8] >> (i % 8)) & 1
	}
}

func ReadElements(dst []Element, src []byte) {
	if len(dst) > len(src)*8 {
		panic("Failed to read elements. Not enough bytes to fill set all elements")
	}

	for i := range len(dst) {
		val := (src[i/8] >> (i % 8)) & 1
		dst[i] = Element(val)
	}
}

func WriteElements(dst []byte, src []Element) {
	if len(dst)*8 < len(src) {
		panic("Failed to write elements. Not enough bytes to write all elements")
	}

	for i := range len(dst) {
		var val byte

		for j := range len(src) - i*8 {
			val |= byte(src[i*8+j]) << j
		}
		dst[i] = val
	}
}

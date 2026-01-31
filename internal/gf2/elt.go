package gf2

import "crypto/rand"

type Elt uint8

func ElsInBytes(bytes int) int {
	return bytes * 8
}

func RandEls(n int) []Elt {
	els := make([]Elt, n)

	buf := make([]byte, (n+7)/8)
	rand.Read(buf)

	for i := range n {
		val := buf[i/8] >> (i % 8)
		els[i] = Elt(val & 1)
	}
	return els
}

func RandNzEls(n int) []Elt {
	els := make([]Elt, n)

	for i := range n {
		els[i] = 1
	}
	return els
}

func Add(a, b Elt) Elt {
	return a ^ b
}

func Sub(a, b Elt) Elt {
	return a ^ b
}

func Mul(a, b Elt) Elt {
	return a & b
}

func Div(a, b Elt) Elt {
	if b == 0 {
		panic("Division by zero")
	}
	return a
}

func ReadEls(dst []Elt, src []byte) {
	if len(dst) > len(src)*8 {
		panic("Failed to read elements. Not enough bytes to fill set all elements")
	}

	for i := range len(dst) {
		val := (src[i/8] >> (i % 8)) & 1
		dst[i] = Elt(val)
	}
}

func WriteEls(dst []byte, src []Elt) {
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

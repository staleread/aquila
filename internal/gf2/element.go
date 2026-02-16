package gf2

import "crypto/rand"

type Element uint8

func ElementsInBytes(bytes int) int {
	return bytes * 8
}

func RandElements(n int) []Element {
	els := make([]Element, n)

	buf := make([]byte, (n+7)/8)
	rand.Read(buf)

	for i := range n {
		val := buf[i/8] >> (i % 8)
		els[i] = Element(val & 1)
	}
	return els
}

func RandNonZeroElements(n int) []Element {
	els := make([]Element, n)

	for i := range n {
		els[i] = 1
	}
	return els
}

func Add(a, b Element) Element {
	return a ^ b
}

func Sub(a, b Element) Element {
	return a ^ b
}

func Mul(a, b Element) Element {
	return a & b
}

func Div(a, b Element) Element {
	if b == 0 {
		panic("Division by zero")
	}
	return a
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

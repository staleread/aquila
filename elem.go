package main

const byteBits = 8

type elem uint8

func randElem() elem {
	return elem(randByte() & 1)
}

func randNzElem() elem {
	return elem(1)
}

func add(a, b elem) elem {
	return a ^ b
}

func sub(a, b elem) elem {
	return a ^ b
}

func mul(a, b elem) elem {
	return a & b
}

func div(a, b elem) elem {
	if b == 0 {
		panic("Division by zero")
	}
	return a
}

func bytesToElemCnt(bytes int) {
	return bytes * 8
}

func readElems(src []byte) []elem {
	els := make([]elem, len(src)*byteBits)

	for i := range len(src) {
		val := src[i]

		for j := range byteBits {
			els[i*byteBits+j] = (elem(val) >> (byteBits - j - 1)) & 1
		}
	}
	return els
}

func writeElems(dst []byte, els []elem) {
	bCnt := len(els) / byteBits

	for i := range bCnt {
		var b byte = 0

		for j := range byteBits {
			b <<= 1
			b |= byte(els[byteBits*i+j])
		}
		dst[i] = b
	}
}

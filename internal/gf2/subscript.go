package gf2

import "crypto/rand"

type Subscript uint8

func RandSubscripts(n int) []Subscript {
	subscripts := make([]Subscript, n)

	buf := make([]byte, n)
	rand.Read(buf)

	for i := range n {
		subscripts[i] = Subscript(buf[i])
	}
	return subscripts
}

package gf2

import "crypto/rand"

type Subscript uint8

func RandSubscripts(n int) []Subscript {
	subs := make([]Subscript, n)

	buf := make([]byte, n)
	rand.Read(buf)

	for i := range n {
		subs[i] = Subscript(buf[i])
	}
	return subs
}

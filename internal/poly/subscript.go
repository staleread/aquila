package poly

import "crypto/rand"

type Subscript uint16

func RandSubscripts(n int) []Subscript {
	subs := make([]Subscript, n)

	// TODO move to a util for generating rand bytes
	buf := make([]byte, n*2)
	rand.Read(buf)

	for i := range n {
		subs[i] = Subscript(buf[i])<<8 | Subscript(buf[n+i])
	}
	return subs
}

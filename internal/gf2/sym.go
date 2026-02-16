package gf2

import "crypto/rand"

type Sym uint8

func RandSyms(n int) []Sym {
	syms := make([]Sym, n)

	buf := make([]byte, n)
	rand.Read(buf)

	for i := range n {
		syms[i] = Sym(buf[i])
	}
	return syms
}

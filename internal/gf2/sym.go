package gf2

import "crypto/rand"

type Sym uint16

func RandSyms(n int) []Sym {
	syms := make([]Sym, n)

	buf := make([]byte, n*2)
	rand.Read(buf)

	for i := range n {
		syms[i] = (Sym(buf[i]) << 8) | Sym(buf[n+i])
	}
	return syms
}

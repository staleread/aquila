package main

import crand "crypto/rand"

var gRand = newRandSrc(1024)

func randByte() byte {
	if gRand.pos >= len(gRand.buf) {
		gRand.refill()
	}
	out := gRand.buf[gRand.pos]
	gRand.pos++

	return out
}

type randSrc struct {
	buf []byte
	pos int
}

func newRandSrc(n int) *randSrc {
	buf := make([]byte, n)
	return &randSrc{buf, n}
}

func (r *randSrc) refill() {
	crand.Read(r.buf[:r.pos])
	r.pos = 0
}

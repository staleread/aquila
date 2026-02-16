package main

import (
	"github.com/staleread/aquila/internal/ca"
	f "github.com/staleread/aquila/internal/gf2"
)

type PublicKey struct {
	bSize int
	ca    *ca.GeneralCA
}

func (k *PublicKey) Encrypt(dst, src []byte) {
	if len(src)%k.bSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]f.Element, f.ElementsInBytes(k.bSize))

	for i := range len(src) / k.bSize {
		from := k.bSize * i
		to := k.bSize * (i + 1)

		f.ReadElements(tmp, src[from:to])

		k.ca.Apply(tmp)

		f.WriteElements(dst[from:to], tmp)
	}
}

package main

import (
	"github.com/staleread/aquila/internal/ca/general"
	"github.com/staleread/aquila/internal/field"
)

type PublicKey struct {
	bSize int
	ca    *general.CA
}

func (k *PublicKey) Encrypt(dst, src []byte) {
	if len(src)%k.bSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]field.Element, field.ElementsInBytes(k.bSize))

	for i := range len(src) / k.bSize {
		from := k.bSize * i
		to := k.bSize * (i + 1)

		field.ReadElements(tmp, src[from:to])

		k.ca.Apply(tmp)

		field.WriteElements(dst[from:to], tmp)
	}
}

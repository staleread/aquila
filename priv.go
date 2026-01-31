package main

import (
	"errors"
	"github.com/staleread/aquila/internal/ca"
	f "github.com/staleread/aquila/internal/gf2"
)

type PrivateKey struct {
	bSize int
	ca    *ca.InvertibleCA
}

type config struct {
	bSize   int
	folds   int
	polyDeg int
	rules   int
}

var configs = map[int]config{
	16: {
		bSize:   16,
		folds:   8,
		polyDeg: 2,
		rules:   16,
	},
	24: {
		bSize:   24,
		folds:   12,
		polyDeg: 3,
		rules:   24,
	},
	32: {
		bSize:   32,
		folds:   16,
		polyDeg: 3,
		rules:   32,
	},
}

func GenerateKey(bSize int) (*PrivateKey, error) {
	cfg, ok := configs[bSize]

	if !ok {
		return nil, errors.New("Unsupported block size")
	}

	caSize := f.ElsInBytes(bSize)
	ca := ca.NewInvertibleCA(caSize, cfg.folds, cfg.polyDeg, cfg.rules)

	return &PrivateKey{bSize, ca}, nil
}

func (k *PrivateKey) Decrypt(dst, src []byte) {
	if len(src)%k.bSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]f.Elt, f.ElsInBytes(k.bSize))

	for i := range len(src) / k.bSize {
		from := k.bSize * i
		to := k.bSize * (i + 1)

		f.ReadEls(tmp, src[from:to])

		k.ca.ApplyInverse(tmp)

		f.WriteEls(dst[from:to], tmp)
	}
}

func (k *PrivateKey) Public() *PublicKey {
	ca := k.ca.ToGeneralCA()
	return &PublicKey{k.bSize, ca}
}

func (k *PrivateKey) encryptTest(dst, src []byte) {
	if len(src)%k.bSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]f.Elt, f.ElsInBytes(k.bSize))

	for i := range len(src) / k.bSize {
		from := k.bSize * i
		to := k.bSize * (i + 1)

		f.ReadEls(tmp, src[from:to])

		k.ca.Apply(tmp)

		f.WriteEls(dst[from:to], tmp)
	}
}

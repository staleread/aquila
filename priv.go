package main

import (
	"errors"
	"github.com/staleread/aquila/internal/ca"
	f "github.com/staleread/aquila/internal/gf2"
)

type PrivateKey struct {
	blockSize int
	ca    *ca.InvertibleCA
}

type config struct {
	blockSize   int
	folds   int
	degree int
	rules   int
}

var configs = map[int]config{
	8: {
		blockSize:   8,
		folds:   8,
		degree: 2,
		rules:   8,
	},
	16: {
		blockSize:   16,
		folds:   8,
		degree: 3,
		rules:   16,
	},
	24: {
		blockSize:   24,
		folds:   12,
		degree: 3,
		rules:   24,
	},
	32: {
		blockSize:   32,
		folds:   16,
		degree: 3,
		rules:   32,
	},
}

func GenerateKey(blockSize int) (*PrivateKey, error) {
	cfg, ok := configs[blockSize]

	if !ok {
		return nil, errors.New("Unsupported block size")
	}

	caSize := f.ElementsInBytes(blockSize)
	ca := ca.NewInvertibleCA(caSize, cfg.folds, cfg.degree, cfg.rules)

	return &PrivateKey{blockSize, ca}, nil
}

func (k *PrivateKey) Decrypt(dst, src []byte) {
	if len(src)%k.blockSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]f.Element, f.ElementsInBytes(k.blockSize))

	for i := range len(src) / k.blockSize {
		from := k.blockSize * i
		to := k.blockSize * (i + 1)

		f.ReadElements(tmp, src[from:to])

		k.ca.ApplyInverse(tmp)

		f.WriteElements(dst[from:to], tmp)
	}
}

func (k *PrivateKey) Public() *PublicKey {
	ca := k.ca.ToGeneralCA()
	return &PublicKey{k.blockSize, ca}
}

func (k *PrivateKey) encryptTest(dst, src []byte) {
	if len(src)%k.blockSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]f.Element, f.ElementsInBytes(k.blockSize))

	for i := range len(src) / k.blockSize {
		from := k.blockSize * i
		to := k.blockSize * (i + 1)

		f.ReadElements(tmp, src[from:to])

		k.ca.Apply(tmp)

		f.WriteElements(dst[from:to], tmp)
	}
}

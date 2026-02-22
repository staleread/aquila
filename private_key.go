package main

import (
	"errors"
	"github.com/staleread/aquila/internal/automata/invertible"
	"github.com/staleread/aquila/internal/field"
)

type PrivateKey struct {
	blockSize int
	ca        *invertible.CA
}

type config struct {
	blockSize int
	folds     int
	degree    int
	rules     int
}

var configs = map[int]config{
	8: {
		blockSize: 8,
		folds:     8,
		degree:    2,
		rules:     8,
	},
	16: {
		blockSize: 16,
		folds:     8,
		degree:    3,
		rules:     16,
	},
	24: {
		blockSize: 24,
		folds:     12,
		degree:    3,
		rules:     24,
	},
	32: {
		blockSize: 32,
		folds:     16,
		degree:    3,
		rules:     32,
	},
}

func GenerateKey(blockSize int) (*PrivateKey, error) {
	cfg, ok := configs[blockSize]

	if !ok {
		return nil, errors.New("Unsupported block size")
	}

	caSize := field.ElementsInBytes(blockSize)
	ca := invertible.NewCA(caSize, cfg.folds, cfg.degree, cfg.rules)

	return &PrivateKey{blockSize, ca}, nil
}

func (k *PrivateKey) Decrypt(dst, src []byte) {
	if len(src)%k.blockSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]field.Element, field.ElementsInBytes(k.blockSize))

	for i := range len(src) / k.blockSize {
		from := k.blockSize * i
		to := k.blockSize * (i + 1)

		field.ReadElements(tmp, src[from:to])

		k.ca.ApplyInverse(tmp)

		field.WriteElements(dst[from:to], tmp)
	}
}

func (k *PrivateKey) Public() *PublicKey {
	ca := k.ca.General()
	return &PublicKey{k.blockSize, ca}
}

func (k *PrivateKey) encryptTest(dst, src []byte) {
	if len(src)%k.blockSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	tmp := make([]field.Element, field.ElementsInBytes(k.blockSize))

	for i := range len(src) / k.blockSize {
		from := k.blockSize * i
		to := k.blockSize * (i + 1)

		field.ReadElements(tmp, src[from:to])

		k.ca.Apply(tmp)

		field.WriteElements(dst[from:to], tmp)
	}
}

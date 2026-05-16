package main

import (
	"crypto/rand"
	"errors"
	"io"

	"github.com/staleread/aquila/internal/automata"
	"github.com/staleread/aquila/internal/automata/invertible"
)

type PrivateKey struct {
	ca *invertible.CA
}

func GenerateKey() (*PrivateKey, error) {
	k := &PrivateKey{
		ca: invertible.NewCA(),
	}

	if err := k.ca.Generate(rand.Reader); err != nil {
		return nil, err
	}
	return k, nil
}

func LoadPrivateKey(src io.Reader) (*PrivateKey, error) {
	k := &PrivateKey{
		ca: invertible.NewCA(),
	}

	if err := k.ca.Load(src); err != nil {
		return nil, err
	}
	return k, nil
}

// Decrypt decrypts the full contents of src into dst.
// Both slices must be sized to a multiple of your block size.
// Returns an error if the lengths are mismatched or invalid.
func (k *PrivateKey) Decrypt(dst, src []byte) error {
	if len(src)%automata.BlockBytes != 0 || len(dst) < len(src) {
		return errors.New("invalid buffer size")
	}

	for i := 0; i < len(src); i += automata.BlockBytes {
		srcBlock := automata.LoadBlock(src[i : i+automata.BlockBytes])
		k.ca.Revert(srcBlock)
		srcBlock.WriteTo(dst[i : i+automata.BlockBytes])
	}
	return nil
}

func (k *PrivateKey) encrypt(dst, src []byte) error {
	if len(src)%automata.BlockBytes != 0 || len(dst) < len(src) {
		return errors.New("invalid buffer size")
	}

	for i := 0; i < len(src); i += automata.BlockBytes {
		srcBlock := automata.LoadBlock(src[i : i+automata.BlockBytes])
		k.ca.Apply(srcBlock)
		srcBlock.WriteTo(dst[i : i+automata.BlockBytes])
	}
	return nil
}

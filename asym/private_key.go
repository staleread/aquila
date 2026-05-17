package asym

import (
	"errors"
	"io"

	"github.com/staleread/aquila/internal/automata"
	"github.com/staleread/aquila/internal/automata/invertible"
)

type PrivateKey struct {
	dec BlockDecrypter
}

func GenerateKey(rnd io.Reader) (*PrivateKey, error) {
	ca := invertible.NewCA()

	if err := ca.Generate(rnd); err != nil {
		return nil, err
	}
	return &PrivateKey{ca}, nil
}

func LoadPrivateKey(src io.Reader) (*PrivateKey, error) {
	ca := invertible.NewCA()

	if err := ca.Load(src); err != nil {
		return nil, err
	}
	return &PrivateKey{ca}, nil
}

// Decrypts the full contents of src into dst.
// Both slices must be sized to a multiple of BlockBytes.
// Returns an error if the lengths are mismatched or invalid.
func (k *PrivateKey) Decrypt(dst, src []byte) error {
	if len(src)%automata.BlockBytes != 0 || len(dst) < len(src) {
		return errors.New("invalid buffer size")
	}

	for i := 0; i < len(src); i += automata.BlockBytes {
		k.dec.Revert(dst[i:i+automata.BlockBytes], src[i:i+automata.BlockBytes])
	}
	return nil
}

// Encrypts the full contents of src into dst.
// Both slices must be sized to a multiple of BlockBytes.
// Returns an error if the lengths are mismatched or invalid.
func (k *PrivateKey) Encrypt(dst, src []byte) error {
	if len(src)%automata.BlockBytes != 0 || len(dst) < len(src) {
		return errors.New("invalid buffer size")
	}

	for i := 0; i < len(src); i += automata.BlockBytes {
		k.dec.Apply(dst[i:i+automata.BlockBytes], src[i:i+automata.BlockBytes])
	}
	return nil
}

func (k *PrivateKey) Public() *PublicKey {
	return nil
}

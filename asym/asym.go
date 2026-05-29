package asym

import (
	"crypto"
	"errors"
	"io"

	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/invertible"
	"github.com/staleread/aquila/internal/automata/state"
)

const BlockSize = state.StateBytes

var _ crypto.Decrypter = (*PrivateKey)(nil)

type PrivateKey struct {
	ca  *invertible.CA
	pub *PublicKey
}

func (k *PrivateKey) Public() crypto.PublicKey {
	if k.pub == nil {
		gen, err := k.ca.DeriveGeneralCA()
		if err != nil {
			return nil
		}
		k.pub = &PublicKey{gen}
	}
	return k.pub
}

func (k *PrivateKey) Encode(dst io.Writer) error {
	return k.ca.Save(dst)
}

func DecodePrivateKey(src io.Reader) (*PrivateKey, error) {
	ca := invertible.NewCA()
	if err := ca.Load(src); err != nil {
		return nil, err
	}
	return &PrivateKey{ca: ca}, nil
}

func (k *PrivateKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	if len(msg)%state.StateBytes != 0 {
		return nil, errors.New("invalid ciphertext length")
	}

	plaintext = make([]byte, len(msg))
	for i := 0; i < len(msg); i += state.StateBytes {
		k.ca.Revert(plaintext[i:i+state.StateBytes], msg[i:i+state.StateBytes])
	}

	return plaintext, nil
}

type PublicKey struct {
	ca *general.CA
}

func (k *PublicKey) Encode(dst io.Writer) error {
	return k.ca.Save(dst)
}

func DecodePublicKey(src io.Reader) (*PublicKey, error) {
	ca, err := general.LoadCA(src)
	if err != nil {
		return nil, err
	}
	return &PublicKey{ca: ca}, nil
}

func (k *PublicKey) Encrypt(rand io.Reader, msg []byte) (ciphertext []byte, err error) {
	if len(msg)%state.StateBytes != 0 {
		return nil, errors.New("invalid plaintext length")
	}

	ciphertext = make([]byte, len(msg))
	for i := 0; i < len(msg); i += state.StateBytes {
		k.ca.Apply(ciphertext[i:i+state.StateBytes], msg[i:i+state.StateBytes])
	}

	return ciphertext, nil
}

func GenerateKeyPair(rnd io.Reader) (*PrivateKey, *PublicKey, error) {
	invertibleCa := invertible.NewCA()

	if err := invertibleCa.Generate(rnd); err != nil {
		return nil, nil, err
	}

	generalCa, err := invertibleCa.DeriveGeneralCA()

	if err != nil {
		return nil, nil, err
	}

	pub := &PublicKey{generalCa}
	return &PrivateKey{invertibleCa, pub}, pub, nil
}

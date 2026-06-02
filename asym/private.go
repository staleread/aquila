package asym

import (
	"crypto"
	"errors"
	"io"

	"github.com/staleread/aquila/internal/automata/invertible"
	"github.com/staleread/aquila/internal/automata/state"
)

var _ crypto.Decrypter = (*PrivateKey)(nil)
var _ crypto.Signer = (*PrivateKey)(nil)

type PrivateKey struct {
	ca  *invertible.CA
	pub *PublicKey
}

func GeneratePrivateKey(rnd io.Reader) (*PrivateKey, error) {
	ca := invertible.NewCA()
	if err := ca.Generate(rnd); err != nil {
		return nil, err
	}
	return &PrivateKey{ca: ca}, nil
}

func (k *PrivateKey) PublicKey() (*PublicKey, error) {
	if k.pub == nil {
		gen, err := k.ca.DeriveGeneralCA()
		if err != nil {
			return nil, err
		}
		k.pub = &PublicKey{gen}
	}
	return k.pub, nil
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

func (k *PrivateKey) Decode(src io.Reader) error {
	ca := invertible.NewCA()
	if err := ca.Load(src); err != nil {
		return err
	}
	k.ca = ca
	return nil
}

func (k *PrivateKey) Encode(dst io.Writer) error {
	return k.ca.Save(dst)
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

func (k *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if len(digest)%state.StateBytes != 0 {
		return nil, errors.New("invalid digest length")
	}

	signature = make([]byte, len(digest))
	for i := 0; i < len(digest); i += state.StateBytes {
		k.ca.Revert(signature[i:i+state.StateBytes], digest[i:i+state.StateBytes])
	}

	return signature, nil
}

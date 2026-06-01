package asym

import (
	"bytes"
	"errors"
	"io"

	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/state"
)

type PublicKey struct {
	ca *general.CA
}

type PublicKeyInfo struct {
	BlockSizeBytes int
	BlockSizeBits  int
	FoldSize       int
	Degree         int
	Compositions   int
	MonomialCounts []int
}

func (k *PublicKey) Describe() PublicKeyInfo {
	return PublicKeyInfo{
		BlockSizeBytes: BlockSize,
		BlockSizeBits:  BlockSize * 8,
		FoldSize:       FoldSize,
		Degree:         Degree,
		Compositions:   Compositions,
		MonomialCounts: k.ca.GetMonomialCounts(),
	}
}

func (k *PublicKey) Decode(src io.Reader) error {
	ca, err := general.LoadCA(src)
	if err != nil {
		return err
	}
	k.ca = ca
	return nil
}

func (k *PublicKey) Encode(dst io.Writer) error {
	return k.ca.Save(dst)
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

func (k *PublicKey) Verify(digest []byte, signature []byte) bool {
	if len(signature) != len(digest) || len(signature)%state.StateBytes != 0 {
		return false
	}

	temp := make([]byte, len(signature))
	for i := 0; i < len(signature); i += state.StateBytes {
		k.ca.Apply(temp[i:i+state.StateBytes], signature[i:i+state.StateBytes])
	}

	return bytes.Equal(temp, digest)
}

func (k *PublicKey) ExportToANF(w io.Writer, input []byte) error {
	if len(input) != BlockSize {
		return errors.New("invalid input block length")
	}
	var s state.State
	s.Read(input)
	return k.ca.ExportToANF(w, s)
}

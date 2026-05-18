package sym

import (
	"crypto/cipher"
	"io"

	"github.com/staleread/aquila/internal/automata/core"
	"github.com/staleread/aquila/internal/automata/invertible"
)

var _ cipher.Block = (*AquilaBlock)(nil)

type AquilaBlock struct {
	ca *invertible.CA
}

func New(rnd io.Reader) (*AquilaBlock, error) {
	ca := invertible.NewCA()

	if err := ca.Generate(rnd); err != nil {
		return nil, err
	}

	return &AquilaBlock{ca}, nil
}

func Decode(src io.Reader) (*AquilaBlock, error) {
	ca := invertible.NewCA()

	if err := ca.Load(src); err != nil {
		return nil, err
	}
	return &AquilaBlock{ca}, nil
}

func (b *AquilaBlock) BlockSize() int          { return core.BlockSize }
func (b *AquilaBlock) Encrypt(dst, src []byte) { b.ca.Apply(dst, src) }
func (b *AquilaBlock) Decrypt(dst, src []byte) { b.ca.Revert(dst, src) }

func (b *AquilaBlock) Encode(dst io.Writer) error { return b.ca.Save(dst) }

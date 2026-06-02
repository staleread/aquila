package asym

import (
	"io"

	"github.com/staleread/aquila/internal/automata/config"
	"github.com/staleread/aquila/internal/automata/math"
	"github.com/staleread/aquila/internal/automata/state"
)

const (
	BlockSize    = state.StateBytes
	FoldSize     = math.VectorSize
	Degree       = math.ConfusionDegree
	Compositions = config.CompositionCount
)

func GenerateKeyPair(rnd io.Reader) (*PrivateKey, *PublicKey, error) {
	priv, err := GeneratePrivateKey(rnd)
	if err != nil {
		return nil, nil, err
	}

	pub, err := priv.PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return priv, pub, nil
}

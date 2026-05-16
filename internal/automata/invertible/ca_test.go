package invertible_test

import (
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata"
	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestCAInvertibility(t *testing.T) {
	ca := invertible.NewCA()
	rng := rand.New(rand.NewSource(42))

	if err := ca.Generate(rng); err != nil {
		t.Fatalf("Failed to generate CA: %v", err)
	}

	for range 10 {
		block := &automata.Block{
			rng.Uint32(),
			rng.Uint32(),
			rng.Uint32(),
		}

		original := *block

		ca.Apply(block)
		if *block == original {
			t.Errorf("Apply did not change the block (highly unlikely)")
		}

		ca.Revert(block)
		if *block != original {
			t.Errorf("Revert failed: original=%v, reverted=%v", original, block)
		}
	}
}

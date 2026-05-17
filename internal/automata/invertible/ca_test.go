package invertible_test

import (
	"bytes"
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

	src := make([]byte, automata.BlockBytes)
	intermediate := make([]byte, automata.BlockBytes)
	reverted := make([]byte, automata.BlockBytes)

	for range 10 {
		rng.Read(src)

		original := make([]byte, automata.BlockBytes)
		copy(original, src)

		ca.Apply(intermediate, src)
		if bytes.Equal(intermediate, src) {
			t.Errorf("Apply did not change the block (highly unlikely)")
		}

		ca.Revert(reverted, intermediate)
		if !bytes.Equal(reverted, original) {
			t.Errorf("Revert failed: original=%v, reverted=%v", original, reverted)
		}
	}
}

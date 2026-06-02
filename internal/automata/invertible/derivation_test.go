package invertible_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/invertible"
	"github.com/staleread/aquila/internal/automata/state"
)

func TestDeriveGeneralCA(t *testing.T) {
	ca := invertible.NewCA()
	rng := rand.New(rand.NewSource(42))

	if err := ca.Generate(rng); err != nil {
		t.Fatalf("Failed to generate CA: %v", err)
	}

	gen, err := ca.DeriveGeneralCA()
	if err != nil {
		t.Fatalf("Failed to derive general CA: %v", err)
	}

	src := make([]byte, state.StateBytes)
	dstInvertible := make([]byte, state.StateBytes)
	dstGeneral := make([]byte, state.StateBytes)

	for i := range 5 {
		rng.Read(src)

		ca.Apply(dstInvertible, src)
		gen.Apply(dstGeneral, src)

		if !bytes.Equal(dstInvertible, dstGeneral) {
			t.Errorf("Inconsistency found at iteration %d:\nInvertible Apply: %x\nGeneral Apply:    %x", i, dstInvertible, dstGeneral)
		}
	}
}

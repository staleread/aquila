package invertible_test

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/invertible"
)

type zeroReader struct{}

func (zeroReader) Read(dst []byte) (int, error) {
	clear(dst)
	return len(dst), nil
}

func TestCAGenerateRejectsZeroAffineShift(t *testing.T) {
	prefix := make([]byte, 2*invertible.StateBytes)
	prefix[invertible.StateBytes] = 1
	rng := rand.New(rand.NewSource(42))
	reader := io.MultiReader(bytes.NewReader(prefix), rng)

	ca := invertible.NewCA()
	if err := ca.Generate(reader); err != nil {
		t.Fatalf("Generate failed after a nonzero shift was supplied: %v", err)
	}

	zero := make([]byte, invertible.StateBytes)
	got := make([]byte, invertible.StateBytes)
	ca.Apply(got, zero)
	if bytes.Equal(got, zero) {
		t.Fatal("generated CA maps zero to zero")
	}
}

func TestCAGenerateRejectsRepeatedZeroAffineShifts(t *testing.T) {
	ca := invertible.NewCA()
	if err := ca.Generate(zeroReader{}); err == nil {
		t.Fatal("Generate accepted repeated zero affine shifts")
	}
}

func TestCAInvertibility(t *testing.T) {
	ca := invertible.NewCA()
	rng := rand.New(rand.NewSource(42))

	if err := ca.Generate(rng); err != nil {
		t.Fatalf("Failed to generate CA: %v", err)
	}

	src := make([]byte, invertible.StateBytes)
	intermediate := make([]byte, invertible.StateBytes)
	reverted := make([]byte, invertible.StateBytes)

	for range 10 {
		rng.Read(src)

		original := make([]byte, invertible.StateBytes)
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

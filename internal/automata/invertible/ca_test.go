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
	validShift := make([]byte, invertible.StateBytes)
	validShift[0] = 1
	reader := io.MultiReader(
		bytes.NewReader(make([]byte, invertible.StateBytes)), // used for V generation and should get rejected as a zero-vector
		bytes.NewReader(validShift),                          // non zero vector and should be accepted as V
		rand.New(rand.NewSource(42)),                         // used for other CA components generation
	)

	ca := invertible.NewCA()
	if err := ca.Generate(reader); err != nil {
		t.Fatalf("Generate failed after a nonzero shift was supplied: %v", err)
	}

	zero := make([]byte, invertible.StateBytes)
	applied := make([]byte, invertible.StateBytes)
	ca.Apply(applied, zero)
	if !bytes.Equal(applied, validShift) {
		t.Errorf("Apply(zero) = %x, want generated affine shift %x", applied, validShift)
	}

	reverted := make([]byte, invertible.StateBytes)
	ca.Revert(reverted, applied)
	if !bytes.Equal(reverted, zero) {
		t.Errorf("Revert(Apply(zero)) = %x, want zero", reverted)
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

	initial := make([]byte, invertible.StateBytes)
	applied := make([]byte, invertible.StateBytes)
	reverted := make([]byte, invertible.StateBytes)

	for range 10 {
		rng.Read(initial)

		original := make([]byte, invertible.StateBytes)
		copy(original, initial)

		ca.Apply(applied, initial)
		if bytes.Equal(applied, initial) {
			t.Errorf("Apply did not change the block (highly unlikely)")
		}

		ca.Revert(reverted, applied)
		if !bytes.Equal(reverted, original) {
			t.Errorf("Revert failed: original=%v, reverted=%v", original, reverted)
		}
	}
}

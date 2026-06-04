package general_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestCACompilationCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping CA compilation test in short mode")
	}

	rng := rand.New(rand.NewSource(42))
	invCA := invertible.NewCA()
	if err := invCA.Generate(rng); err != nil {
		t.Fatalf("Failed to generate invertible CA: %v", err)
	}

	genCA, err := invCA.DeriveGeneralCA()
	if err != nil {
		t.Fatalf("Failed to compile CA: %v", err)
	}

	src := make([]byte, invertible.StateBytes)
	got := make([]byte, invertible.StateBytes)
	want := make([]byte, invertible.StateBytes)

	for range 5 {
		rng.Read(src)

		invCA.Apply(want, src)
		genCA.Apply(got, src)

		if !bytes.Equal(got, want) {
			t.Errorf("Compiled CA output mismatch!\ngot:  %x\nwant: %x", got, want)
		}
	}
}

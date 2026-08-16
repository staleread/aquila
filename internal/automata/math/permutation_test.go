package math_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}

func TestPermutationGenerateRejectsBiasedByte(t *testing.T) {
	var arena [math.PermutationBytes]byte
	var entropy [math.PermutationBytes]byte

	index := 0
	bound := math.PermutationSize
	for 256%bound == 0 {
		index++
		bound--
	}
	entropy[index] = 255 // Rejected because the bound does not divide 256.
	entropy[index+1] = 0 // Accepted and selects the current index.

	perm := math.NewPermutation(arena[:])
	if err := perm.Generate(bytes.NewReader(entropy[:])); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if got := perm.Data[index]; got != math.Subscript(index) {
		t.Fatalf("rejected byte selected index: got %d, want %d", got, index)
	}
}

func TestPermutationGeneratePropagatesReaderError(t *testing.T) {
	var arena [math.PermutationBytes]byte
	errExpected := errors.New("entropy unavailable")

	err := math.NewPermutation(arena[:]).Generate(errorReader{errExpected})
	if !errors.Is(err, errExpected) {
		t.Fatalf("Generate error = %v, want %v", err, errExpected)
	}
}

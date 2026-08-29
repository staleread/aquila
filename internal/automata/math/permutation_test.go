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

func TestPermutationGatherScatterUsesLocalToGlobalDirection(t *testing.T) {
	var arena [math.PermutationBytes]byte
	perm := math.NewPermutation(arena[:])

	for i := range perm.Data {
		perm.Data[i] = math.Subscript(i)
	}

	// A 3-cycle, not a swap: pi(0)=2, pi(1)=0, pi(2)=1.
	perm.Data[0] = 2
	perm.Data[1] = 0
	perm.Data[2] = 1

	// Global coordinate 0 occurs at local coordinate 1, since pi(1)=0.
	input := math.NewBitset(0)
	got := perm.Gather(&input, 0)
	want := math.Vector(1 << 1)
	if got != want {
		t.Fatalf("Gather = %b, want %b", got, want)
	}

	// Conversely, local coordinate 1 must scatter to global coordinate 0.
	var output math.Bitset
	perm.Scatter(&output, 0, want)
	if output.Compare(input) != 0 {
		t.Fatalf("Scatter result differs from expected global coordinate")
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

//go:build deg3

package math_test

import (
	"bytes"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
	"github.com/staleread/aquila/internal/automata/state"
)

func TestConfusionMapEval(t *testing.T) {
	perm := getIdentityPermutation()

	degArena := make([]byte, math.ConfusionMapBytes)
	m := math.NewConfusionMap(degArena)

	// y0 = x0*x1*x2 + x3*x4
	fixture := make([]byte, math.ConfusionMapBytes)
	fixture[0], fixture[1], fixture[2] = 0, 1, 2
	fixture[3], fixture[4] = 3, 4
	fixtureGen := bytes.NewReader(fixture)

	maxSub := state.Subscript(math.VectorSize)

	if err := m.Generate(fixtureGen, maxSub, perm); err != nil {
		t.Fatalf("Failed to generate degusion map: %v", err)
	}

	// Verify initialization
	expectedIds := []byte{0, 1, 2, 3, 4}
	for i, exp := range expectedIds {
		if degArena[i] != exp {
			t.Fatalf("Index at %d: got %d, expected %d", i, degArena[i], exp)
		}
	}

	tests := []struct {
		name        string
		activeBits  []state.Subscript
		expectedBit uint8
	}{
		{
			name:        "All zeros -> 0^0 = 0",
			activeBits:  []state.Subscript{},
			expectedBit: 0,
		},
		{
			name:        "Only bit 5 is set -> 0^0 = 0",
			activeBits:  []state.Subscript{5},
			expectedBit: 0,
		},
		{
			name:        "Bits 3 and 4 are set -> 0^1 = 1",
			activeBits:  []state.Subscript{3, 4},
			expectedBit: 1,
		},
		{
			name:        "Bits 0, 1, 2, 3, 4 are set -> 1^1 = 0",
			activeBits:  []state.Subscript{0, 1, 2, 3, 4},
			expectedBit: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := state.State{}

			for _, bit := range tc.activeBits {
				b.SetAt(bit, 1)
			}

			actual := m.Eval(b)
			actualBit := uint8(actual & 1)

			if actualBit != tc.expectedBit {
				t.Errorf("Expected %d, got %d", tc.expectedBit, actualBit)
			}
		})
	}
}

func getIdentityPermutation() *math.Permutation {
	permArena := make([]byte, math.PermutationBytes)

	for i := range len(permArena) {
		permArena[i] = state.Subscript(i)
	}
	return math.NewPermutation(permArena)
}

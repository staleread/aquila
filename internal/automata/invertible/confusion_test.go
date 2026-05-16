package invertible_test

import (
	"bytes"
	"testing"

	"github.com/staleread/aquila/internal/automata"
	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestConfusionMapEval(t *testing.T) {
	perm := getIdentityPermutation()

	confArena := make([]byte, invertible.ConfusionMapBytes)
	m := invertible.NewConfusionMap(confArena)

	// y0 = x0*x1*x2 + x3*x4
	fixture := make([]byte, invertible.ConfusionMapBytes)
	fixture[0], fixture[1], fixture[2] = 0, 1, 2
	fixture[3], fixture[4] = 3, 4
	fixtureGen := bytes.NewReader(fixture)

	maxSub := invertible.Subscript(invertible.VectorSize)

	if err := m.Generate(fixtureGen, maxSub, perm); err != nil {
		t.Fatalf("Failed to generate confusion map: %v", err)
	}

	// Verify initialization
	expectedIds := []byte{0, 1, 2, 3, 4}
	for i, exp := range expectedIds {
		if confArena[i] != exp {
			t.Fatalf("Index at %d: got %d, expected %d", i, confArena[i], exp)
		}
	}

	tests := []struct {
		name        string
		activeBits  []uint8
		expectedBit uint8
	}{
		{
			name:        "All zeros -> 0^0 = 0",
			activeBits:  []uint8{},
			expectedBit: 0,
		},
		{
			name:        "Only bit 5 is set -> 0^0 = 0",
			activeBits:  []uint8{5},
			expectedBit: 0,
		},
		{
			name:        "Bits 3 and 4 are set -> 0^1 = 1",
			activeBits:  []uint8{3, 4},
			expectedBit: 1,
		},
		{
			name:        "Bits 0, 1, 2, 3, 4 are set -> 1^1 = 0",
			activeBits:  []uint8{0, 1, 2, 3, 4},
			expectedBit: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := &automata.Block{}

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

func getIdentityPermutation() *invertible.Permutation {
	permArena := make([]byte, invertible.PermutationBytes)

	for i := range len(permArena) {
		permArena[i] = invertible.Subscript(i)
	}
	return invertible.NewPermutation(permArena)
}

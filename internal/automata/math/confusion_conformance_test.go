package math_test

import (
	"bytes"
	"slices"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestConfusionMapGenerateUsesPermutedPredecessorCoordinates(t *testing.T) {
	permArena := make([]byte, math.PermutationBytes)
	perm := math.NewPermutation(permArena)
	for i := range perm.Data {
		perm.Data[i] = math.Subscript(i)
	}

	// A nonidentity local-to-global mapping for predecessor coordinates.
	perm.Data[0] = 2
	perm.Data[1] = 0
	perm.Data[2] = 1

	maxSub := math.Subscript(math.VectorSize)
	arena := make([]byte, math.ConfusionMapBytes)
	fixture := make([]byte, math.ConfusionMapBytes)
	for i := range math.ConfusionDegree {
		fixture[i] = byte(i)
	}

	m := math.NewConfusionMap(arena)
	if err := m.Generate(bytes.NewReader(fixture), maxSub, perm); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// The first generated monomial chooses local predecessors 0..d-1 and must
	// store their corresponding global coordinates pi(0)..pi(d-1).
	for i := range math.ConfusionDegree {
		if got, want := m.Data[i], perm.Data[i]; got != want {
			t.Fatalf("first monomial index %d = %d, want global coordinate %d", i, got, want)
		}
	}

	for _, sub := range m.Data {
		if !slices.Contains(perm.Data[:maxSub], sub) {
			t.Fatalf("nonlinear index %d is not a global image of a predecessor coordinate", sub)
		}
	}

	// Eval receives global state coordinates. Keep only the first monomial
	// active, proving that its permutation-translated indices are used.
	inactive := perm.Data[maxSub-1]
	for i := math.ConfusionDegree; i < len(m.Data); i++ {
		m.Data[i] = inactive
	}
	active := math.NewBitset(perm.Data[:math.ConfusionDegree]...)
	if got := m.Eval(active) & 1; got != 1 {
		t.Fatalf("Eval first output bit = %d, want 1 with its global predecessors active", got)
	}
}

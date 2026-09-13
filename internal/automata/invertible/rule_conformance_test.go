//go:build !comp0

package invertible

import (
	"bytes"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestRulesUseLocalPermutationOrderAndAscendingComposition(t *testing.T) {
	ca := NewCA()
	for i := range RulesCount {
		initializeIdentityRule(ca.getRule(i), t)
	}

	// R_0 maps global x_0 to x_0 xor x_1.
	setFirstFoldUpperShear(ca.getRule(0), t)

	// R_1 has the same first-fold coordinate set as R_0, but its local order is
	// rotated: local 0,1 correspond to global 1,2. The same local shear thus
	// maps global x_1 to x_1 xor x_2.
	rule1 := ca.getRule(1)
	perm := rule1.getPermutation()
	perm.Data[0] = 1
	perm.Data[1] = 2
	perm.Data[2] = 0
	setFirstFoldUpperShear(rule1, t)

	// With input x_2, ascending R_0 then R_1 produces x_1 xor x_2. Reversing
	// the Rule order would additionally set x_0.
	input := make([]byte, StateBytes)
	input[0] = 1 << 2
	got := make([]byte, StateBytes)
	ca.Apply(got, input)
	if got[0] != 0b00000110 {
		t.Fatalf("Apply(x_2) first byte = %08b, want 00000110", got[0])
	}
	for i, b := range got[1:] {
		if b != 0 {
			t.Fatalf("Apply(x_2) byte %d = %08b, want 0", i+1, b)
		}
	}

	reverted := make([]byte, StateBytes)
	ca.Revert(reverted, got)
	if !bytes.Equal(reverted, input) {
		t.Fatalf("Revert(Apply(x_2)) = %08b, want %08b", reverted, input)
	}
}

func initializeIdentityRule(rule Rule, t *testing.T) {
	t.Helper()

	perm := rule.getPermutation()
	for i := range perm.Data {
		perm.Data[i] = math.Subscript(i)
	}
	for i := range FoldsCount {
		fold := rule.getFold(i)
		if err := fold.sle.Generate(bytes.NewReader(make([]byte, math.SLEBytes))); err != nil {
			t.Fatalf("initialize identity linear fold %d: %v", i, err)
		}
	}
}

func setFirstFoldUpperShear(rule Rule, t *testing.T) {
	t.Helper()

	// The first SLE row is an upper-triangular equation x_0 xor x_1. Generate
	// then adds all diagonal entries, yielding an invertible unit-LU factor.
	arena := make([]byte, math.SLEBytes)
	arena[0] = 1 << 1
	fold := rule.getFold(0)
	if err := fold.sle.Generate(bytes.NewReader(arena)); err != nil {
		t.Fatalf("initialize upper shear: %v", err)
	}
}

package core_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestBlock96(t *testing.T) {
	var b core.Block96

	if b.OnesCount() != 0 {
		t.Errorf("Expected 0 ones, got %d", b.OnesCount())
	}

	b.SetAt(10, 1)
	b.SetAt(40, 1)
	b.SetAt(95, 1)

	if b.At(10) != 1 || b.At(40) != 1 || b.At(95) != 1 {
		t.Errorf("Bits not set correctly")
	}

	if b.OnesCount() != 3 {
		t.Errorf("Expected 3 ones, got %d", b.OnesCount())
	}

	if b.TrailingZeros() != 10 {
		t.Errorf("Expected trailing zeros be 10, got %d", b.TrailingZeros())
	}

	if b.TrailingZerosAfter(11) != 40 {
		t.Errorf("Expected next trailing zeros be 40, got %d", b.TrailingZerosAfter(11))
	}

	if b.TrailingZerosAfter(41) != 95 {
		t.Errorf("Expected next trailing zeros be 95, got %d", b.TrailingZerosAfter(41))
	}

	if b.TrailingZerosAfter(96) != 96 {
		t.Errorf("Expected out of bounds next trailing zeros be 96, got %d", b.TrailingZerosAfter(96))
	}

	var other core.Block96
	other.SetAt(10, 1)

	if !b.Contains(other) {
		t.Errorf("b should contain other")
	}

	other.SetAt(20, 1)
	if b.Contains(other) {
		t.Errorf("b should not contain other")
	}
}

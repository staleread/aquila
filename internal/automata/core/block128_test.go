package core_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestBlock128(t *testing.T) {
	var b core.Block128

	if b.OnesCount() != 0 {
		t.Errorf("Expected 0 ones, got %d", b.OnesCount())
	}

	b.SetAt(10, 1)
	b.SetAt(70, 1)
	b.SetAt(127, 1)

	if b.At(10) != 1 || b.At(70) != 1 || b.At(127) != 1 {
		t.Errorf("Bits not set correctly")
	}

	if b.OnesCount() != 3 {
		t.Errorf("Expected 3 ones, got %d", b.OnesCount())
	}

	if b.TrailingZeros() != 10 {
		t.Errorf("Expected trailing zeros be 10, got %d", b.TrailingZeros())
	}

	if b.TrailingZerosAfter(11) != 70 {
		t.Errorf("Expected next trailing zeros be 70, got %d", b.TrailingZerosAfter(11))
	}

	if b.TrailingZerosAfter(71) != 127 {
		t.Errorf("Expected next trailing zeros be 127, got %d", b.TrailingZerosAfter(71))
	}

	if b.TrailingZerosAfter(128) != 128 {
		t.Errorf("Expected out of bounds next trailing zeros be 128, got %d", b.TrailingZerosAfter(128))
	}

	var other core.Block128
	other.SetAt(10, 1)

	if !b.Contains(other) {
		t.Errorf("b should contain other")
	}

	other.SetAt(20, 1)
	if b.Contains(other) {
		t.Errorf("b should not contain other")
	}
}

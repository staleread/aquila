package automata

import (
	"testing"
)

func TestBlockSetAt(t *testing.T) {
	var b Block

	idx := 90
	bit := Word(1)

	b.SetAt(idx, bit)

	if b.At(idx) != 1 {
		t.Errorf("Expected bit %d to be 1, got %d", idx, b.At(idx))
	}
}

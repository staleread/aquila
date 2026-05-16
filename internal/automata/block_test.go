package automata

import (
	"testing"
)

func TestBlock_SetAt(t *testing.T) {
	var b Block

	idx := uint8(120)
	bit := Bit(1)

	b.SetAt(idx, bit)

	if b.At(idx) != 1 {
		t.Errorf("Expected bit %d to be 1, got %d", idx, b.At(idx))
	}
}

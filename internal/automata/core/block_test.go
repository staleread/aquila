package core_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestBlockSetAt(t *testing.T) {
	var b core.Block

	idx := 90
	bit := core.Word(1)

	b.SetAt(idx, bit)

	if b.At(idx) != 1 {
		t.Errorf("Expected bit %d to be 1, got %d", idx, b.At(idx))
	}
}

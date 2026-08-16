package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestBitsetIsZero(t *testing.T) {
	var zero math.Bitset
	if !zero.IsZero() {
		t.Fatal("zero bitset reported as nonzero")
	}

	nonzero := math.NewBitset(0)
	if nonzero.IsZero() {
		t.Fatal("nonzero bitset reported as zero")
	}
}

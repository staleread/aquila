//go:build block96 || (!block8 && !block16 && !block32 && !block64 && !block128)

package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestBitset_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
		readBit math.Subscript
		want    uint8
	}{
		{"first bit of first word", []math.Subscript{0}, 0, 1},
		{"last bit of first word", []math.Subscript{31}, 31, 1},
		{"first bit of second word", []math.Subscript{32}, 32, 1},
		{"first bit of third word", []math.Subscript{64}, 64, 1},
		{"last bit of third word", []math.Subscript{95}, 95, 1},
		{"unset bit reads zero", []math.Subscript{5}, 10, 0},
		{"multiple bits set, read one", []math.Subscript{3, 40, 80}, 40, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := math.NewBitset(tt.setBits...)
			got := s.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBitset_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero math", nil},
		{"single bit word 0", []math.Subscript{7}},
		{"single bit word 1", []math.Subscript{40}},
		{"single bit word 2", []math.Subscript{80}},
		{"bits across all words", []math.Subscript{10, 40, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := math.NewBitset(tt.setBits...)

			buf := make([]byte, 12)
			src.Write(buf)

			var dst math.Bitset
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBitset_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b math.Bitset
		want int
	}{
		{"equal", math.NewBitset(5), math.NewBitset(5), 0},
		{"a < b first word", math.NewBitset(0), math.NewBitset(1), -1},
		{"a > b first word", math.NewBitset(1), math.NewBitset(0), 1},
		{"a < b second word", math.NewBitset(32), math.NewBitset(33), -1},
		{"a > b third word", math.NewBitset(65), math.NewBitset(64), 1},
		{"both zero", math.Bitset{}, math.Bitset{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Compare(tt.b)
			if got != tt.want {
				t.Errorf("Compare = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBitset_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     math.Bitset
		wantBits []math.Subscript
	}{
		{"disjoint bits", math.NewBitset(0), math.NewBitset(1), []math.Subscript{0, 1}},
		{"overlapping bits", math.NewBitset(0, 1), math.NewBitset(1, 2), []math.Subscript{0, 1, 2}},
		{"across words", math.NewBitset(31), math.NewBitset(32), []math.Subscript{31, 32}},
		{"with zero", math.NewBitset(10, 50), math.Bitset{}, []math.Subscript{10, 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := math.NewBitset(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBitset_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     math.Bitset
		other math.Bitset
		want  bool
	}{
		{"exact match", math.NewBitset(10, 50), math.NewBitset(10, 50), true},
		{"proper superset", math.NewBitset(10, 50, 80), math.NewBitset(10, 50), true},
		{"missing one bit", math.NewBitset(10), math.NewBitset(10, 50), false},
		{"other is zero", math.NewBitset(5), math.Bitset{}, true},
		{"both zero", math.Bitset{}, math.Bitset{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Contains(tt.other)
			if got != tt.want {
				t.Errorf("Contains = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitset_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero math", nil},
		{"single bit word 0", []math.Subscript{0}},
		{"single bit word 1", []math.Subscript{40}},
		{"single bit word 2", []math.Subscript{80}},
		{"multiple bits all words", []math.Subscript{0, 31, 32, 63, 64, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.Subscripts(nil)

			if len(got) != len(tt.setBits) {
				t.Fatalf("Subscripts len = %d, want %d", len(got), len(tt.setBits))
			}
			for i, s := range tt.setBits {
				if got[i] != s {
					t.Errorf("Subscripts[%d] = %d, want %d", i, got[i], s)
				}
			}
		})
	}
}

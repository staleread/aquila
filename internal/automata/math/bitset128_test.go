//go:build block128

package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/state"
)

func TestBitset_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
		readBit math.Subscript
		want    uint8
	}{
		{"first bit of first word", []math.Subscript{0}, 0, 1},
		{"last bit of first word", []math.Subscript{63}, 63, 1},
		{"first bit of second word", []math.Subscript{64}, 64, 1},
		{"last bit of second word", []math.Subscript{127}, 127, 1},
		{"unset bit reads zero", []math.Subscript{5}, 10, 0},
		{"multiple bits set, read one", []math.Subscript{3, 70, 120}, 70, 1},
		{"clear a previously set bit", nil, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.At(tt.readBit)
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
		{"zero state", nil},
		{"single bit word 0", []math.Subscript{7}},
		{"single bit word 1", []math.Subscript{100}},
		{"bits in both words", []math.Subscript{10, 70, 127}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := math.NewBitset(tt.setBits...)

			buf := make([]byte, 16)
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
		{"a < b same word", math.NewBitset(0), math.NewBitset(1), -1},
		{"a > b same word", math.NewBitset(1), math.NewBitset(0), 1},
		{"a < b second word", math.NewBitset(64), math.NewBitset(65), -1},
		{"a > b second word", math.NewBitset(65), math.NewBitset(64), 1},
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
		{"across words", math.NewBitset(63), math.NewBitset(64), []math.Subscript{63, 64}},
		{"with zero", math.NewBitset(10, 20), math.Bitset{}, []math.Subscript{10, 20}},
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
		{"exact match", math.NewBitset(10, 70), math.NewBitset(10, 70), true},
		{"proper superset", math.NewBitset(10, 70, 127), math.NewBitset(10, 70), true},
		{"missing one bit", math.NewBitset(10), math.NewBitset(10, 70), false},
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
		{"zero state", nil},
		{"single bit word 0", []math.Subscript{0}},
		{"single bit word 1", []math.Subscript{65}},
		{"multiple bits both words", []math.Subscript{0, 31, 64, 127}},
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

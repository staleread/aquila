//go:build block32

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
		{"first bit", []math.Subscript{0}, 0, 1},
		{"last bit", []math.Subscript{31}, 31, 1},
		{"unset bit reads zero", []math.Subscript{2}, 5, 0},
		{"multiple bits set, read one", []math.Subscript{1, 3, 20}, 3, 1},
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
		{"zero bitset", nil},
		{"single bit", []math.Subscript{20}},
		{"multiple bits", []math.Subscript{1, 16, 31}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := math.NewBitset(tt.setBits...)

			buf := make([]byte, 4)
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
		{"a < b", math.NewBitset(0), math.NewBitset(1), -1},
		{"a > b", math.NewBitset(1), math.NewBitset(0), 1},
		{"both zero", math.Bitset(0), math.Bitset(0), 0},
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
		{"with zero", math.NewBitset(2, 28), math.Bitset(0), []math.Subscript{2, 28}},
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
		{"exact match", math.NewBitset(2, 28), math.NewBitset(2, 28), true},
		{"proper superset", math.NewBitset(2, 4, 28), math.NewBitset(2, 28), true},
		{"missing one bit", math.NewBitset(2), math.NewBitset(2, 28), false},
		{"other is zero", math.NewBitset(5), math.Bitset(0), true},
		{"both zero", math.Bitset(0), math.Bitset(0), true},
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
		{"single bit", []math.Subscript{25}},
		{"multiple bits", []math.Subscript{0, 25, 31}},
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

package core_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestBlock96_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
		readBit core.Subscript
		want    uint8
	}{
		{"first bit of first word", []core.Subscript{0}, 0, 1},
		{"last bit of first word", []core.Subscript{31}, 31, 1},
		{"first bit of second word", []core.Subscript{32}, 32, 1},
		{"first bit of third word", []core.Subscript{64}, 64, 1},
		{"last bit of third word", []core.Subscript{95}, 95, 1},
		{"unset bit reads zero", []core.Subscript{5}, 10, 0},
		{"multiple bits set, read one", []core.Subscript{3, 40, 80}, 40, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block96With(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBlock96_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"zero block", nil},
		{"single bit word 0", []core.Subscript{7}},
		{"single bit word 1", []core.Subscript{40}},
		{"single bit word 2", []core.Subscript{80}},
		{"bits across all words", []core.Subscript{10, 40, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := block96With(tt.setBits...)

			buf := make([]byte, 12)
			src.Write(buf)

			var dst core.Block96
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBlock96_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b core.Block96
		want int
	}{
		{"equal", block96With(5), block96With(5), 0},
		{"a < b first word", block96With(0), block96With(1), -1},
		{"a > b first word", block96With(1), block96With(0), 1},
		{"a < b second word", block96With(32), block96With(33), -1},
		{"a > b third word", block96With(65), block96With(64), 1},
		{"both zero", core.Block96{}, core.Block96{}, 0},
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

func TestBlock96_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     core.Block96
		wantBits []core.Subscript
	}{
		{"disjoint bits", block96With(0), block96With(1), []core.Subscript{0, 1}},
		{"overlapping bits", block96With(0, 1), block96With(1, 2), []core.Subscript{0, 1, 2}},
		{"across words", block96With(31), block96With(32), []core.Subscript{31, 32}},
		{"with zero", block96With(10, 50), core.Block96{}, []core.Subscript{10, 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := block96With(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBlock96_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     core.Block96
		other core.Block96
		want  bool
	}{
		{"exact match", block96With(10, 50), block96With(10, 50), true},
		{"proper superset", block96With(10, 50, 80), block96With(10, 50), true},
		{"missing one bit", block96With(10), block96With(10, 50), false},
		{"other is zero", block96With(5), core.Block96{}, true},
		{"both zero", core.Block96{}, core.Block96{}, true},
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

func TestBlock96_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"empty block", nil},
		{"single bit word 0", []core.Subscript{0}},
		{"single bit word 1", []core.Subscript{40}},
		{"single bit word 2", []core.Subscript{80}},
		{"multiple bits all words", []core.Subscript{0, 31, 32, 63, 64, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block96With(tt.setBits...)
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

func block96With(subs ...core.Subscript) core.Block96 {
	var b core.Block96
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

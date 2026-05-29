package core_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestBlock128_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
		readBit core.Subscript
		want    uint8
	}{
		{"first bit of first word", []core.Subscript{0}, 0, 1},
		{"last bit of first word", []core.Subscript{63}, 63, 1},
		{"first bit of second word", []core.Subscript{64}, 64, 1},
		{"last bit of second word", []core.Subscript{127}, 127, 1},
		{"unset bit reads zero", []core.Subscript{5}, 10, 0},
		{"multiple bits set, read one", []core.Subscript{3, 70, 120}, 70, 1},
		{"clear a previously set bit", nil, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block128With(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBlock128_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"zero block", nil},
		{"single bit word 0", []core.Subscript{7}},
		{"single bit word 1", []core.Subscript{100}},
		{"bits in both words", []core.Subscript{10, 70, 127}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := block128With(tt.setBits...)

			buf := make([]byte, 16)
			src.Write(buf)

			var dst core.Block128
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBlock128_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b core.Block128
		want int
	}{
		{"equal", block128With(5), block128With(5), 0},
		{"a < b same word", block128With(0), block128With(1), -1},
		{"a > b same word", block128With(1), block128With(0), 1},
		{"a < b second word", block128With(64), block128With(65), -1},
		{"a > b second word", block128With(65), block128With(64), 1},
		{"both zero", core.Block128{}, core.Block128{}, 0},
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

func TestBlock128_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     core.Block128
		wantBits []core.Subscript
	}{
		{"disjoint bits", block128With(0), block128With(1), []core.Subscript{0, 1}},
		{"overlapping bits", block128With(0, 1), block128With(1, 2), []core.Subscript{0, 1, 2}},
		{"across words", block128With(63), block128With(64), []core.Subscript{63, 64}},
		{"with zero", block128With(10, 20), core.Block128{}, []core.Subscript{10, 20}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := block128With(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBlock128_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     core.Block128
		other core.Block128
		want  bool
	}{
		{"exact match", block128With(10, 70), block128With(10, 70), true},
		{"proper superset", block128With(10, 70, 127), block128With(10, 70), true},
		{"missing one bit", block128With(10), block128With(10, 70), false},
		{"other is zero", block128With(5), core.Block128{}, true},
		{"both zero", core.Block128{}, core.Block128{}, true},
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

func TestBlock128_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"empty block", nil},
		{"single bit word 0", []core.Subscript{0}},
		{"single bit word 1", []core.Subscript{65}},
		{"multiple bits both words", []core.Subscript{0, 31, 64, 127}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block128With(tt.setBits...)
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

func block128With(subs ...core.Subscript) core.Block128 {
	var b core.Block128
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

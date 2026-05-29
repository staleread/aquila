package core_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestBlock16_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
		readBit core.Subscript
		want    uint8
	}{
		{"first bit", []core.Subscript{0}, 0, 1},
		{"last bit", []core.Subscript{15}, 15, 1},
		{"unset bit reads zero", []core.Subscript{2}, 5, 0},
		{"multiple bits set, read one", []core.Subscript{1, 3, 10}, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block16With(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBlock16_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"zero block", nil},
		{"single bit", []core.Subscript{10}},
		{"multiple bits", []core.Subscript{1, 8, 15}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := block16With(tt.setBits...)

			buf := make([]byte, 2)
			src.Write(buf)

			var dst core.Block16
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBlock16_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b core.Block16
		want int
	}{
		{"equal", block16With(5), block16With(5), 0},
		{"a < b", block16With(0), block16With(1), -1},
		{"a > b", block16With(1), block16With(0), 1},
		{"both zero", core.Block16(0), core.Block16(0), 0},
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

func TestBlock16_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     core.Block16
		wantBits []core.Subscript
	}{
		{"disjoint bits", block16With(0), block16With(1), []core.Subscript{0, 1}},
		{"overlapping bits", block16With(0, 1), block16With(1, 2), []core.Subscript{0, 1, 2}},
		{"with zero", block16With(2, 12), core.Block16(0), []core.Subscript{2, 12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := block16With(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBlock16_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     core.Block16
		other core.Block16
		want  bool
	}{
		{"exact match", block16With(2, 12), block16With(2, 12), true},
		{"proper superset", block16With(2, 4, 12), block16With(2, 12), true},
		{"missing one bit", block16With(2), block16With(2, 12), false},
		{"other is zero", block16With(5), core.Block16(0), true},
		{"both zero", core.Block16(0), core.Block16(0), true},
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

func TestBlock16_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"empty block", nil},
		{"single bit", []core.Subscript{7}},
		{"multiple bits", []core.Subscript{0, 7, 15}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block16With(tt.setBits...)
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

func block16With(subs ...core.Subscript) core.Block16 {
	var b core.Block16
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

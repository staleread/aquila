package core_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestBlock32_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
		readBit core.Subscript
		want    uint8
	}{
		{"first bit", []core.Subscript{0}, 0, 1},
		{"last bit", []core.Subscript{31}, 31, 1},
		{"unset bit reads zero", []core.Subscript{2}, 5, 0},
		{"multiple bits set, read one", []core.Subscript{1, 3, 20}, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block32With(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBlock32_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"zero block", nil},
		{"single bit", []core.Subscript{20}},
		{"multiple bits", []core.Subscript{1, 16, 31}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := block32With(tt.setBits...)

			buf := make([]byte, 4)
			src.Write(buf)

			var dst core.Block32
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBlock32_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b core.Block32
		want int
	}{
		{"equal", block32With(5), block32With(5), 0},
		{"a < b", block32With(0), block32With(1), -1},
		{"a > b", block32With(1), block32With(0), 1},
		{"both zero", core.Block32(0), core.Block32(0), 0},
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

func TestBlock32_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     core.Block32
		wantBits []core.Subscript
	}{
		{"disjoint bits", block32With(0), block32With(1), []core.Subscript{0, 1}},
		{"overlapping bits", block32With(0, 1), block32With(1, 2), []core.Subscript{0, 1, 2}},
		{"with zero", block32With(2, 28), core.Block32(0), []core.Subscript{2, 28}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := block32With(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBlock32_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     core.Block32
		other core.Block32
		want  bool
	}{
		{"exact match", block32With(2, 28), block32With(2, 28), true},
		{"proper superset", block32With(2, 4, 28), block32With(2, 28), true},
		{"missing one bit", block32With(2), block32With(2, 28), false},
		{"other is zero", block32With(5), core.Block32(0), true},
		{"both zero", core.Block32(0), core.Block32(0), true},
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

func TestBlock32_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []core.Subscript
	}{
		{"empty block", nil},
		{"single bit", []core.Subscript{25}},
		{"multiple bits", []core.Subscript{0, 25, 31}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := block32With(tt.setBits...)
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

func block32With(subs ...core.Subscript) core.Block32 {
	var b core.Block32
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

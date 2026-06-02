//go:build block16

package state_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/state"
)

func TestState_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []state.Subscript
		readBit state.Subscript
		want    uint8
	}{
		{"first bit", []state.Subscript{0}, 0, 1},
		{"last bit", []state.Subscript{15}, 15, 1},
		{"unset bit reads zero", []state.Subscript{2}, 5, 0},
		{"multiple bits set, read one", []state.Subscript{1, 3, 10}, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := state.NewState(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestState_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []state.Subscript
	}{
		{"zero state", nil},
		{"single bit", []state.Subscript{10}},
		{"multiple bits", []state.Subscript{1, 8, 15}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := state.NewState(tt.setBits...)

			buf := make([]byte, 2)
			src.Write(buf)

			var dst state.State
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestState_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b state.State
		want int
	}{
		{"equal", state.NewState(5), state.NewState(5), 0},
		{"a < b", state.NewState(0), state.NewState(1), -1},
		{"a > b", state.NewState(1), state.NewState(0), 1},
		{"both zero", state.State(0), state.State(0), 0},
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

func TestState_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     state.State
		wantBits []state.Subscript
	}{
		{"disjoint bits", state.NewState(0), state.NewState(1), []state.Subscript{0, 1}},
		{"overlapping bits", state.NewState(0, 1), state.NewState(1, 2), []state.Subscript{0, 1, 2}},
		{"with zero", state.NewState(2, 12), state.State(0), []state.Subscript{2, 12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := state.NewState(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestState_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     state.State
		other state.State
		want  bool
	}{
		{"exact match", state.NewState(2, 12), state.NewState(2, 12), true},
		{"proper superset", state.NewState(2, 4, 12), state.NewState(2, 12), true},
		{"missing one bit", state.NewState(2), state.NewState(2, 12), false},
		{"other is zero", state.NewState(5), state.State(0), true},
		{"both zero", state.State(0), state.State(0), true},
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

func TestState_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []state.Subscript
	}{
		{"zero state", nil},
		{"single bit", []state.Subscript{7}},
		{"multiple bits", []state.Subscript{0, 7, 15}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := state.NewState(tt.setBits...)
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

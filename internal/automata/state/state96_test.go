//go:build block96 || (!block16 && !block32 && !block128)

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
		{"first bit of first word", []state.Subscript{0}, 0, 1},
		{"last bit of first word", []state.Subscript{31}, 31, 1},
		{"first bit of second word", []state.Subscript{32}, 32, 1},
		{"first bit of third word", []state.Subscript{64}, 64, 1},
		{"last bit of third word", []state.Subscript{95}, 95, 1},
		{"unset bit reads zero", []state.Subscript{5}, 10, 0},
		{"multiple bits set, read one", []state.Subscript{3, 40, 80}, 40, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.NewState(tt.setBits...)
			got := s.At(tt.readBit)
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
		{"single bit word 0", []state.Subscript{7}},
		{"single bit word 1", []state.Subscript{40}},
		{"single bit word 2", []state.Subscript{80}},
		{"bits across all words", []state.Subscript{10, 40, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := state.NewState(tt.setBits...)

			buf := make([]byte, 12)
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
		{"a < b first word", state.NewState(0), state.NewState(1), -1},
		{"a > b first word", state.NewState(1), state.NewState(0), 1},
		{"a < b second word", state.NewState(32), state.NewState(33), -1},
		{"a > b third word", state.NewState(65), state.NewState(64), 1},
		{"both zero", state.State{}, state.State{}, 0},
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
		{"across words", state.NewState(31), state.NewState(32), []state.Subscript{31, 32}},
		{"with zero", state.NewState(10, 50), state.State{}, []state.Subscript{10, 50}},
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
		{"exact match", state.NewState(10, 50), state.NewState(10, 50), true},
		{"proper superset", state.NewState(10, 50, 80), state.NewState(10, 50), true},
		{"missing one bit", state.NewState(10), state.NewState(10, 50), false},
		{"other is zero", state.NewState(5), state.State{}, true},
		{"both zero", state.State{}, state.State{}, true},
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
		{"single bit word 0", []state.Subscript{0}},
		{"single bit word 1", []state.Subscript{40}},
		{"single bit word 2", []state.Subscript{80}},
		{"multiple bits all words", []state.Subscript{0, 31, 32, 63, 64, 95}},
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

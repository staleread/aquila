package general_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata"
	"github.com/staleread/aquila/internal/automata/general"
)

func TestNewMonomial(t *testing.T) {
	words := []automata.Word{1, 2, 3}
	m := general.NewMonomial(words)

	if m[0] != 1 || m[1] != 2 || m[2] != 3 {
		t.Errorf("NewMonomial failed: got %v, want %v", m, words)
	}
}

func TestCompareMonomials(t *testing.T) {
	tests := []struct {
		a, b general.Monomial
		want int
	}{
		{general.Monomial{1, 0, 0}, general.Monomial{1, 0, 0}, 0},
		{general.Monomial{1, 0, 0}, general.Monomial{2, 0, 0}, -1},
		{general.Monomial{2, 0, 0}, general.Monomial{1, 0, 0}, 1},
		{general.Monomial{0, 1, 0}, general.Monomial{0, 2, 0}, -1},
		{general.Monomial{0, 0, 1}, general.Monomial{0, 0, 0}, 1},
		{general.Monomial{1, 1, 0}, general.Monomial{1, 0, 1}, 1},
	}

	for _, tt := range tests {
		got := general.CompareMonomials(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("CompareMonomials(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMonomial_FirstSubscript(t *testing.T) {
	tests := []struct {
		m    general.Monomial
		want int
	}{
		{general.Monomial{0, 0, 0}, automata.BlockSize},
		{general.Monomial{1, 0, 0}, 0},
		{general.Monomial{1 << 31, 0, 0}, 31},
		{general.Monomial{0, 1, 0}, 32},
		{general.Monomial{0, 0, 1}, 64},
		{general.Monomial{0, 0, 1 << 31}, 95},
		{general.Monomial{0, 1 << 5, 1}, 37},
	}

	for _, tt := range tests {
		got := tt.m.FirstSubscript()
		if got != tt.want {
			t.Errorf("FirstSubscript(%v) = %d, want %d", tt.m, got, tt.want)
		}
	}
}

func TestMonomial_NextSubscript(t *testing.T) {
	// x0 * x5 * x34 * x63 * x64
	m := general.Monomial{1 | (1 << 5), (1 << 2) | (1 << 31), 1}

	expected := []int{0, 5, 34, 63, 64}

	curr := m.FirstSubscript()
	var got []int
	for curr < automata.BlockSize {
		got = append(got, curr)
		curr = m.NextSubscript(curr + 1)
	}

	if len(got) != len(expected) {
		t.Fatalf("NextSubscript iteration failed: got %v, want %v", got, expected)
	}

	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("At index %d: got %d, want %d", i, got[i], expected[i])
		}
	}
}

func TestMonomial_Mul(t *testing.T) {
	m1 := general.Monomial{1, 0, 0}
	m2 := general.Monomial{0, 2, 0}
	want := general.Monomial{1, 2, 0}

	got := m1.Mul(m2)
	if got != want {
		t.Errorf("Mul failed: got %v, want %v", got, want)
	}
}

func TestMonomial_Eval(t *testing.T) {
	// x0 * x31 * x33
	m := general.Monomial{1 | (1 << 31), 2, 0}

	tests := []struct {
		b    *automata.Block
		want automata.Bit
	}{
		{&automata.Block{1 | (1 << 31), 2, 0}, 1},                // x0 * x31 * x33
		{&automata.Block{1 | (1 << 31) | (1 << 5), 2 | 4, 1}, 1}, // x0 * x31 * x33 * x5 * x34
		{&automata.Block{1, 2, 0}, 0},                            // x0 * x33
		{&automata.Block{1 | (1 << 31), 0, 0}, 0},                // x0 * x31
	}

	for _, tt := range tests {
		got := m.Eval(tt.b)
		if got != tt.want {
			t.Errorf("Eval(%v) = %d, want %d", tt.b, got, tt.want)
		}
	}
}

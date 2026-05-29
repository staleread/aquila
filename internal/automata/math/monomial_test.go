package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
	"github.com/staleread/aquila/internal/automata/math"
)

func TestNewMonomial(t *testing.T) {
	subs := []core.Subscript{0, 1, 2}
	m := math.NewMonomial(subs...)

	got := m.Subscripts(nil)
	if len(got) != len(subs) {
		t.Fatalf("NewMonomial subscripts length: got %d, want %d", len(got), len(subs))
	}
	for i, s := range subs {
		if got[i] != s {
			t.Errorf("NewMonomial subscript[%d]: got %d, want %d", i, got[i], s)
		}
	}
}

func TestCompareMonomials(t *testing.T) {
	tests := []struct {
		name string
		a, b math.Monomial
		want int
	}{
		{
			name: "equal",
			a:    math.NewMonomial(0),
			b:    math.NewMonomial(0),
			want: 0,
		},
		{
			name: "a < b (lower bit set in a)",
			a:    math.NewMonomial(0),
			b:    math.NewMonomial(1),
			want: -1,
		},
		{
			name: "a > b (higher bit set in a)",
			a:    math.NewMonomial(1),
			b:    math.NewMonomial(0),
			want: 1,
		},
		{
			name: "a < b (second word differs)",
			a:    math.NewMonomial(64),
			b:    math.NewMonomial(65),
			want: -1,
		},
		{
			name: "a > b (empty vs non-empty second word)",
			a:    math.NewMonomial(64),
			b:    math.NewMonomial(),
			want: 1,
		},
		{
			name: "a > b (higher bit wins over lower bit in same word)",
			a:    math.NewMonomial(0, 1),
			b:    math.NewMonomial(0, 2),
			want: -1,
		},
	}

	for _, tt := range tests {
		got := math.CompareMonomials(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("CompareMonomials %s: got %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestMonomial_Mul(t *testing.T) {
	// x0 * x2 = {0, 2}
	m1 := math.NewMonomial(0)
	m2 := math.NewMonomial(2)
	want := math.NewMonomial(0, 2)

	got := m1.Mul(m2)
	if math.CompareMonomials(got, want) != 0 {
		t.Errorf("Mul(%v, %v) = %v, want %v", m1, m2, got, want)
	}
}

func TestMonomial_Eval(t *testing.T) {
	// monomial: x0 * x31 * x65  (spans both 64-bit words)
	m := math.NewMonomial(0, 31, 65)

	b0 := blockWith(0, 31, 65)        // exact match
	b1 := blockWith(0, 31, 65, 5, 66) // superset
	b2 := blockWith(0, 65)            // missing x31
	b3 := blockWith(0, 31)            // missing x65

	tests := []struct {
		name string
		b    core.Block
		want uint8
	}{
		{"exact match", b0, 1},
		{"superset", b1, 1},
		{"missing x31", b2, 0},
		{"missing x65", b3, 0},
	}

	for _, tt := range tests {
		got := m.Eval(tt.b)
		if got != tt.want {
			t.Errorf("Eval %s: got %d, want %d", tt.name, got, tt.want)
		}
	}
}

func blockWith(subs ...core.Subscript) core.Block {
	var b core.Block
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

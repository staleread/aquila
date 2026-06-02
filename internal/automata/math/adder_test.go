package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestAddPolynomials(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   []math.Monomial
		expected []math.Monomial
	}{
		{
			name:     "Both empty",
			p1:       []math.Monomial{},
			p2:       []math.Monomial{},
			expected: []math.Monomial{},
		},
		{
			name:     "One empty",
			p1:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			p2:       []math.Monomial{},
			expected: []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
		},
		{
			name:     "Disjoint sets",
			p1:       []math.Monomial{math.NewMonomial(3), math.NewMonomial(1)},
			p2:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			expected: []math.Monomial{math.NewMonomial(3), math.NewMonomial(2), math.NewMonomial(1), math.NewMonomial(0)},
		},
		{
			name:     "Identical (full cancellation)",
			p1:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			p2:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			expected: []math.Monomial{},
		},
		{
			name:     "Partial overlap",
			p1:       []math.Monomial{math.NewMonomial(3), math.NewMonomial(2)},
			p2:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			expected: []math.Monomial{math.NewMonomial(3), math.NewMonomial(0)}, // x2 cancels
		},
		{
			name:     "Interleaved",
			p1:       []math.Monomial{math.NewMonomial(4), math.NewMonomial(2), math.NewMonomial(0)},
			p2:       []math.Monomial{math.NewMonomial(3), math.NewMonomial(2), math.NewMonomial(1)},
			expected: []math.Monomial{math.NewMonomial(4), math.NewMonomial(3), math.NewMonomial(1), math.NewMonomial(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := math.AddPolynomials(nil, tt.p1, tt.p2)

			if len(got) != len(tt.expected) {
				t.Fatalf("Expected length %d, got %d. Got: %v", len(tt.expected), len(got), got)
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("At index %d: got %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

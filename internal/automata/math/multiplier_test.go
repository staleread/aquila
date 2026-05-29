package math_test

import (
	"slices"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func sortPolynomial(p []math.Monomial) {
	slices.SortFunc(p, func(a, b math.Monomial) int {
		return math.CompareMonomials(b, a) // Descending
	})
}

func TestPolynomialMultiplier_Multiply(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   []math.Monomial
		expected []math.Monomial
	}{
		{
			name:     "Empty polynomials",
			p1:       []math.Monomial{},
			p2:       []math.Monomial{math.NewMonomial(0)},
			expected: []math.Monomial{},
		},
		{
			name:     "Single monomial multiplication",
			p1:       []math.Monomial{math.NewMonomial(0)},    // x0
			p2:       []math.Monomial{math.NewMonomial(1)},    // x1
			expected: []math.Monomial{math.NewMonomial(0, 1)}, // x0*x1
		},
		{
			name: "Basic distribution (x0 + x1) * x2",
			p1:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)}, // x1, x0 (descending)
			p2:   []math.Monomial{math.NewMonomial(2)},                      // x2
			// (x1+x0)*x2 = x1x2 + x0x2
			expected: []math.Monomial{math.NewMonomial(1, 2), math.NewMonomial(0, 2)},
		},
		{
			name: "Cancellation x0 * (x0 + x1) = x0 + x0x1",
			p1:   []math.Monomial{math.NewMonomial(0)},
			p2:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)}, // x1, x0
			// x0*x1 = x0x1, x0*x0 = x0
			expected: []math.Monomial{math.NewMonomial(0, 1), math.NewMonomial(0)},
		},
		{
			name: "Duplicate cancellation (x0 + x1) * (x0 + x1) = x0 + x1 + x0x1 + x0x1 = x0 + x1",
			p1:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)},
			p2:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)},
			// x1*x1=x1, x1*x0=x0x1, x0*x1=x0x1, x0*x0=x0 → x0x1 cancels
			expected: []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure inputs are sorted as required
			sortPolynomial(tt.p1)
			sortPolynomial(tt.p2)

			got := math.MultiplyPolynomials(nil, tt.p1, tt.p2)

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

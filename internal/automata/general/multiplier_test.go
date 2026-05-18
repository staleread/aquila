package general_test

import (
	"slices"
	"testing"

	"github.com/staleread/aquila/internal/automata/general"
)

func sortPolynomial(p []general.Monomial) {
	slices.SortFunc(p, func(a, b general.Monomial) int {
		return general.CompareMonomials(b, a) // Descending
	})
}

func TestPolynomialMultiplier_Multiply(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   []general.Monomial
		expected []general.Monomial
	}{
		{
			name:     "Empty polynomials",
			p1:       []general.Monomial{},
			p2:       []general.Monomial{{1, 0, 0}},
			expected: []general.Monomial{},
		},
		{
			name:     "Single monomial multiplication",
			p1:       []general.Monomial{{1, 0, 0}}, // x0
			p2:       []general.Monomial{{2, 0, 0}}, // x1
			expected: []general.Monomial{{3, 0, 0}}, // x0*x1 = bit 0 | bit 1 = 3
		},
		{
			name: "Basic distribution (x0 + x1) * x2",
			p1:   []general.Monomial{{2, 0, 0}, {1, 0, 0}}, // x1, x0 (descending)
			p2:   []general.Monomial{{4, 0, 0}},            // x2
			// (x1+x0)*x2 = x1x2 + x0x2
			// x1x2 = 2|4 = 6
			// x0x2 = 1|4 = 5
			expected: []general.Monomial{{6, 0, 0}, {5, 0, 0}},
		},
		{
			name: "Cancellation x0 * (x0 + x1) = x0 + x0x1",
			p1:   []general.Monomial{{1, 0, 0}},
			p2:   []general.Monomial{{2, 0, 0}, {1, 0, 0}}, // x1, x0
			// x0*x1 = 3
			// x0*x0 = 1
			expected: []general.Monomial{{3, 0, 0}, {1, 0, 0}},
		},
		{
			name: "Duplicate cancellation (x0 + x1) * (x0 + x1) = x0 + x1 + x0x1 + x0x1 = x0 + x1",
			p1:   []general.Monomial{{2, 0, 0}, {1, 0, 0}},
			p2:   []general.Monomial{{2, 0, 0}, {1, 0, 0}},
			// x1*x1 = 2
			// x1*x0 = 3
			// x0*x1 = 3
			// x0*x0 = 1
			// Sorted descending before cancellation: 3, 3, 2, 1
			// After cancellation: 2, 1
			expected: []general.Monomial{{2, 0, 0}, {1, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &general.PolynomialMultiplier{}

			// Ensure inputs are sorted as required
			sortPolynomial(tt.p1)
			sortPolynomial(tt.p2)

			got := m.MultiplyPolynomials(nil, tt.p1, tt.p2)

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

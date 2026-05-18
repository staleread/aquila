package general_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/general"
)

func TestAddPolynomials(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   []general.Monomial
		expected []general.Monomial
	}{
		{
			name:     "Both empty",
			p1:       []general.Monomial{},
			p2:       []general.Monomial{},
			expected: []general.Monomial{},
		},
		{
			name:     "One empty",
			p1:       []general.Monomial{{4, 0, 0}, {1, 0, 0}},
			p2:       []general.Monomial{},
			expected: []general.Monomial{{4, 0, 0}, {1, 0, 0}},
		},
		{
			name:     "Disjoint sets",
			p1:       []general.Monomial{{8, 0, 0}, {2, 0, 0}},
			p2:       []general.Monomial{{4, 0, 0}, {1, 0, 0}},
			expected: []general.Monomial{{8, 0, 0}, {4, 0, 0}, {2, 0, 0}, {1, 0, 0}},
		},
		{
			name:     "Identical (full cancellation)",
			p1:       []general.Monomial{{4, 0, 0}, {1, 0, 0}},
			p2:       []general.Monomial{{4, 0, 0}, {1, 0, 0}},
			expected: []general.Monomial{},
		},
		{
			name:     "Partial overlap",
			p1:       []general.Monomial{{8, 0, 0}, {4, 0, 0}},
			p2:       []general.Monomial{{4, 0, 0}, {1, 0, 0}},
			expected: []general.Monomial{{8, 0, 0}, {1, 0, 0}}, // 4 cancels
		},
		{
			name:     "Interleaved",
			p1:       []general.Monomial{{16, 0, 0}, {4, 0, 0}, {1, 0, 0}},
			p2:       []general.Monomial{{8, 0, 0}, {4, 0, 0}, {2, 0, 0}},
			expected: []general.Monomial{{16, 0, 0}, {8, 0, 0}, {2, 0, 0}, {1, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := general.AddPolynomials(nil, tt.p1, tt.p2)

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

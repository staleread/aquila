package general_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata"
	"github.com/staleread/aquila/internal/automata/general"
)

func TestAddPolynomials(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   []automata.Word
		expected []automata.Word
	}{
		{
			name:     "Both empty",
			p1:       []automata.Word{},
			p2:       []automata.Word{},
			expected: []automata.Word{},
		},
		{
			name:     "One empty",
			p1:       []automata.Word{4, 0, 0, 1, 0, 0},
			p2:       []automata.Word{},
			expected: []automata.Word{4, 0, 0, 1, 0, 0},
		},
		{
			name:     "Disjoint sets",
			p1:       []automata.Word{8, 0, 0, 2, 0, 0},
			p2:       []automata.Word{4, 0, 0, 1, 0, 0},
			expected: []automata.Word{8, 0, 0, 4, 0, 0, 2, 0, 0, 1, 0, 0},
		},
		{
			name:     "Identical (full cancellation)",
			p1:       []automata.Word{4, 0, 0, 1, 0, 0},
			p2:       []automata.Word{4, 0, 0, 1, 0, 0},
			expected: []automata.Word{},
		},
		{
			name:     "Partial overlap",
			p1:       []automata.Word{8, 0, 0, 4, 0, 0},
			p2:       []automata.Word{4, 0, 0, 1, 0, 0},
			expected: []automata.Word{8, 0, 0, 1, 0, 0}, // 4 cancels
		},
		{
			name:     "Interleaved",
			p1:       []automata.Word{16, 0, 0, 4, 0, 0, 1, 0, 0},
			p2:       []automata.Word{8, 0, 0, 4, 0, 0, 2, 0, 0},
			expected: []automata.Word{16, 0, 0, 8, 0, 0, 2, 0, 0, 1, 0, 0},
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
					t.Errorf("At index %d: got %d, want %d", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

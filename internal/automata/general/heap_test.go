package general_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/general"
)

func TestMonomialMaxHeap(t *testing.T) {
	// m3 < m2 < m1
	m1 := general.Monomial{1, 0, 0}
	m2 := general.Monomial{0, 1, 0}
	m3 := general.Monomial{0, 0, 1}

	tests := []struct {
		name string
		ops  []struct {
			action string
			item   general.MonomialHeapItem
			want   general.Monomial
			len    int
		}
	}{
		{
			name: "Basic Operations",
			ops: []struct {
				action string
				item   general.MonomialHeapItem
				want   general.Monomial
				len    int
			}{
				{action: "push", item: general.MonomialHeapItem{Prod: m2}, len: 1},
				{action: "peek", want: m2, len: 1},
				{action: "push", item: general.MonomialHeapItem{Prod: m1}, len: 2},
				{action: "peek", want: m1, len: 2},
				{action: "push", item: general.MonomialHeapItem{Prod: m3}, len: 3},
				{action: "peek", want: m1, len: 3},
				{action: "pop", want: m1, len: 2},
				{action: "peek", want: m2, len: 2},
				{action: "pop", want: m2, len: 1},
				{action: "peek", want: m3, len: 1},
				{action: "pop", want: m3, len: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &general.MonomialMaxHeap{}

			for i, op := range tt.ops {
				switch op.action {
				case "push":
					h.Push(op.item)
				case "peek":
					got := h.Peek().Prod
					if !got.Equals(op.want) {
						t.Errorf("op %d (%s): Peek() = %v, want %v", i, op.action, got, op.want)
					}
				case "pop":
					got := h.Pop().Prod
					if !got.Equals(op.want) {
						t.Errorf("op %d (%s): Pop() = %v, want %v", i, op.action, got, op.want)
					}
				}
				if h.Len() != op.len {
					t.Errorf("op %d (%s): Len() = %d, want %d", i, op.action, h.Len(), op.len)
				}
			}
		})
	}
}

func TestMonomialMaxHeap_Property(t *testing.T) {
	h := &general.MonomialMaxHeap{}
	items := []general.Monomial{
		{10, 0, 0},
		{5, 0, 0},
		{15, 0, 0},
		{2, 0, 0},
		{20, 0, 0},
		{7, 0, 0},
	}

	for i, m := range items {
		h.Push(general.MonomialHeapItem{Prod: m, FirstIdx: i, SecondIdx: i})
	}

	if h.Len() != len(items) {
		t.Errorf("Expected heap length %d, got %d", len(items), h.Len())
	}

	var last general.Monomial
	isFirst := true

	for h.Len() > 0 {
		current := h.Pop().Prod

		if !isFirst && general.CompareMonomials(last, current) == -1 {
			t.Errorf("Heap property violated: %v came after %v", current, last)
		}
		last = current
		isFirst = false
	}
}

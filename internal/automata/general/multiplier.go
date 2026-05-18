package general

type PolynomialMultiplier struct {
	heap MonomialMaxHeap
}

// Multiplies two polynomials and appends the resulting monomials to the 'dst' buffer.
// Returns the updated slice header to handle standard Go append semantics.
func (m *PolynomialMultiplier) MultiplyPolynomials(dst, p1, p2 []Monomial) []Monomial {
	if len(p1) == 0 || len(p2) == 0 {
		return dst
	}

	m.heap.size = 0

	m.heap.Push(MonomialHeapItem{
		Prod: p1[0].Mul(p2[0]),
		I:    0,
		J:    0,
	})

	for m.heap.Len() > 0 {
		popped := m.heap.Pop()
		m.pushSuccessors(popped, p1, p2)

		duplicatesCnt := 1

		for m.heap.Len() > 0 && m.heap.Peek().Prod == popped.Prod {
			m.heap.Pop()
			duplicatesCnt++
		}

		if duplicatesCnt%2 != 0 {
			dst = append(dst, popped.Prod)
		}
	}
	return dst
}

func (m *PolynomialMultiplier) pushSuccessors(item MonomialHeapItem, p1, p2 []Monomial) {
	if item.J == 0 && item.I+1 < len(p1) {
		m.heap.Push(MonomialHeapItem{
			Prod: p1[item.I+1].Mul(p2[0]),
			I:    item.I + 1,
			J:    0,
		})
	}

	if item.J+1 < len(p2) {
		m.heap.Push(MonomialHeapItem{
			Prod: p1[item.I].Mul(p2[item.J+1]),
			I:    item.I,
			J:    item.J + 1,
		})
	}
}

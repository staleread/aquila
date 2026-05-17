package general

import "github.com/staleread/aquila/internal/automata"

type PolynomialMultiplier struct {
	heap MonomialMaxHeap
}

// Multiplies two polynomials and appends the resulting monomials to the 'dst' buffer.
// Returns the updated slice header to handle standard Go append semantics.
func (m *PolynomialMultiplier) MultiplyPolynomials(dst, p1, p2 []automata.Word) []automata.Word {
	if len(p1) == 0 || len(p2) == 0 {
		return dst
	}

	m.heap.size = 0

	lenM1 := len(p1) / automata.BlockWords
	lenM2 := len(p2) / automata.BlockWords

	p1m1 := NewMonomial(p1[0:automata.BlockWords])
	p2m1 := NewMonomial(p2[0:automata.BlockWords])

	m.heap.Push(MonomialHeapItem{
		Prod: p1m1.Mul(p2m1),
		I:    0,
		J:    0,
	})

	for m.heap.Len() > 0 {
		popped := m.heap.Pop()
		m.pushSuccessors(popped, lenM1, lenM2, p1, p2)

		duplicatesCnt := 1

		for m.heap.Len() > 0 && m.heap.Peek().Prod == popped.Prod {
			duplicate := m.heap.Pop()
			duplicatesCnt++

			m.pushSuccessors(duplicate, lenM1, lenM2, p1, p2)
		}

		if duplicatesCnt%2 != 0 {
			dst = append(dst, popped.Prod[:]...)
		}
	}
	return dst
}

func (m *PolynomialMultiplier) pushSuccessors(item MonomialHeapItem, lenM1, lenM2 int, p1, p2 []automata.Word) {
	if item.J == 0 && item.I+1 < lenM1 {
		nextIOffset := (item.I + 1) * automata.BlockWords
		p1Monom := NewMonomial(p1[nextIOffset : nextIOffset+automata.BlockWords])
		p2First := NewMonomial(p2[0:automata.BlockWords])

		m.heap.Push(MonomialHeapItem{
			Prod: p1Monom.Mul(p2First),
			I:    item.I + 1,
			J:    0,
		})
	}

	if item.J+1 < lenM2 {
		iOffset := item.I * automata.BlockWords
		nextJOffset := (item.J + 1) * automata.BlockWords

		p1Monom := NewMonomial(p1[iOffset : iOffset+automata.BlockWords])
		p2Monom := NewMonomial(p2[nextJOffset : nextJOffset+automata.BlockWords])

		m.heap.Push(MonomialHeapItem{
			Prod: p1Monom.Mul(p2Monom),
			I:    item.I,
			J:    item.J + 1,
		})
	}
}

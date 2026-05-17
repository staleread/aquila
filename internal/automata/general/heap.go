package general

type MonomialHeapItem struct {
	Prod Monomial
	I    int
	J    int
}

// Max heap for storing monomial products.
// Used for memory-efficient polynomial multiplication
type MonomialMaxHeap struct {
	data [MonomialHeapCap]MonomialHeapItem
	size int
}

func (h *MonomialMaxHeap) Len() int {
	return h.size
}

func (h *MonomialMaxHeap) Push(prod MonomialHeapItem) {
	h.data[h.size] = prod
	h.swim()
	h.size++
}

func (h *MonomialMaxHeap) Peek() MonomialHeapItem {
	return h.data[0]
}

func (h *MonomialMaxHeap) Pop() MonomialHeapItem {
	h.size--
	h.swap(0, h.size)
	h.sink(0)
	return h.data[h.size]
}

func (h *MonomialMaxHeap) more(i, j int) bool {
	return CompareMonomials(h.data[i].Prod, h.data[j].Prod) == 1
}

func (h *MonomialMaxHeap) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *MonomialMaxHeap) sink(parent int) {
	for {
		child1 := 2*parent + 1

		if child1 >= h.size {
			break
		}

		child := child1
		if child2 := child1 + 1; child2 < h.size && h.more(child2, child1) {
			child = child2
		}

		if !h.more(child, parent) {
			break
		}
		h.swap(child, parent)
		parent = child
	}
}

func (h *MonomialMaxHeap) swim() {
	child := h.size

	for {
		parent := (child - 1) / 2

		if child == parent || !h.more(child, parent) {
			break
		}
		h.swap(child, parent)
		child = parent
	}
}

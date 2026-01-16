package main

type idx uint16

func randIdx() idx {
	return idx(randByte()) | idx(randByte())<<8
}

func orderedIds(n int) []idx {
	ids := make([]idx, n)

	for i := range n {
		ids[i] = idx(i)
	}
	return ids
}

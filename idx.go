package main

type idx uint16

func randIdx() idx {
	return idx(randByte()) | idx(randByte())<<8
}

func readRandUniqueIds(ids []idx, maxIdx idx) {
	n := len(ids)
	added := make(map[idx]struct{})

	for len(added) < n {
		id := randIdx() % maxIdx

		if _, ok := added[id]; !ok {
			ids[len(added)] = id
			added[id] = struct{}{}
		}
	}
}

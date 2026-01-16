package main

type monom []idx

func randMonom(deg int, maxIdx idx) monom {
	m := make(monom, deg)
	added := make(map[idx]struct{})

	for len(added) < deg {
		id := randIdx() % maxIdx

		if _, ok := added[id]; !ok {
			m[len(added)] = id
			added[id] = struct{}{}
		}
	}
	return m
}

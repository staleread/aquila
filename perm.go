package main

type perm []int

func randPerm(n int) perm {
	p := make(perm, n-1)

	// Fisher–Yates shuffle
	for i := range n - 1 {
		p[i] = (int(randIdx())%n - i) + i
	}
	return p
}

func (v vec) permute(p perm) {
	for i := range len(v) - 1 {
		j := p[i]
		v[i], v[j] = v[j], v[i]
	}
}

func (v vec) permuteBack(p perm) {
	for i := len(v) - 2; i >= 0; i-- {
		j := p[i]
		v[j], v[i] = v[i], v[j]
	}
}

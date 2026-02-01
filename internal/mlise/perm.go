package mlise

import f "github.com/staleread/aquila/internal/gf2"

type perm []int

func randPerm(n int) perm {
	p := make(perm, n-1)
	rands := f.RandSyms(n - 1)

	// Fisher–Yates shuffle
	for i := range n - 1 {
		p[i] = int(rands[i])%(n-i) + i
	}
	return p
}

func (p perm) permute(v []f.Elt) {
	if len(v) != len(p)+1 {
		panic("Array size doesn't match the permutation size")
	}

	for i, j := range p {
		v[i], v[j] = v[j], v[i]
	}
}

func (p perm) permuteBack(v []f.Elt) {
	if len(v) != len(p)+1 {
		panic("Array size doesn't match the permutation size")
	}

	for i := len(p) - 1; i >= 0; i-- {
		j := p[i]
		v[i], v[j] = v[j], v[i]
	}
}

func (p perm) ids() []f.Sym {
	n := len(p) + 1
	ids := make([]f.Sym, n)

	for i := range f.Sym(n) {
		ids[i] = i
	}

	for i, j := range p {
		ids[i], ids[j] = ids[j], ids[i]
	}
	return ids
}

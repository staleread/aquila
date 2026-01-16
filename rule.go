package main

type ruleFold struct {
	lin  *sle
	nlin *snle
}

type rule struct {
	size  int
	p     perm
	folds []ruleFold
}

func randRule(size, folds, deg int) *rule {
	p := randPerm(size)
	n := size / folds
	rFolds := make([]ruleFold, folds)

	rFolds[0] = ruleFold{
		lin:  randSle(n),
		nlin: emptySnle(n),
	}

	for i := 1; i < folds; i++ {
		maxIdx := idx(n * i)

		rFolds[i] = ruleFold{
			lin:  randSle(n),
			nlin: randSnle(n, deg, maxIdx),
		}
	}
	return &rule{size, p, rFolds}
}

func (r *rule) decrypt(pt, ct vec) {
	n := r.size / len(r.folds)

	for i, f := range r.folds {
		noise := zeros(n)
		xPrev := pt[:n*i]
		bCurr := ct[n*i : n*i+n]

		f.nlin.eval(xPrev, noise)
		bCurr.sub(noise)

		xCurr := pt[n*i : n*i+n]

		f.lin.solve(xCurr, bCurr)
	}
	permuteBack(pt, r.p)
}

func (r *rule) encrypt(pt, ct vec) {
	n := r.size / len(r.folds)

	permute(pt, r.p)

	for i, f := range r.folds {
		xCurr := pt[n*i : n*i+n]
		bCurr := ct[n*i : n*i+n]

		f.lin.eval(xCurr, bCurr)

		noise := zeros(n)
		xPrev := pt[:n*i]

		f.nlin.eval(xPrev, noise)
		bCurr.add(noise)
	}
}

func (r *rule) toSnle() *snle {
	n := r.size / len(r.folds)
	se := emptySnle(r.size)

	ids := orderedIds(r.size)
	permute(ids, r.p)

	for i, f := range r.folds {
		lin := f.lin.coefs()
		nlin := f.nlin.data

		for j := range n {
			p := poly{}

			// Non linear part
			for _, m := range nlin[j] {
				if len(m) == 0 {
					continue
				}
				mp := make(monom, len(m))

				for l, id := range m {
					mp[l] = ids[id]
				}
				p = append(p, mp)
			}

			// Linear part
			for k := range n {
				val := lin.at(j, k)

				if val != 0 {
					coef := ids[n*i+k]
					p = append(p, monom{coef})
				}
			}
			se.data[n*i+j] = p
		}
	}
	return se
}

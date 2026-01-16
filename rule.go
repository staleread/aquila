package main

type ruleFold struct {
	lin   *sle
	noise *snle
}

type rule struct {
	n     int
	p     perm
	folds []ruleFold
}

func randRule(s, n, deg int) *rule {
	p := randPerm(s * n)
	folds := make([]ruleFold, s)

	folds[0] = ruleFold{
		lin:   randSle(n),
		noise: emptySnle(n),
	}

	for i := 1; i < s; i++ {
		maxIdx := idx(n * i)

		folds[i] = ruleFold{
			lin:   randSle(n),
			noise: randSnle(n, deg, maxIdx),
		}
	}
	return &rule{n: n * s, p: p, folds: folds}
}

func (r *rule) decrypt(pt, ct vec) {
	s := len(r.folds)
	n := r.n / s

	permute(ct, r.p)

	for i, f := range r.folds {
		noise := zeros(n)
		xPrev := pt[:n*i]
		bCurr := pt[n*i : n*i+n]

		f.noise.eval(xPrev, noise)
		bCurr.sub(noise)

		xCurr := ct[n*i : n*i+n]

		f.lin.solve(xCurr, bCurr)
	}
}

func (r *rule) encrypt(pt, ct vec) {
	n := r.n / len(r.folds)

	for i, f := range r.folds {
		xCurr := pt[n*i : n*i+n]
		bCurr := ct[n*i : n*i+n]

		f.lin.eval(xCurr, bCurr)

		noise := zeros(n)
		xPrev := pt[:n*i]

		f.noise.eval(xPrev, noise)
		bCurr.add(noise)
	}

	permuteBack(ct, r.p)
}

func (r *rule) toSnle() *snle {
	n := r.n / len(r.folds)
	res := emptySnle(r.n)

	ids := orderedIds(r.n)
	permute(ids, r.p)

	for i, f := range r.folds {
		lin := f.lin.coefs()
		nlin := f.noise.data

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

			res.data[n*i+j] = p
		}
	}
	return res
}

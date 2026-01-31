package la

import f "github.com/staleread/aquila/internal/gf2"

type Vec []f.Elt

func ZeroVec(n int) Vec {
	return make(Vec, n)
}

func RandVec(n int) Vec {
	return Vec(f.RandEls(n))
}

func (a Vec) Add(b Vec) {
	for i := range len(a) {
		a[i] = f.Add(a[i], b[i])
	}
}

func (a Vec) Sub(b Vec) {
	for i := range len(a) {
		a[i] = f.Sub(a[i], b[i])
	}
}

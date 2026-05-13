package gf2

type Vector []Element

func (v Vector) FillZero() {
	for i := range len(v) {
		v[i] = Zero
	}
}

func (v Vector) FillOnes() {
	for i := range len(v) {
		v[i] = One
	}
}

func (v Vector) FillRand() {
	FillRand(v)
}

func (v Vector) Add(other Vector) {
	for i := range len(v) {
		v[i] ^= other[i]
	}
}

func (v Vector) Sub(other Vector) {
	for i := range len(v) {
		v[i] ^= other[i]
	}
}

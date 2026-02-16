package la

import f "github.com/staleread/aquila/internal/gf2"

type Vector []f.Element

func ZeroVector(n int) Vector {
	return make(Vector, n)
}

func RandVector(n int) Vector {
	return Vector(f.RandElements(n))
}

func (a Vector) Add(b Vector) {
	for i := range len(a) {
		a[i] = f.Add(a[i], b[i])
	}
}

func (a Vector) Sub(b Vector) {
	for i := range len(a) {
		a[i] = f.Sub(a[i], b[i])
	}
}

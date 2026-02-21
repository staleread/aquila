package linalg

import "github.com/staleread/aquila/internal/field"

type Vector []field.Element

func ZeroVector(n int) Vector {
	return make(Vector, n)
}

func RandVector(n int) Vector {
	return Vector(field.RandElements(n))
}

func (a Vector) Add(b Vector) {
	for i := range len(a) {
		a[i] = field.Add(a[i], b[i])
	}
}

func (a Vector) Sub(b Vector) {
	for i := range len(a) {
		a[i] = field.Sub(a[i], b[i])
	}
}

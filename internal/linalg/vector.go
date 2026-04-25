package linalg

import "github.com/staleread/aquila/internal/field"

type Vector []field.Element

func ZeroVector(n int) Vector {
	return make(Vector, n)
}

func RandVector(n int) Vector {
	return Vector(field.RandElements(n))
}

func (self Vector) Add(other Vector) {
	for i := range len(self) {
		self[i] = field.Add(self[i], other[i])
	}
}

func (self Vector) Sub(other Vector) {
	for i := range len(self) {
		self[i] = field.Sub(self[i], other[i])
	}
}

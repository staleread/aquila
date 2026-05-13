package gf2

import "fmt"

type SquareMatrix struct {
	n    int
	data Vector
}

func NewSquareMatrix(n int, data Vector) (*SquareMatrix, error) {
	if len(data) != n*n {
		return nil, fmt.Errorf("matrix data size %d does not match dimensions %dx%d", len(data), n, n)
	}

	return &SquareMatrix{
		n:    n,
		data: data,
	}, nil
}

func (m *SquareMatrix) At(i, j int) Element {
	return m.data[m.n*i+j]
}

func (m *SquareMatrix) Set(i, j int, val Element) {
	m.data[m.n*i+j] = val
}

func (m *SquareMatrix) SubstituteForward(x, b Vector) {
	for i := range m.n {
		num := b[i]

		for j := range i {
			num ^= m.At(i, j) & x[j]
		}
		x[i] = num
	}
}

func (m *SquareMatrix) SubstituteBackward(x, b Vector) {
	n := m.n

	for i := n - 1; i >= 0; i-- {
		num := b[i]

		for j := i + 1; j < n; j++ {
			num ^= m.At(i, j) & x[j]
		}
		x[i] = num
	}
}

func (m *SquareMatrix) MultiplyLower(dst, src Vector) {
	for i := range m.n {
		sum := src[i]

		for j := range i {
			sum ^= m.At(i, j) & src[j]
		}
		dst[i] = sum
	}
}

func (m *SquareMatrix) MultiplyUpper(dst, src Vector) {
	n := m.n
	for i := range n {
		sum := src[i]

		for j := i + 1; j < n; j++ {
			sum ^= m.At(i, j) & src[j]
		}
		dst[i] = sum
	}
}

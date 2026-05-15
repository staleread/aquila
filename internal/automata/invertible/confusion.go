package invertible

import (
	"io"
	"slices"
)

type Subscript = uint8

type ConfusionMap struct {
	data []Subscript
}

func NewConfusionMap(arena []byte) *ConfusionMap {
	return &ConfusionMap{data: arena}
}

func (m *ConfusionMap) Generate(rnd io.Reader, foldIdx int, perm Permutation) error {
	if _, err := io.ReadFull(rnd, m.data); err != nil {
		return err
	}

	maxSub := Subscript(foldIdx * VectorSize)
	for i := range m.data {
		m.data[i] &= maxSub
	}

	subIdx := 0
	for range VectorSize {
		for j := ConfusionDegree; j > 0; j-- {
			upperBound := maxSub - Subscript(j)

			for k := range j {
				candidateIdx := m.data[subIdx] % (upperBound + Subscript(k) + 1)
				candidate := perm[candidateIdx]

				duplicate := slices.Contains(m.data[subIdx-k:subIdx], candidate)

				if duplicate {
					m.data[subIdx] = perm[upperBound+Subscript(k)]
				} else {
					m.data[subIdx] = candidate
				}
				subIdx++
			}
		}
	}
	return nil
}

func (m *ConfusionMap) Eval(state []uint64) Vector {
	var res Vector
	subIdx := 0

	for i := range VectorSize {
		var sum Vector

		for j := range ConfusionDegree {
			prod := Vector(1)
			subCnt := ConfusionDegree - j

			for range subCnt {
				idx := uint8(m.data[subIdx])

				bit := (state[idx/64] >> (idx % 64)) & 1
				prod &= Vector(bit)

				subIdx++
			}
			sum ^= prod
		}
		res |= sum << i
	}
	return res
}

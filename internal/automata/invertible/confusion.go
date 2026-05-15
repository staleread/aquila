package invertible

import "slices"

type ConfusionMap struct {
	data []Subscript
}

func NewConfusionMap(arena []Subscript) *ConfusionMap {
	return &ConfusionMap{
		data: arena,
	}
}

func (m *ConfusionMap) FillRand(foldIdx int, perm Permutation) {
	maxSub := Subscript(foldIdx * VectorSize)
	FillRandSubscripts(m.data, maxSub)

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

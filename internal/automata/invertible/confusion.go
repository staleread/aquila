package invertible

import "slices"

type ConfusionMap []Subscript

func (p ConfusionMap) FillRand(foldIdx int, perm Permutation) {
	maxSub := Subscript(foldIdx * VectorSize)
	FillRandSubscripts(p, maxSub)

	subIdx := 0
	for range VectorSize {
		for j := ConfusionDegree; j > 0; j-- {
			upperBound := maxSub - Subscript(j)

			for k := range j {
				candidateIdx := p[subIdx] % (upperBound + Subscript(k) + 1)
				candidate := perm[candidateIdx]

				duplicate := slices.Contains(p[subIdx-k:subIdx], candidate)

				if duplicate {
					p[subIdx] = perm[upperBound+Subscript(k)]
				} else {
					p[subIdx] = candidate
				}
				subIdx++
			}
		}
	}
}

func (p ConfusionMap) Eval(state []uint64) Vector {
	var res Vector
	subIdx := 0

	for i := range VectorSize {
		var sum Vector

		for j := range ConfusionDegree {
			prod := Vector(1)
			subCnt := ConfusionDegree - j

			for range subCnt {
				idx := uint8(p[subIdx])

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

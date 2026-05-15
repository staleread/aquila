package invertible

import (
	"github.com/staleread/aquila/internal/automata"
	"slices"
)

type ConfusionMap []Subscript

func (p ConfusionMap) FillRand(maxSub Subscript) {
	FillRandSubscripts(p, maxSub)

	subIdx := 0
	for range Dim {
		for j := automata.ConfusionMapDegree; j > 0; j-- {
			upperBound := maxSub - Subscript(j)

			for k := range j {
				candidate := p[subIdx] % (upperBound + Subscript(k) + 1)
				duplicate := slices.Contains(p[subIdx-k:subIdx], candidate)

				if duplicate {
					p[subIdx] = upperBound + Subscript(k)
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

	for i := range Dim {
		var sum Vector

		for j := range automata.ConfusionMapDegree {
			prod := Vector(1)
			subCnt := automata.ConfusionMapDegree - j

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

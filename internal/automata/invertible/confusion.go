package invertible

import (
	"io"
	"slices"

	"github.com/staleread/aquila/internal/automata"
)

type Subscript = uint8

type ConfusionMap struct {
	data []Subscript
}

func NewConfusionMap(arena []byte) *ConfusionMap {
	return &ConfusionMap{data: arena}
}

func (m *ConfusionMap) Generate(rnd io.Reader, maxSub Subscript, perm *Permutation) error {
	if _, err := io.ReadFull(rnd, m.data); err != nil {
		return err
	}

	for i := range m.data {
		m.data[i] %= maxSub
	}

	subIdx := 0
	for range VectorSize {
		for j := ConfusionDegree; j > 1; j-- {
			upperBound := int(maxSub) - j

			for k := range j {
				candidateIdx := int(m.data[subIdx]) % (upperBound + k + 1)
				candidate := perm.data[candidateIdx]

				duplicate := slices.Contains(m.data[subIdx-k:subIdx], candidate)

				if duplicate {
					m.data[subIdx] = perm.data[upperBound+k]
				} else {
					m.data[subIdx] = candidate
				}
				subIdx++
			}
		}
	}
	return nil
}

func (m *ConfusionMap) Eval(state *automata.Block) Vector {
	if m == nil {
		return Vector(0)
	}

	var res Vector
	subIdx := 0

	for i := range VectorSize {
		var sum Vector

		for j := range ConfusionDegree - 1 {
			prod := Vector(1)
			subCnt := ConfusionDegree - j

			for range subCnt {
				bit := state.At(int(m.data[subIdx]))
				prod &= Vector(bit)

				subIdx++
			}
			sum ^= prod
		}
		res |= sum << i
	}
	return res
}

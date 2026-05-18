package math

import (
	"io"
	"slices"

	"github.com/staleread/aquila/internal/automata/core"
)

type Subscript = uint8

type ConfusionMap struct {
	Data []Subscript
}

func NewConfusionMap(arena []byte) *ConfusionMap {
	return &ConfusionMap{Data: arena}
}

func (m *ConfusionMap) Generate(rnd io.Reader, maxSub Subscript, perm *Permutation) error {
	if _, err := io.ReadFull(rnd, m.Data); err != nil {
		return err
	}

	for i := range m.Data {
		m.Data[i] %= maxSub
	}

	subIdx := 0
	for range VectorSize {
		for j := ConfusionDegree; j > 1; j-- {
			upperBound := int(maxSub) - j

			for k := range j {
				candidateIdx := int(m.Data[subIdx]) % (upperBound + k + 1)
				candidate := perm.Data[candidateIdx]

				duplicate := slices.Contains(m.Data[subIdx-k:subIdx], candidate)

				if duplicate {
					m.Data[subIdx] = perm.Data[upperBound+k]
				} else {
					m.Data[subIdx] = candidate
				}
				subIdx++
			}
		}
	}
	return nil
}

func (m *ConfusionMap) Eval(state *core.Block) Vector {
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
				bit := state.At(int(m.Data[subIdx]))
				prod &= Vector(bit)

				subIdx++
			}
			sum ^= prod
		}
		res |= sum << i
	}
	return res
}

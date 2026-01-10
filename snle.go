package main

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type monom []VarIdx
type poly []monom

type SFoldSNLE struct {
	s    Size
	data []poly
}

func RandSFoldSNLE(n, s, deg Size) *SFoldSNLE {
	idxPerPoly := nSum(deg)
	pCnt := n * s

	randIds := randVars(idxPerPoly * (pCnt - n)) // No polynomials for the first fold
	data := make([]poly, pCnt)

	for i := n; i < pCnt; i++ {
		p := poly{}

		for j := range deg {
			m := monom{}
			maxId := VarIdx(i / n * n) // Restrict the max idx to ones from the previous folds
			jSum := nSum(j)

			for k := range j + 1 {
				id := randIds[idxPerPoly*(i-n)+jSum+k] % maxId
				m = append(m, id)
			}
			p = append(p, m)
		}
		data[i] = p
	}
	return &SFoldSNLE{s, data}
}

func (snle *SFoldSNLE) Eval(x, b *BVec) {
	s := snle.s
	data := snle.data

	for i := range Size(len(data)) {
		bRow := i / s
		bBit := s - (i % s) - 1

		var sum Batch = 0

		for _, m := range data[i] {
			var factor Batch = 1

			for _, id := range m {
				xRow := id / s
				xBit := s - (id % s) - 1

				factor &= x.data[xRow] >> xBit
			}
			sum ^= factor
		}
		b.data[bRow] |= sum << bBit
	}
}

func (snle *SFoldSNLE) String() string {
	data := snle.data
	sb := strings.Builder{}

	for i := range Size(len(data)) {
		sb.WriteString(fmt.Sprintf("y%d = ", i))

		for j, m := range data[i] {
			if j > 0 {
				sb.WriteString(" + ")
			}
			for k, id := range m {
				if k > 0 {
					sb.WriteRune('*')
				}
				sb.WriteString(fmt.Sprintf("x%d", id))
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func randVars(n Size) []VarIdx {
	buffSize := n * VarIdxBytes
	buff := make([]byte, buffSize)
	rand.Read(buff)

	vars := make([]VarIdx, n)

	for buffIdx := range buffSize {
		varIdx := buffIdx / VarIdxBytes

		shift := (buffIdx % VarIdxBytes) * 8
		vars[varIdx] |= VarIdx(buff[buffIdx]) << shift

		if shift == 0 {
			vars[varIdx] &= MaxVarIdx
		}
	}
	return vars
}

func nSum(n Size) Size {
	return n * (n + 1) / 2
}

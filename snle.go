package main

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type monom []VarIdx
type poly []monom

type SNLE struct {
	data []poly
}

func EmptySNLE(n Size) *SNLE {
	return &SNLE{data: make([]poly, n)}
}

func RandSNLE(n, deg Size, maxVarIdx VarIdx) *SNLE {
	idxPerPoly := nSum(deg)

	ids := randVarIds(idxPerPoly * n)
	data := make([]poly, n)

	for i := range n {
		var jSum Size = 0
		p := poly{}

		for j := range deg {
			jSum += j
			m := monom{}

			for k := range j + 1 {
				idx := ids[idxPerPoly*i+jSum+k] % maxVarIdx
				m = append(m, idx)
			}
			p = append(p, m)
		}
		data[i] = p
	}
	return &SNLE{data}
}

func (snle *SNLE) Eval(x, b Vec) {
	for i, p := range snle.data {
		var sum Elem = 0

		for _, m := range p {
			var factor Elem = 1

			for _, idx := range m {
				factor &= x[idx]
			}
			sum ^= factor
		}
		// GF(2) addition in case b is not empty
		b[i] ^= sum
	}
}

func (snle *SNLE) String() string {
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

func randVarIds(n Size) []VarIdx {
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

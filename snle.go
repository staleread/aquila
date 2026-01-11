package main

import "crypto/rand"

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
		// Use GF(2) addition(=subraction) instead of overriding b in case b is
		// a non-zero vector.
		b[i] ^= sum
	}
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

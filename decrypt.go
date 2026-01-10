package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

type fold struct {
	sle  *SLE
	snle *SNLE
}

type DecRule struct {
	perm  []VarIdx
	folds []fold
}

func RandDecRule(s, n, deg Size) *DecRule {
	perm := varPermutation(n * s)
	folds := make([]fold, s)

	folds[0] = fold{
		sle:  RandSLE(n),
		snle: EmptySNLE(n),
	}

	for i := Size(1); i < s; i++ {
		folds[i] = fold{
			sle:  RandSLE(n),
			snle: RandSNLE(n, deg, n*i),
		}
	}
	return &DecRule{perm, folds}
}

func (rule *DecRule) Apply(ct Vec) {
	s := len(rule.folds)
	n := len(rule.perm) / s

	ct.Permute(rule.perm)

	for i, f := range rule.folds {
		xPrev := ct[:n*i]
		xCurr := ct[n*i : n*i+n]

		f.snle.Eval(xPrev, xCurr)
		f.sle.Solve(xCurr, xCurr)
	}
}

func (rule *DecRule) String() string {
	s := len(rule.folds)
	n := len(rule.perm) / s
	idxPad := Size(math.Log10(float64(n*s+1)) + 1)

	sb := strings.Builder{}

	for i, f := range rule.folds {
		sle := f.sle.Coefs().data
		snle := f.snle.data

		for j := range n {
			permJ := rule.perm[n*i+j]
			sb.WriteString(fmt.Sprintf("y%-*d = ", idxPad, permJ+1))

			// Linear part
			for k := range n {
				if k > 0 {
					sb.WriteString(" + ")
				}
				var varStr string

				if sle[n*j+k] != 0 {
					varStr = fmt.Sprintf("x%d", n*i+k+1)
				}
				sb.WriteString(fmt.Sprintf("%-*s", idxPad+2, varStr))
			}

			// Non-linear part
			for _, m := range snle[j] {
				sb.WriteString(" + ")

				for l, id := range m {
					if l > 0 {
						sb.WriteRune('*')
					}
					sb.WriteString(fmt.Sprintf("x%d", id+1))
				}
			}
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func varPermutation(n Size) []VarIdx {
	perm := make([]VarIdx, n)

	for i := range n {
		perm[i] = i
	}

	rand.Shuffle(len(perm), func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})

	return perm
}

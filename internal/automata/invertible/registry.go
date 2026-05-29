package invertible

import (
	"slices"

	"github.com/staleread/aquila/internal/automata/core"
	"github.com/staleread/aquila/internal/automata/math"
)

type SymbolicRegistry struct {
	arena    []math.Monomial
	offsets  [RulesCount][core.BlockSize + 1]uint32
	invPerms [RulesCount][core.BlockSize]uint8
}

func (e *SymbolicRegistry) GetPolynomial(ruleIdx, bitIdx int) []math.Monomial {
	k := e.invPerms[ruleIdx][bitIdx]
	start := e.offsets[ruleIdx][k]
	end := e.offsets[ruleIdx][k+1]
	return e.arena[start:end]
}

func CompileRegistry(ca *CA) *SymbolicRegistry {
	e := &SymbolicRegistry{}
	e.arena = make([]math.Monomial, 0, MaxCAMonomials)
	monomials := make([]math.Monomial, 0, SymbolicPolynomialSize)

	for r := range RulesCount {
		rule := ca.getRule(r)
		perm := rule.getPermutation()

		var sleCoefs [math.VectorSize]math.Vector

		for k := range core.BlockSize {
			f := k / math.VectorSize
			i := k % math.VectorSize
			b := perm.Data[k]

			e.invPerms[r][b] = uint8(k)
			e.offsets[r][k] = uint32(len(e.arena))

			fold := rule.getFold(f)

			fold.sle.Coefs(sleCoefs[:])
			row := sleCoefs[i]

			monomials = monomials[:0]

			// SLE part
			foldOffset := f * math.VectorSize
			for j := range math.VectorSize {
				if (row>>j)&1 == 1 {
					var m math.Monomial
					m.SetAt(perm.Data[foldOffset+j], 1)
					monomials = append(monomials, m)
				}
			}

			// Confusion part
			if f > 0 && len(fold.confusion.Data) > 0 {
				const SubsPerBit = math.ConfusionMapBytes / math.VectorSize
				cursor := i * SubsPerBit

				for j := range math.ConfusionDegree - 1 {
					degree := math.ConfusionDegree - j
					var m math.Monomial
					for range degree {
						m.SetAt(fold.confusion.Data[cursor], 1)
						cursor++
					}
					monomials = append(monomials, m)
				}
			}

			slices.SortFunc(monomials, func(a, b math.Monomial) int {
				return math.CompareMonomials(b, a)
			})

			e.arena = append(e.arena, monomials...)
		}
		e.offsets[r][core.BlockSize] = uint32(len(e.arena))
	}

	return e
}

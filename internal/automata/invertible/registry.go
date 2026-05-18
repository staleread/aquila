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

			monomials := make([]math.Monomial, 0, SymbolicPolynomialSize)

			// SLE part
			foldOffset := f * math.VectorSize
			for j := range math.VectorSize {
				if (row>>j)&1 == 1 {
					m := math.Monomial{}
					sub := perm.Data[foldOffset+j]
					m[sub/32] |= (1 << (sub % 32))
					monomials = append(monomials, m)
				}
			}

			// Confusion part
			if f > 0 && fold.confusion != nil {
				const SubsPerBit = 5
				subIdx := i * SubsPerBit

				// Monomial degree 3
				m3 := math.Monomial{}
				for k := range 3 {
					sub := fold.confusion.Data[subIdx+k]
					m3[sub/32] |= (1 << (sub % 32))
				}
				monomials = append(monomials, m3)

				// Monomial degree 2
				m2 := math.Monomial{}
				for k := 3; k < 5; k++ {
					sub := fold.confusion.Data[subIdx+k]
					m2[sub/32] |= (1 << (sub % 32))
				}
				monomials = append(monomials, m2)
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

package sparse

import (
	"fmt"
	"github.com/staleread/aquila/internal/poly"
	"github.com/staleread/aquila/internal/poly/compiled"
	"strings"
)

type Polyset []Polynomial

func (set Polyset) Compile() compiled.Polyset {
	cmp := compiled.Polyset{
		Subscripts:   make([]poly.Subscript, 0),
		PolyCount:    len(set),
		PolyOffsets:  make([]int, 1),
		MonomOffsets: make([]int, 1),
	}

	var monomOffset int
	var polyOffset int

	for _, poly := range set {
		for monom := range poly.Monomials() {
			for sub := range monom.Subscripts() {
				cmp.Subscripts = append(cmp.Subscripts, sub)
				monomOffset++
			}
			cmp.MonomOffsets = append(cmp.MonomOffsets, monomOffset)
			polyOffset++
		}
		cmp.PolyOffsets = append(cmp.PolyOffsets, polyOffset)
	}
	return cmp
}

func (set Polyset) String() string {
	sb := strings.Builder{}

	for i, poly := range set {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		firstMonomial := true
		for monom := range poly.Monomials() {
			if !firstMonomial {
				sb.WriteString(" + ")
			}
			firstMonomial = false

			firstSubscript := true
			for sub := range monom.Subscripts() {
				if !firstSubscript {
					sb.WriteRune('*')
				}
				firstSubscript = false

				fmt.Fprintf(&sb, "x%d", sub+1)
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

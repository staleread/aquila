package sparse

import (
	"fmt"
	"github.com/staleread/aquila/internal/poly"
	"github.com/staleread/aquila/internal/poly/compiled"
	"strings"
)

type Polyset struct {
	polynoms  []Polynomial
	monomPool *monomialInternPool
}

func (set *Polyset) ComposeWith(other *Polyset) {
	// prodCache := make(map[*Monomial]Polynomial)
	cmpPool := newMonomialInternPool()

	for i, p := range set.polynoms {
		sum := Polynomial{}

		for m := range p.Monomials() {
			// if cached, hit := prodCache[m]; hit {
			// 	sum.AddTo(cached)
			// 	continue
			// }

			prod := Polynomial{}
			isFirst := true

			for s := range m.Subscripts() {
				if isFirst {
					prod.AddTo(other.polynoms[s])
					isFirst = false
					continue
				}
				cmpPool.mulPolynomialBy(prod, other.polynoms[s])
			}
			sum.AddTo(prod)
			// prodCache[m] = prod
		}
		set.polynoms[i] = sum
	}
}

func (set *Polyset) Compile() *compiled.Polyset {
	cmp := &compiled.Polyset{
		Subscripts:   make([]poly.Subscript, 0),
		PolyCount:    len(set.polynoms),
		PolyOffsets:  make([]int, 1),
		MonomOffsets: make([]int, 1),
	}

	var monomOffset int
	var polyOffset int

	for _, poly := range set.polynoms {
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

func (set *Polyset) String() string {
	sb := strings.Builder{}

	for i, poly := range set.polynoms {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		isFirstMonomial := true
		for monom := range poly.Monomials() {
			if !isFirstMonomial {
				sb.WriteString(" + ")
			}
			isFirstMonomial = false

			isFirstSubscript := true
			for sub := range monom.Subscripts() {
				if !isFirstSubscript {
					sb.WriteRune('*')
				}
				isFirstSubscript = false

				fmt.Fprintf(&sb, "x%d", sub+1)
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

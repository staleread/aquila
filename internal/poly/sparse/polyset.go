package sparse

import (
	"github.com/staleread/aquila/internal/poly"
	"github.com/staleread/aquila/internal/poly/compiled"
	"strings"
)

type Polyset struct {
	data []Polynomial
}

func NewPolyset(polys []Polynomial) Polyset {
	return Polyset{data: polys}
}

func (self Polyset) Compose(other Polyset) {
	for i, poly := range self.data {
		sum := NewPolynomial()

		for monom := range poly.Monomials() {
			prod := NewPolynomial()
			isFirst := true

			for sub := range monom.Subscripts() {
				if isFirst {
					prod = other.data[sub].Clone()
					isFirst = false
				} else {
					prod.Mul(other.data[sub])
				}
			}
			sum.Add(prod)
		}
		self.data[i] = sum
	}
}

func (self Polyset) Compile() *compiled.Polyset {
	cmp := &compiled.Polyset{
		Subscripts:   make([]poly.Subscript, 0),
		PolyCount:    len(self.data),
		PolyOffsets:  make([]int, 1),
		MonomOffsets: make([]int, 1),
	}

	var monomOffset int
	var polyOffset int

	for _, poly := range self.data {
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

func (self Polyset) Write(sb *strings.Builder) {
	for i, poly := range self.data {
		poly.Write(sb, i+1)
		sb.WriteRune('\n')
	}
}

func (self Polyset) String() string {
	sb := &strings.Builder{}
	self.Write(sb)

	return sb.String()
}

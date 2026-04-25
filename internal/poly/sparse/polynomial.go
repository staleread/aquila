package sparse

import (
	"fmt"
	"iter"
	"maps"
	"strings"
)

type Polynomial struct {
	data map[uint64]Monomial
}

func NewPolynomial() Polynomial {
	return Polynomial{data: make(map[uint64]Monomial)}
}

func (self Polynomial) Monomials() iter.Seq[Monomial] {
	return maps.Values(self.data)
}

func (self Polynomial) Toggle(monom Monomial) {
	hash := monom.Hash()
	found, hit := self.data[hash]

	if !hit {
		self.data[hash] = monom
	} else if monom.Equals(found) {
		delete(self.data, hash)
	} else {
		panic("monomial hash collision detected")
	}
}

func (self Polynomial) Clone() Polynomial {
	clone := NewPolynomial()
	maps.Copy(clone.data, self.data)
	return clone
}

func (self Polynomial) Add(other Polynomial) {
	for _, monom := range other.data {
		self.Toggle(monom)
	}
}

func (self Polynomial) Mul(other Polynomial) {
	result := NewPolynomial()
	for _, sMonom := range self.data {
		for _, oMonom := range other.data {
			result.Toggle(sMonom.Product(oMonom))
		}
	}
	clear(self.data)
	maps.Copy(self.data, result.data)
}

func (self Polynomial) Write(sb *strings.Builder, idx int) {
	isFirst := true

	if idx > 0 {
		fmt.Fprintf(sb, "y%d = ", idx)
	} else {
		fmt.Fprint(sb, "y = ")
	}

	for monom := range self.Monomials() {
		if !isFirst {
			sb.WriteString(" + ")
		} else {
			isFirst = false
		}
		monom.Write(sb)
	}
}

func (self Polynomial) String() string {
	sb := &strings.Builder{}
	self.Write(sb, 0)

	return sb.String()
}

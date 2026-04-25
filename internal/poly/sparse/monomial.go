package sparse

import (
	"fmt"
	"iter"
	"strings"

	"github.com/staleread/aquila/internal/poly"
)

const monomChunkBits = 64

type Monomial struct {
	hash uint64
	data []uint64
}

func EmptyMonomial(words int) Monomial {
	return Monomial{0, make([]uint64, words)}
}

func NewMonomial(words int, subs ...poly.Subscript) Monomial {
	monom := EmptyMonomial(words)

	for _, sub := range subs {
		monom.Toggle(sub)
	}
	return monom
}

func hashByte(s uint8) uint64 {
	x := uint64(s) + 1
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}

func (self Monomial) Len() int {
	return len(self.data) * monomChunkBits
}

func (self Monomial) Has(sub poly.Subscript) bool {
	chunk := self.data[len(self.data)-int(sub)/monomChunkBits-1]
	bitValue := chunk >> (sub % monomChunkBits) & 1

	return bitValue == 1
}

func (self Monomial) Toggle(pos poly.Subscript) {
	chunkIdx := len(self.data) - int(pos)/monomChunkBits - 1
	tmp := uint64(1) << (pos % monomChunkBits)

	self.data[chunkIdx] ^= tmp
}

func (self Monomial) Equals(other Monomial) bool {
	n := len(self.data)

	if n != len(other.data) {
		return false
	}

	for i := range n {
		if self.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

func (self Monomial) Hash() uint64 {
	const chunkBytes = monomChunkBits / 8

	if self.hash != 0 {
		return self.hash
	}

	n := len(self.data)
	var hash uint64

	for i := range n {
		tmp := self.data[n-i-1]

		for j := range chunkBytes {
			byteToHash := uint8((tmp >> (j * 8)) & 0xff)
			hash = hash*31 + hashByte(byteToHash)
		}
	}

	self.hash = hash
	return hash
}

func (self Monomial) Mul(other Monomial) {
	n := len(self.data)

	if n != len(other.data) {
		panic("Size of monomials do not match")
	}

	for i := range n {
		self.data[i] |= other.data[i]
	}
}

func (self Monomial) Product(other Monomial) Monomial {
	n := len(self.data)

	if n != len(other.data) {
		panic("Size of monomials do not match")
	}

	monom := EmptyMonomial(n)

	for i := range n {
		monom.data[i] = self.data[i] | other.data[i]
	}
	return monom
}

func (self Monomial) Clone() Monomial {
	n := len(self.data)

	monom := EmptyMonomial(n)

	for i := range n {
		monom.data[i] = self.data[i]
	}
	return monom
}

func (self Monomial) Write(sb *strings.Builder) {
	n := len(self.data)
	isFirst := true

	for i := range n {
		tmp := self.data[n-i-1]

		for j := range monomChunkBits {
			isPresent := (tmp>>j)&1 == 1

			if !isPresent {
				continue
			}

			if !isFirst {
				sb.WriteString("*")
			} else {
				isFirst = false
			}

			fmt.Fprintf(sb, "x%d", i*monomChunkBits+j+1)
		}
	}
}

func (self Monomial) String() string {
	sb := &strings.Builder{}
	self.Write(sb)

	return sb.String()
}

func (self Monomial) Subscripts() iter.Seq[poly.Subscript] {
	n := len(self.data)

	return func(yield func(poly.Subscript) bool) {
		for i := range n {
			tmp := self.data[n-i-1]

			for j := range monomChunkBits {
				isPresent := (tmp>>j)&1 == 1

				if !isPresent {
					continue
				}

				sub := poly.Subscript(i*monomChunkBits + j)

				if !yield(sub) {
					return
				}
			}
		}
	}
}

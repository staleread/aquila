
#align(center, [*ДОДАТОК А*])
#align(center, [Програмний код])
#linebreak()

Файл lib/internal/automata/math/adder.go

```
package math

// Adds two pre-sorted polinomials using two-pointer merge
// Appends the resulting monomials to the 'dst' buffer.
func AddPolynomials(dst, p1, p2 []Monomial) []Monomial {
	var ptr1, ptr2 int

	for ptr1 < len(p1) && ptr2 < len(p2) {
		m1 := p1[ptr1]
		m2 := p2[ptr2]

		cmp := CompareMonomials(m1, m2)

		switch {
		case cmp > 0:
			dst = append(dst, m1)
			ptr1++

		case cmp < 0:
			dst = append(dst, m2)
			ptr2++

		default:
			// Identical terms cancel each other (a ^ a = 0), so skip
			ptr1++
			ptr2++
		}
	}

	if ptr1 < len(p1) {
		dst = append(dst, p1[ptr1:]...)
	}

	if ptr2 < len(p2) {
		dst = append(dst, p2[ptr2:]...)
	}
	return dst
}
```

#linebreak()
Файл lib/internal/automata/math/adder_test.go

```
package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestAddPolynomials(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   []math.Monomial
		expected []math.Monomial
	}{
		{
			name:     "Both empty",
			p1:       []math.Monomial{},
			p2:       []math.Monomial{},
			expected: []math.Monomial{},
		},
		{
			name:     "One empty",
			p1:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			p2:       []math.Monomial{},
			expected: []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
		},
		{
			name:     "Disjoint sets",
			p1:       []math.Monomial{math.NewMonomial(3), math.NewMonomial(1)},
			p2:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			expected: []math.Monomial{math.NewMonomial(3), math.NewMonomial(2), math.NewMonomial(1), math.NewMonomial(0)},
		},
		{
			name:     "Identical (full cancellation)",
			p1:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			p2:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			expected: []math.Monomial{},
		},
		{
			name:     "Partial overlap",
			p1:       []math.Monomial{math.NewMonomial(3), math.NewMonomial(2)},
			p2:       []math.Monomial{math.NewMonomial(2), math.NewMonomial(0)},
			expected: []math.Monomial{math.NewMonomial(3), math.NewMonomial(0)}, // x2 cancels
		},
		{
			name:     "Interleaved",
			p1:       []math.Monomial{math.NewMonomial(4), math.NewMonomial(2), math.NewMonomial(0)},
			p2:       []math.Monomial{math.NewMonomial(3), math.NewMonomial(2), math.NewMonomial(1)},
			expected: []math.Monomial{math.NewMonomial(4), math.NewMonomial(3), math.NewMonomial(1), math.NewMonomial(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := math.AddPolynomials(nil, tt.p1, tt.p2)

			if len(got) != len(tt.expected) {
				t.Fatalf("Expected length %d, got %d. Got: %v", len(tt.expected), len(got), got)
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("At index %d: got %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
```

Файл lib/internal/automata/math/bitset128.go

```
//go:build block128

package math

import (
	"encoding/binary"
	"math/bits"
)

const (
	BitsetSize  = 128
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset [2]uint64

func NewBitset(subs ...Subscript) Bitset {
	var s Bitset
	for _, i := range subs {
		s.SetAt(i, 1)
	}
	return s
}

func (s *Bitset) Read(src []byte) {
	_ = src[15] // BCE check

	s[0] = binary.LittleEndian.Uint64(src[0:8])
	s[1] = binary.LittleEndian.Uint64(src[8:16])
}

func (s *Bitset) Write(dst []byte) {
	_ = dst[15] // BCE check

	binary.LittleEndian.PutUint64(dst[0:8], s[0])
	binary.LittleEndian.PutUint64(dst[8:16], s[1])
}

func (s Bitset) At(sub Subscript) uint8 {
	return uint8((s[sub/64] >> (sub % 64)) & 1)
}

func (s *Bitset) SetAt(sub Subscript, bit uint8) {
	wordIdx := sub / 64
	shift := sub % 64

	s[wordIdx] = (s[wordIdx] &^ (uint64(1) << shift)) | uint64(bit&1)<<shift
}

func (s Bitset) Compare(other Bitset) int {
	for i := range len(s) {
		if s[i] == other[i] {
			continue
		}
		if s[i] < other[i] {
			return -1
		}
		return 1
	}
	return 0
}

func (s Bitset) Or(other Bitset) Bitset {
	return Bitset{
		s[0] | other[0],
		s[1] | other[1],
	}
}

func (s *Bitset) XorWith(other Bitset) {
	s[0] ^= other[0]
	s[1] ^= other[1]
}

func (s Bitset) Contains(other Bitset) bool {
	return s[0]&other[0] == other[0] && s[1]&other[1] == other[1]
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	for i, w := range s {
		wordOffset := i * 64
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			dst = append(dst, Subscript(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
}
```

Файл lib/internal/automata/math/bitset128_test.go

```
//go:build block128

package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestBitset_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
		readBit math.Subscript
		want    uint8
	}{
		{"first bit of first word", []math.Subscript{0}, 0, 1},
		{"last bit of first word", []math.Subscript{63}, 63, 1},
		{"first bit of second word", []math.Subscript{64}, 64, 1},
		{"last bit of second word", []math.Subscript{127}, 127, 1},
		{"unset bit reads zero", []math.Subscript{5}, 10, 0},
		{"multiple bits set, read one", []math.Subscript{3, 70, 120}, 70, 1},
		{"clear a previously set bit", nil, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBitset_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero state", nil},
		{"single bit word 0", []math.Subscript{7}},
		{"single bit word 1", []math.Subscript{100}},
		{"bits in both words", []math.Subscript{10, 70, 127}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := math.NewBitset(tt.setBits...)

			buf := make([]byte, 16)
			src.Write(buf)

			var dst math.Bitset
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBitset_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b math.Bitset
		want int
	}{
		{"equal", math.NewBitset(5), math.NewBitset(5), 0},
		{"a < b same word", math.NewBitset(0), math.NewBitset(1), -1},
		{"a > b same word", math.NewBitset(1), math.NewBitset(0), 1},
		{"a < b second word", math.NewBitset(64), math.NewBitset(65), -1},
		{"a > b second word", math.NewBitset(65), math.NewBitset(64), 1},
		{"both zero", math.Bitset{}, math.Bitset{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Compare(tt.b)
			if got != tt.want {
				t.Errorf("Compare = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBitset_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     math.Bitset
		wantBits []math.Subscript
	}{
		{"disjoint bits", math.NewBitset(0), math.NewBitset(1), []math.Subscript{0, 1}},
		{"overlapping bits", math.NewBitset(0, 1), math.NewBitset(1, 2), []math.Subscript{0, 1, 2}},
		{"across words", math.NewBitset(63), math.NewBitset(64), []math.Subscript{63, 64}},
		{"with zero", math.NewBitset(10, 20), math.Bitset{}, []math.Subscript{10, 20}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := math.NewBitset(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBitset_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     math.Bitset
		other math.Bitset
		want  bool
	}{
		{"exact match", math.NewBitset(10, 70), math.NewBitset(10, 70), true},
		{"proper superset", math.NewBitset(10, 70, 127), math.NewBitset(10, 70), true},
		{"missing one bit", math.NewBitset(10), math.NewBitset(10, 70), false},
		{"other is zero", math.NewBitset(5), math.Bitset{}, true},
		{"both zero", math.Bitset{}, math.Bitset{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Contains(tt.other)
			if got != tt.want {
				t.Errorf("Contains = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitset_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero state", nil},
		{"single bit word 0", []math.Subscript{0}},
		{"single bit word 1", []math.Subscript{65}},
		{"multiple bits both words", []math.Subscript{0, 31, 64, 127}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.Subscripts(nil)

			if len(got) != len(tt.setBits) {
				t.Fatalf("Subscripts len = %d, want %d", len(got), len(tt.setBits))
			}
			for i, s := range tt.setBits {
				if got[i] != s {
					t.Errorf("Subscripts[%d] = %d, want %d", i, got[i], s)
				}
			}
		})
	}
}
```

Файл lib/internal/automata/math/bitset16.go

```
//go:build block16

package math

import (
	"encoding/binary"
	"math/bits"
)

const (
	BitsetSize  = 16
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset uint16

func NewBitset(subs ...Subscript) Bitset {
	var b Bitset
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

func (s *Bitset) Read(src []byte) {
	_ = src[1] // BCE check
	*s = Bitset(binary.LittleEndian.Uint16(src[0:2]))
}

func (s Bitset) Write(dst []byte) {
	_ = dst[1] // BCE check
	binary.LittleEndian.PutUint16(dst[0:2], uint16(s))
}

func (s Bitset) At(idx Subscript) uint8 {
	return uint8((s >> idx) & 1)
}

func (s *Bitset) SetAt(idx Subscript, bit uint8) {
	*s = (*s &^ (Bitset(1) << idx)) | (Bitset(bit&1) << idx)
}

func (s Bitset) Compare(other Bitset) int {
	if s == other {
		return 0
	}
	if s < other {
		return -1
	}
	return 1
}

func (s Bitset) Or(other Bitset) Bitset {
	return s | other
}

func (s *Bitset) XorWith(other Bitset) {
	*s ^= other
}

func (s Bitset) Contains(other Bitset) bool {
	return s&other == other
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	w := uint16(s)
	for w != 0 {
		tz := bits.TrailingZeros16(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}
```

Файл lib/internal/automata/math/bitset16_test.go

```
//go:build block16

package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestBitset_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
		readBit math.Subscript
		want    uint8
	}{
		{"first bit", []math.Subscript{0}, 0, 1},
		{"last bit", []math.Subscript{15}, 15, 1},
		{"unset bit reads zero", []math.Subscript{2}, 5, 0},
		{"multiple bits set, read one", []math.Subscript{1, 3, 10}, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBitset_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero state", nil},
		{"single bit", []math.Subscript{10}},
		{"multiple bits", []math.Subscript{1, 8, 15}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := math.NewBitset(tt.setBits...)

			buf := make([]byte, 2)
			src.Write(buf)

			var dst math.Bitset
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBitset_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b math.Bitset
		want int
	}{
		{"equal", math.NewBitset(5), math.NewBitset(5), 0},
		{"a < b", math.NewBitset(0), math.NewBitset(1), -1},
		{"a > b", math.NewBitset(1), math.NewBitset(0), 1},
		{"both zero", math.Bitset(0), math.Bitset(0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Compare(tt.b)
			if got != tt.want {
				t.Errorf("Compare = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBitset_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     math.Bitset
		wantBits []math.Subscript
	}{
		{"disjoint bits", math.NewBitset(0), math.NewBitset(1), []math.Subscript{0, 1}},
		{"overlapping bits", math.NewBitset(0, 1), math.NewBitset(1, 2), []math.Subscript{0, 1, 2}},
		{"with zero", math.NewBitset(2, 12), math.Bitset(0), []math.Subscript{2, 12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := math.NewBitset(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBitset_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     math.Bitset
		other math.Bitset
		want  bool
	}{
		{"exact match", math.NewBitset(2, 12), math.NewBitset(2, 12), true},
		{"proper superset", math.NewBitset(2, 4, 12), math.NewBitset(2, 12), true},
		{"missing one bit", math.NewBitset(2), math.NewBitset(2, 12), false},
		{"other is zero", math.NewBitset(5), math.Bitset(0), true},
		{"both zero", math.Bitset(0), math.Bitset(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Contains(tt.other)
			if got != tt.want {
				t.Errorf("Contains = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitset_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero state", nil},
		{"single bit", []math.Subscript{7}},
		{"multiple bits", []math.Subscript{0, 7, 15}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.Subscripts(nil)

			if len(got) != len(tt.setBits) {
				t.Fatalf("Subscripts len = %d, want %d", len(got), len(tt.setBits))
			}
			for i, s := range tt.setBits {
				if got[i] != s {
					t.Errorf("Subscripts[%d] = %d, want %d", i, got[i], s)
				}
			}
		})
	}
}
```

Файл lib/internal/automata/math/bitset32.go

```
//go:build block32

package math

import (
	"encoding/binary"
	"math/bits"
)

const (
	BitsetSize  = 32
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset uint32

func NewBitset(subs ...Subscript) Bitset {
	var b Bitset
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

func (s *Bitset) Read(src []byte) {
	_ = src[3] // BCE check
	*s = Bitset(binary.LittleEndian.Uint32(src[0:4]))
}

func (s Bitset) Write(dst []byte) {
	_ = dst[3] // BCE check
	binary.LittleEndian.PutUint32(dst[0:4], uint32(s))
}

func (s Bitset) At(idx Subscript) uint8 {
	return uint8((s >> idx) & 1)
}

func (s *Bitset) SetAt(idx Subscript, bit uint8) {
	*s = (*s &^ (Bitset(1) << idx)) | (Bitset(bit&1) << idx)
}

func (s Bitset) Compare(other Bitset) int {
	if s == other {
		return 0
	}
	if s < other {
		return -1
	}
	return 1
}

func (s Bitset) Or(other Bitset) Bitset {
	return s | other
}

func (s *Bitset) XorWith(other Bitset) {
	*s ^= other
}

func (s Bitset) Contains(other Bitset) bool {
	return s&other == other
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	w := uint32(s)
	for w != 0 {
		tz := bits.TrailingZeros32(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}
```

Файл lib/internal/automata/math/bitset32_test.go

```
//go:build block32

package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestBitset_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
		readBit math.Subscript
		want    uint8
	}{
		{"first bit", []math.Subscript{0}, 0, 1},
		{"last bit", []math.Subscript{31}, 31, 1},
		{"unset bit reads zero", []math.Subscript{2}, 5, 0},
		{"multiple bits set, read one", []math.Subscript{1, 3, 20}, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBitset_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero bitset", nil},
		{"single bit", []math.Subscript{20}},
		{"multiple bits", []math.Subscript{1, 16, 31}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := math.NewBitset(tt.setBits...)

			buf := make([]byte, 4)
			src.Write(buf)

			var dst math.Bitset
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBitset_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b math.Bitset
		want int
	}{
		{"equal", math.NewBitset(5), math.NewBitset(5), 0},
		{"a < b", math.NewBitset(0), math.NewBitset(1), -1},
		{"a > b", math.NewBitset(1), math.NewBitset(0), 1},
		{"both zero", math.Bitset(0), math.Bitset(0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Compare(tt.b)
			if got != tt.want {
				t.Errorf("Compare = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBitset_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     math.Bitset
		wantBits []math.Subscript
	}{
		{"disjoint bits", math.NewBitset(0), math.NewBitset(1), []math.Subscript{0, 1}},
		{"overlapping bits", math.NewBitset(0, 1), math.NewBitset(1, 2), []math.Subscript{0, 1, 2}},
		{"with zero", math.NewBitset(2, 28), math.Bitset(0), []math.Subscript{2, 28}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := math.NewBitset(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBitset_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     math.Bitset
		other math.Bitset
		want  bool
	}{
		{"exact match", math.NewBitset(2, 28), math.NewBitset(2, 28), true},
		{"proper superset", math.NewBitset(2, 4, 28), math.NewBitset(2, 28), true},
		{"missing one bit", math.NewBitset(2), math.NewBitset(2, 28), false},
		{"other is zero", math.NewBitset(5), math.Bitset(0), true},
		{"both zero", math.Bitset(0), math.Bitset(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Contains(tt.other)
			if got != tt.want {
				t.Errorf("Contains = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitset_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero state", nil},
		{"single bit", []math.Subscript{25}},
		{"multiple bits", []math.Subscript{0, 25, 31}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.Subscripts(nil)

			if len(got) != len(tt.setBits) {
				t.Fatalf("Subscripts len = %d, want %d", len(got), len(tt.setBits))
			}
			for i, s := range tt.setBits {
				if got[i] != s {
					t.Errorf("Subscripts[%d] = %d, want %d", i, got[i], s)
				}
			}
		})
	}
}
```

Файл lib/internal/automata/math/bitset64.go

```
//go:build block64

package math

import (
	"encoding/binary"
	"math/bits"
)

const (
	BitsetSize  = 64
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset uint64

func NewBitset(subs ...Subscript) Bitset {
	var b Bitset
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

func (s *Bitset) Read(src []byte) {
	_ = src[7] // BCE check
	*s = Bitset(binary.LittleEndian.Uint64(src[0:8]))
}

func (s Bitset) Write(dst []byte) {
	_ = dst[7] // BCE check
	binary.LittleEndian.PutUint64(dst[0:8], uint64(s))
}

func (s Bitset) At(idx Subscript) uint8 {
	return uint8((s >> idx) & 1)
}

func (s *Bitset) SetAt(idx Subscript, bit uint8) {
	*s = (*s &^ (Bitset(1) << idx)) | (Bitset(bit&1) << idx)
}

func (s Bitset) Compare(other Bitset) int {
	if s == other {
		return 0
	}
	if s < other {
		return -1
	}
	return 1
}

func (s Bitset) Or(other Bitset) Bitset {
	return s | other
}

func (s *Bitset) XorWith(other Bitset) {
	*s ^= other
}

func (s Bitset) Contains(other Bitset) bool {
	return s&other == other
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	w := uint64(s)
	for w != 0 {
		tz := bits.TrailingZeros64(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}
```

Файл lib/internal/automata/math/bitset8.go

```
//go:build block8

package math

import "math/bits"

const (
	BitsetSize  = 8
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset uint8

func NewBitset(subs ...Subscript) Bitset {
	var b Bitset
	for _, i := range subs {
		b.SetAt(i, 1)
	}
	return b
}

func (s *Bitset) Read(src []byte) {
	_ = src[0] // BCE check
	*s = Bitset(src[0])
}

func (s Bitset) Write(dst []byte) {
	_ = dst[0] // BCE check
	dst[0] = uint8(s)
}

func (s Bitset) At(idx Subscript) uint8 {
	return uint8((s >> idx) & 1)
}

func (s *Bitset) SetAt(idx Subscript, bit uint8) {
	*s = (*s &^ (Bitset(1) << idx)) | (Bitset(bit&1) << idx)
}

func (s Bitset) Compare(other Bitset) int {
	if s == other {
		return 0
	}
	if s < other {
		return -1
	}
	return 1
}

func (s Bitset) Or(other Bitset) Bitset {
	return s | other
}

func (s *Bitset) XorWith(other Bitset) {
	*s ^= other
}

func (s Bitset) Contains(other Bitset) bool {
	return s&other == other
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	w := uint8(s)
	for w != 0 {
		tz := bits.TrailingZeros8(w)
		dst = append(dst, Subscript(tz))
		w &= w - 1
	}
	return dst
}
```

Файл lib/internal/automata/math/bitset96.go

```
//go:build block96 || (!block8 && !block16 && !block32 && !block64 && !block128)

package math

import (
	"encoding/binary"
	"math/bits"
)

const (
	BitsetSize  = 96
	BitsetBytes = BitsetSize / 8
)

type Subscript = uint8
type Bitset [3]uint32

func NewBitset(subs ...Subscript) Bitset {
	var s Bitset
	for _, i := range subs {
		s.SetAt(i, 1)
	}
	return s
}

func (s *Bitset) Read(src []byte) {
	_ = src[11] // BCE check

	s[0] = binary.LittleEndian.Uint32(src[0:4])
	s[1] = binary.LittleEndian.Uint32(src[4:8])
	s[2] = binary.LittleEndian.Uint32(src[8:12])
}

func (s *Bitset) Write(dst []byte) {
	_ = dst[11] // BCE check

	binary.LittleEndian.PutUint32(dst[0:4], s[0])
	binary.LittleEndian.PutUint32(dst[4:8], s[1])
	binary.LittleEndian.PutUint32(dst[8:12], s[2])
}

func (s Bitset) At(idx Subscript) uint8 {
	return uint8((s[idx/32] >> (idx % 32)) & 1)
}

func (s *Bitset) SetAt(idx Subscript, bit uint8) {
	wordIdx := idx / 32
	shift := idx % 32

	s[wordIdx] = (s[wordIdx] &^ (uint32(1) << shift)) | uint32(bit&1)<<shift
}

func (s Bitset) Compare(other Bitset) int {
	for i := range len(s) {
		if s[i] == other[i] {
			continue
		}
		if s[i] < other[i] {
			return -1
		}
		return 1
	}
	return 0
}

func (s Bitset) Or(other Bitset) Bitset {
	return Bitset{
		s[0] | other[0],
		s[1] | other[1],
		s[2] | other[2],
	}
}

func (s *Bitset) XorWith(other Bitset) {
	s[0] ^= other[0]
	s[1] ^= other[1]
	s[2] ^= other[2]
}

func (s Bitset) Contains(other Bitset) bool {
	return s[0]&other[0] == other[0] && s[1]&other[1] == other[1] && s[2]&other[2] == other[2]
}

func (s Bitset) Subscripts(dst []Subscript) []Subscript {
	for i, w := range s {
		wordOffset := i * 32
		for w != 0 {
			tz := bits.TrailingZeros32(w)
			dst = append(dst, Subscript(wordOffset+tz))
			w &= w - 1
		}
	}
	return dst
}
```

Файл lib/internal/automata/math/bitset96_test.go

```
//go:build block96 || (!block8 && !block16 && !block32 && !block64 && !block128)

package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestBitset_SetAt_At(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
		readBit math.Subscript
		want    uint8
	}{
		{"first bit of first word", []math.Subscript{0}, 0, 1},
		{"last bit of first word", []math.Subscript{31}, 31, 1},
		{"first bit of second word", []math.Subscript{32}, 32, 1},
		{"first bit of third word", []math.Subscript{64}, 64, 1},
		{"last bit of third word", []math.Subscript{95}, 95, 1},
		{"unset bit reads zero", []math.Subscript{5}, 10, 0},
		{"multiple bits set, read one", []math.Subscript{3, 40, 80}, 40, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := math.NewBitset(tt.setBits...)
			got := s.At(tt.readBit)
			if got != tt.want {
				t.Errorf("At(%d) = %d, want %d", tt.readBit, got, tt.want)
			}
		})
	}
}

func TestBitset_ReadFrom_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero math", nil},
		{"single bit word 0", []math.Subscript{7}},
		{"single bit word 1", []math.Subscript{40}},
		{"single bit word 2", []math.Subscript{80}},
		{"bits across all words", []math.Subscript{10, 40, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := math.NewBitset(tt.setBits...)

			buf := make([]byte, 12)
			src.Write(buf)

			var dst math.Bitset
			dst.Read(buf)

			if src != dst {
				t.Errorf("round-trip mismatch: wrote %v, read back %v", src, dst)
			}
		})
	}
}

func TestBitset_Compare(t *testing.T) {
	tests := []struct {
		name string
		a, b math.Bitset
		want int
	}{
		{"equal", math.NewBitset(5), math.NewBitset(5), 0},
		{"a < b first word", math.NewBitset(0), math.NewBitset(1), -1},
		{"a > b first word", math.NewBitset(1), math.NewBitset(0), 1},
		{"a < b second word", math.NewBitset(32), math.NewBitset(33), -1},
		{"a > b third word", math.NewBitset(65), math.NewBitset(64), 1},
		{"both zero", math.Bitset{}, math.Bitset{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Compare(tt.b)
			if got != tt.want {
				t.Errorf("Compare = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestBitset_Or(t *testing.T) {
	tests := []struct {
		name     string
		a, b     math.Bitset
		wantBits []math.Subscript
	}{
		{"disjoint bits", math.NewBitset(0), math.NewBitset(1), []math.Subscript{0, 1}},
		{"overlapping bits", math.NewBitset(0, 1), math.NewBitset(1, 2), []math.Subscript{0, 1, 2}},
		{"across words", math.NewBitset(31), math.NewBitset(32), []math.Subscript{31, 32}},
		{"with zero", math.NewBitset(10, 50), math.Bitset{}, []math.Subscript{10, 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)
			want := math.NewBitset(tt.wantBits...)
			if got != want {
				t.Errorf("Or = %v, want %v", got, want)
			}
		})
	}
}

func TestBitset_Contains(t *testing.T) {
	tests := []struct {
		name  string
		b     math.Bitset
		other math.Bitset
		want  bool
	}{
		{"exact match", math.NewBitset(10, 50), math.NewBitset(10, 50), true},
		{"proper superset", math.NewBitset(10, 50, 80), math.NewBitset(10, 50), true},
		{"missing one bit", math.NewBitset(10), math.NewBitset(10, 50), false},
		{"other is zero", math.NewBitset(5), math.Bitset{}, true},
		{"both zero", math.Bitset{}, math.Bitset{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Contains(tt.other)
			if got != tt.want {
				t.Errorf("Contains = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitset_Subscripts(t *testing.T) {
	tests := []struct {
		name    string
		setBits []math.Subscript
	}{
		{"zero math", nil},
		{"single bit word 0", []math.Subscript{0}},
		{"single bit word 1", []math.Subscript{40}},
		{"single bit word 2", []math.Subscript{80}},
		{"multiple bits all words", []math.Subscript{0, 31, 32, 63, 64, 95}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := math.NewBitset(tt.setBits...)
			got := b.Subscripts(nil)

			if len(got) != len(tt.setBits) {
				t.Fatalf("Subscripts len = %d, want %d", len(got), len(tt.setBits))
			}
			for i, s := range tt.setBits {
				if got[i] != s {
					t.Errorf("Subscripts[%d] = %d, want %d", i, got[i], s)
				}
			}
		})
	}
}
```

Файл lib/internal/automata/math/confusion.go

```
package math

import (
	"io"
	"slices"
)

const ConfusionMapBytes = (ConfusionDegree*(ConfusionDegree+1)/2 - 1) * VectorSize

type ConfusionMap struct {
	Data []Subscript
}

func NewConfusionMap(arena []byte) ConfusionMap {
	return ConfusionMap{Data: []Subscript(arena)}
}

func EmptyConfusionMap() ConfusionMap {
	return ConfusionMap{Data: nil}
}

func (m *ConfusionMap) Generate(rnd io.Reader, maxSub Subscript, perm *Permutation) error {
	if _, err := io.ReadFull(rnd, m.Data); err != nil {
		return err
	}

	for i := range m.Data {
		m.Data[i] %= maxSub
	}

	subIdx := 0
	for range VectorSize {
		for j := ConfusionDegree; j > 1; j-- {
			upperBound := int(maxSub) - j

			for k := range j {
				candidateIdx := int(m.Data[subIdx]) % (upperBound + k + 1)
				candidate := perm.Data[candidateIdx]

				duplicate := slices.Contains(m.Data[subIdx-k:subIdx], candidate)

				if duplicate {
					m.Data[subIdx] = perm.Data[upperBound+k]
				} else {
					m.Data[subIdx] = candidate
				}
				subIdx++
			}
		}
	}
	return nil
}

func (m *ConfusionMap) Eval(s Bitset) Vector {
	if m == nil || len(m.Data) == 0 {
		return Vector(0)
	}

	var res Vector
	subIdx := 0

	for i := range VectorSize {
		var sum Vector

		for j := range ConfusionDegree - 1 {
			prod := Vector(1)
			subCnt := ConfusionDegree - j

			for range subCnt {
				bit := s.At(m.Data[subIdx])
				prod &= Vector(bit)

				subIdx++
			}
			sum ^= prod
		}
		res |= sum << i
	}
	return res
}
```

Файл lib/internal/automata/math/confusion2.go

```
//go:build deg2 || (!deg3 && !deg4)

package math

const ConfusionDegree = 2
```

Файл lib/internal/automata/math/confusion3.go

```
//go:build deg3

package math

const ConfusionDegree = 3
```

Файл lib/internal/automata/math/confusion3_test.go

```
//go:build deg3

package math_test

import (
	"bytes"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestConfusionMapEval(t *testing.T) {
	perm := getIdentityPermutation()

	degArena := make([]byte, math.ConfusionMapBytes)
	m := math.NewConfusionMap(degArena)

	// y0 = x0*x1*x2 + x3*x4
	fixture := make([]byte, math.ConfusionMapBytes)
	fixture[0], fixture[1], fixture[2] = 0, 1, 2
	fixture[3], fixture[4] = 3, 4
	fixtureGen := bytes.NewReader(fixture)

	maxSub := math.Subscript(math.VectorSize)

	if err := m.Generate(fixtureGen, maxSub, perm); err != nil {
		t.Fatalf("Failed to generate degusion map: %v", err)
	}

	// Verify initialization
	expectedIds := []byte{0, 1, 2, 3, 4}
	for i, exp := range expectedIds {
		if degArena[i] != exp {
			t.Fatalf("Index at %d: got %d, expected %d", i, degArena[i], exp)
		}
	}

	tests := []struct {
		name        string
		activeBits  []math.Subscript
		expectedBit uint8
	}{
		{
			name:        "All zeros -> 0^0 = 0",
			activeBits:  []math.Subscript{},
			expectedBit: 0,
		},
		{
			name:        "Only bit 5 is set -> 0^0 = 0",
			activeBits:  []math.Subscript{5},
			expectedBit: 0,
		},
		{
			name:        "Bits 3 and 4 are set -> 0^1 = 1",
			activeBits:  []math.Subscript{3, 4},
			expectedBit: 1,
		},
		{
			name:        "Bits 0, 1, 2, 3, 4 are set -> 1^1 = 0",
			activeBits:  []math.Subscript{0, 1, 2, 3, 4},
			expectedBit: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b math.Bitset

			for _, bit := range tc.activeBits {
				b.SetAt(bit, 1)
			}

			actual := m.Eval(b)
			actualBit := uint8(actual & 1)

			if actualBit != tc.expectedBit {
				t.Errorf("Expected %d, got %d", tc.expectedBit, actualBit)
			}
		})
	}
}

func getIdentityPermutation() *math.Permutation {
	permArena := make([]byte, math.PermutationBytes)

	for i := range len(permArena) {
		permArena[i] = math.Subscript(i)
	}
	return math.NewPermutation(permArena)
}
```

Файл lib/internal/automata/math/confusion4.go

```
//go:build deg4

package math

const ConfusionDegree = 4
```

Файл lib/internal/automata/math/monomial.go

```
package math

type Monomial Bitset

var IdentityMonomial Monomial

func NewMonomial(subs ...Subscript) Monomial {
	var b Bitset
	for _, sub := range subs {
		b.SetAt(sub, 1)
	}
	return Monomial(b)
}

func CompareMonomials(a, b Monomial) int {
	return Bitset(a).Compare(Bitset(b))
}

func (m Monomial) Mul(other Monomial) Monomial {
	return Monomial(Bitset(m).Or(Bitset(other)))
}

func (m Monomial) Subscripts(dst []uint8) []uint8 {
	return Bitset(m).Subscripts(dst)
}

func (m *Monomial) SetAt(sub Subscript, bit uint8) {
	(*Bitset)(m).SetAt(sub, bit)
}

func (m Monomial) At(sub Subscript) uint8 {
	return Bitset(m).At(sub)
}

func (m Monomial) Eval(b Bitset) uint8 {
	if m == IdentityMonomial || b.Contains(Bitset(m)) {
		return 1
	}
	return 0
}
```

Файл lib/internal/automata/math/monomial_test.go

```
package math_test

import (
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestNewMonomial(t *testing.T) {
	subs := []math.Subscript{0, 1, 2}
	m := math.NewMonomial(subs...)

	got := m.Subscripts(nil)
	if len(got) != len(subs) {
		t.Fatalf("NewMonomial subscripts length: got %d, want %d", len(got), len(subs))
	}
	for i, s := range subs {
		if got[i] != s {
			t.Errorf("NewMonomial subscript[%d]: got %d, want %d", i, got[i], s)
		}
	}
}

func TestCompareMonomials(t *testing.T) {
	tests := []struct {
		name string
		a, b math.Monomial
		want int
	}{
		{
			name: "equal",
			a:    math.NewMonomial(0),
			b:    math.NewMonomial(0),
			want: 0,
		},
		{
			name: "a < b (lower bit set in a)",
			a:    math.NewMonomial(0),
			b:    math.NewMonomial(1),
			want: -1,
		},
		{
			name: "a > b (higher bit set in a)",
			a:    math.NewMonomial(1),
			b:    math.NewMonomial(0),
			want: 1,
		},
	}

	for _, tt := range tests {
		got := math.CompareMonomials(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("CompareMonomials %s: got %d, want %d", tt.name, got, tt.want)
		}
	}

	if math.BitsetSize >= 96 {
		multiTests := []struct {
			name string
			a, b math.Monomial
			want int
		}{
			{
				name: "a < b (second word differs)",
				a:    math.NewMonomial(64),
				b:    math.NewMonomial(65),
				want: -1,
			},
			{
				name: "a > b (empty vs non-empty second word)",
				a:    math.NewMonomial(64),
				b:    math.NewMonomial(),
				want: 1,
			},
		}
		for _, tt := range multiTests {
			got := math.CompareMonomials(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("CompareMonomials %s: got %d, want %d", tt.name, got, tt.want)
			}
		}
	}
}

func TestMonomial_Mul(t *testing.T) {
	// x0 * x2 = {0, 2}
	m1 := math.NewMonomial(0)
	m2 := math.NewMonomial(2)
	want := math.NewMonomial(0, 2)

	got := m1.Mul(m2)
	if math.CompareMonomials(got, want) != 0 {
		t.Errorf("Mul(%v, %v) = %v, want %v", m1, m2, got, want)
	}
}

func TestMonomial_Eval(t *testing.T) {
	var m math.Monomial
	var b0, b1, b2, b3 math.Bitset
	var tests []struct {
		name string
		b    math.Bitset
		want uint8
	}

	if math.BitsetSize >= 96 {
		m = math.NewMonomial(0, 31, 65)
		b0 = math.NewBitset(0, 31, 65)        // exact match
		b1 = math.NewBitset(0, 31, 65, 5, 66) // superset
		b2 = math.NewBitset(0, 65)            // missing x31
		b3 = math.NewBitset(0, 31)            // missing x65

		tests = []struct {
			name string
			b    math.Bitset
			want uint8
		}{
			{"exact match", b0, 1},
			{"superset", b1, 1},
			{"missing x31", b2, 0},
			{"missing x65", b3, 0},
		}
	} else if math.BitsetSize >= 16 {
		m = math.NewMonomial(0, 5, 15)
		b0 = math.NewBitset(0, 5, 15)
		b1 = math.NewBitset(0, 5, 15, 2, 7)
		b2 = math.NewBitset(0, 15)
		b3 = math.NewBitset(0, 5)

		tests = []struct {
			name string
			b    math.Bitset
			want uint8
		}{
			{"exact match", b0, 1},
			{"superset", b1, 1},
			{"missing x5", b2, 0},
			{"missing x15", b3, 0},
		}
	} else {
		m = math.NewMonomial(0, 2, 7)
		b0 = math.NewBitset(0, 2, 7)
		b1 = math.NewBitset(0, 2, 7, 5)
		b2 = math.NewBitset(0, 7)
		b3 = math.NewBitset(0, 2)

		tests = []struct {
			name string
			b    math.Bitset
			want uint8
		}{
			{"exact match", b0, 1},
			{"superset", b1, 1},
			{"missing x2", b2, 0},
			{"missing x7", b3, 0},
		}
	}

	for _, tt := range tests {
		got := m.Eval(tt.b)
		if got != tt.want {
			t.Errorf("Eval %s: got %d, want %d", tt.name, got, tt.want)
		}
	}
}
```

Файл lib/internal/automata/math/multiplier.go

```
package math

import "slices"

// Multiplies two polynomials and overwrites the 'dst' buffer with the resulting
// polynomial. The resulting monomials are sorted and deduplicated.
func MultiplyPolynomials(dst, p1, p2 []Monomial) []Monomial {
	if len(p1) == 0 || len(p2) == 0 {
		return dst[:0]
	}

	total := len(p1) * len(p2)

	if cap(dst) < total {
		dst = make([]Monomial, total)
	} else {
		dst = dst[:total]
	}

	for i, m1 := range p1 {
		for j, m2 := range p2 {
			dst[i*len(p2)+j] = m1.Mul(m2)
		}
	}

	slices.SortFunc(dst, func(a, b Monomial) int {
		return CompareMonomials(b, a)
	})

	writeIdx := 0
	for readIdx := 0; readIdx < len(dst); {
		curr := dst[readIdx]
		count := 1
		readIdx++

		for readIdx < len(dst) && dst[readIdx] == curr {
			count++
			readIdx++
		}

		if count%2 != 0 {
			dst[writeIdx] = curr
			writeIdx++
		}
	}

	return dst[:writeIdx]
}
```

Файл lib/internal/automata/math/multiplier_test.go

```
package math_test

import (
	"slices"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func sortPolynomial(p []math.Monomial) {
	slices.SortFunc(p, func(a, b math.Monomial) int {
		return math.CompareMonomials(b, a) // Descending
	})
}

func TestPolynomialMultiplier_Multiply(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   []math.Monomial
		expected []math.Monomial
	}{
		{
			name:     "Empty polynomials",
			p1:       []math.Monomial{},
			p2:       []math.Monomial{math.NewMonomial(0)},
			expected: []math.Monomial{},
		},
		{
			name:     "Single monomial multiplication",
			p1:       []math.Monomial{math.NewMonomial(0)},    // x0
			p2:       []math.Monomial{math.NewMonomial(1)},    // x1
			expected: []math.Monomial{math.NewMonomial(0, 1)}, // x0*x1
		},
		{
			name: "Basic distribution (x0 + x1) * x2",
			p1:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)}, // x1, x0 (descending)
			p2:   []math.Monomial{math.NewMonomial(2)},                      // x2
			// (x1+x0)*x2 = x1x2 + x0x2
			expected: []math.Monomial{math.NewMonomial(1, 2), math.NewMonomial(0, 2)},
		},
		{
			name: "Cancellation x0 * (x0 + x1) = x0 + x0x1",
			p1:   []math.Monomial{math.NewMonomial(0)},
			p2:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)}, // x1, x0
			// x0*x1 = x0x1, x0*x0 = x0
			expected: []math.Monomial{math.NewMonomial(0, 1), math.NewMonomial(0)},
		},
		{
			name: "Duplicate cancellation (x0 + x1) * (x0 + x1) = x0 + x1 + x0x1 + x0x1 = x0 + x1",
			p1:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)},
			p2:   []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)},
			// x1*x1=x1, x1*x0=x0x1, x0*x1=x0x1, x0*x0=x0 → x0x1 cancels
			expected: []math.Monomial{math.NewMonomial(1), math.NewMonomial(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure inputs are sorted as required
			sortPolynomial(tt.p1)
			sortPolynomial(tt.p2)

			got := math.MultiplyPolynomials(nil, tt.p1, tt.p2)

			if len(got) != len(tt.expected) {
				t.Fatalf("Expected length %d, got %d. Got: %v", len(tt.expected), len(got), got)
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("At index %d: got %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
```

Файл lib/internal/automata/math/permutation.go

```
package math

import "io"

const PermutationSize = BitsetSize
const PermutationBytes = PermutationSize

type Permutation struct {
	Data []Subscript
}

func NewPermutation(arena []byte) *Permutation {
	return &Permutation{
		Data: arena,
	}
}

func (p *Permutation) Generate(rnd io.Reader, tmp []byte) error {
	const n Subscript = PermutationSize

	for i := range n {
		p.Data[i] = i
	}

	if _, err := io.ReadFull(rnd, tmp); err != nil {
		return err
	}

	for i := range n - 1 {
		j := tmp[i]%(n-i) + i
		p.Data[i], p.Data[j] = p.Data[j], p.Data[i]
	}
	return nil
}

func (p *Permutation) Gather(b *Bitset, foldIdx int) Vector {
	var res Vector
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := p.Data[offset+i]
		res |= Vector(b.At(idx)) << i
	}
	return res
}

func (p *Permutation) Scatter(b *Bitset, foldIdx int, v Vector) {
	offset := foldIdx * VectorSize

	for i := range VectorSize {
		idx := p.Data[offset+i]
		bit := uint8((v >> i) & 1)
		b.SetAt(idx, bit)
	}
}
```

Файл lib/internal/automata/math/sle.go

```
package math

import (
	"io"
	"unsafe"
)

const (
	SLESize  = VectorSize * VectorSize
	SLEBytes = VectorBytes * VectorSize
)

type SLE struct {
	data []Vector
}

func NewSLE(arena []byte) SLE {
	data := unsafe.Slice((*Vector)(unsafe.Pointer(&arena[0])), len(arena)/VectorBytes)

	return SLE{data}
}

func (s *SLE) Generate(rnd io.Reader) error {
	byteView := unsafe.Slice((*byte)(unsafe.Pointer(&s.data[0])), len(s.data)*VectorBytes)
	if _, err := io.ReadFull(rnd, byteView); err != nil {
		return err
	}

	for i := range VectorSize {
		s.data[i] |= 1 << i
	}
	return nil
}

func (s *SLE) Solve(b Vector) Vector {
	return s.substituteBackward(s.substituteForward(b))
}

func (s *SLE) Eval(x Vector) Vector {
	return s.multiplyLower(s.multiplyUpper(x))
}

func (s *SLE) Coefs(dst []Vector) {
	for i := range VectorSize {
		iMask := ^Vector((1 << i) - 1)
		dstRow := s.data[i] & iMask

		for k := range i {
			if (s.data[i]>>k)&1 != 1 {
				continue
			}
			kMask := ^Vector((1 << k) - 1)
			dstRow ^= s.data[k] & kMask
		}
		dst[i] = dstRow
	}
}

func (s *SLE) substituteForward(v Vector) Vector {
	var res Vector

	for i, row := range s.data {
		res |= ((v >> i) ^ row.Dot(res)) & 1 << i
	}
	return res
}

func (s *SLE) substituteBackward(v Vector) Vector {
	var res Vector

	for i := VectorSize - 1; i >= 0; i-- {
		row := s.data[i]
		res |= ((v >> i) ^ row.Dot(res)) & 1 << i
	}
	return res
}

func (s *SLE) multiplyLower(v Vector) Vector {
	var res Vector

	for i, row := range s.data {
		mask := Vector((1 << (i + 1)) - 1)
		res |= row.Dot(v&mask) << i
	}
	return res
}

func (s *SLE) multiplyUpper(v Vector) Vector {
	var res Vector

	for i, row := range s.data {
		mask := ^Vector((1 << i) - 1)
		res |= row.Dot(v&mask) << i
	}
	return res
}
```

Файл lib/internal/automata/math/sle_test.go

```
package math_test

import (
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/math"
)

func TestSLEEvalSolve(t *testing.T) {
	arena := make([]byte, math.SLEBytes)
	sle := math.NewSLE(arena)

	rng := rand.New(rand.NewSource(42))

	if err := sle.Generate(rng); err != nil {
		t.Fatalf("Failed to generate SLE: %v", err)
	}

	for range 100 {
		x := math.Vector(rng.Uint32())

		y := sle.Eval(x)
		xPrime := sle.Solve(y)

		if x != xPrime {
			t.Errorf("Eval/Solve mismatch: x=%04x, y=%04x, x'=%04x", x, y, xPrime)
		}
	}
}
```

Файл lib/internal/automata/math/vector16.go

```
//go:build fold16 || (!fold8 && !fold32)

package math

import "math/bits"

const (
	VectorSize  = 16
	VectorBytes = VectorSize / 8
)

type Vector uint16

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount16(uint16(v & other))
	return Vector(ones & 1)
}
```

Файл lib/internal/automata/math/vector32.go

```
//go:build fold32

package math

import "math/bits"

const (
	VectorSize  = 32
	VectorBytes = VectorSize / 8
)

type Vector uint32

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount32(uint32(v & other))
	return Vector(ones & 1)
}
```

Файл lib/internal/automata/math/vector8.go

```
//go:build fold8

package math

import "math/bits"

const (
	VectorSize  = 8
	VectorBytes = VectorSize / 8
)

type Vector uint8

func (v Vector) Dot(other Vector) Vector {
	ones := bits.OnesCount8(uint8(v & other))
	return Vector(ones & 1)
}
```

Файл lib/internal/automata/invertible/ca.go

```
package invertible

import (
	"fmt"
	"io"

	"github.com/staleread/aquila/internal/automata/config"
	"github.com/staleread/aquila/internal/automata/math"
)

const (
	StateSize            = math.BitsetSize
	StateBytes           = math.BitsetBytes
	InitialArenaCapacity = 7_929
	RuleBytes            = RuleFoldsBytes + math.PermutationBytes
	CABytes              = StateBytes + RuleBytes*RulesCount
)

type State = math.Bitset

type CA struct {
	arena []byte
	shift State
}

func NewCA() *CA {
	return &CA{
		arena: make([]byte, CABytes),
	}
}

func (ca *CA) Load(src io.Reader) error {
	if err := config.Current.Check(src); err != nil {
		return err
	}
	if _, err := io.ReadFull(src, ca.arena); err != nil {
		return err
	}
	ca.shift.Read(ca.arena[:StateBytes])
	return nil
}

func (ca *CA) Save(dst io.Writer) error {
	if err := config.Current.Write(dst); err != nil {
		return err
	}
	_, err := dst.Write(ca.arena)
	return err
}

func (ca *CA) Generate(rnd io.Reader) error {
	if _, err := io.ReadFull(rnd, ca.arena[:StateBytes]); err != nil {
		return fmt.Errorf("failed to generate affine shift: %w", err)
	}
	ca.shift.Read(ca.arena[:StateBytes])

	for i := range RulesCount {
		rule := ca.getRule(i)

		if err := rule.Generate(rnd); err != nil {
			return fmt.Errorf("failed to generate rule %d: %w", i, err)
		}
	}
	return nil
}

func (ca *CA) Apply(dst, src []byte) {
	var block State
	block.Read(src)

	for i := range RulesCount {
		rule := ca.getRule(i)
		rule.Apply(&block)
	}

	block.XorWith(ca.shift)
	block.Write(dst)
}

func (ca *CA) Revert(dst, src []byte) {
	var block State
	block.Read(src)

	block.XorWith(ca.shift)

	for i := RulesCount - 1; i >= 0; i-- {
		rule := ca.getRule(i)
		rule.Revert(&block)
	}

	block.Write(dst)
}

func (ca *CA) getRule(idx int) Rule {
	offset := StateBytes + idx*RuleBytes

	return Rule{
		arena: ca.arena[offset : offset+RuleBytes],
	}
}
```

Файл lib/internal/automata/invertible/ca_test.go

```
package invertible_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestCAInvertibility(t *testing.T) {
	ca := invertible.NewCA()
	rng := rand.New(rand.NewSource(42))

	if err := ca.Generate(rng); err != nil {
		t.Fatalf("Failed to generate CA: %v", err)
	}

	src := make([]byte, invertible.StateBytes)
	intermediate := make([]byte, invertible.StateBytes)
	reverted := make([]byte, invertible.StateBytes)

	for range 10 {
		rng.Read(src)

		original := make([]byte, invertible.StateBytes)
		copy(original, src)

		ca.Apply(intermediate, src)
		if bytes.Equal(intermediate, src) {
			t.Errorf("Apply did not change the block (highly unlikely)")
		}

		ca.Revert(reverted, intermediate)
		if !bytes.Equal(reverted, original) {
			t.Errorf("Revert failed: original=%v, reverted=%v", original, reverted)
		}
	}
}
```

Файл lib/internal/automata/invertible/derivation.go

```
package invertible

import (
	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/math"
)

func (ca *CA) DeriveGeneralCA() (*general.CA, error) {
	polynomials := CompileRegistry(ca)
	masterArena := make([]math.Monomial, 0, InitialArenaCapacity*StateSize)
	var offsets [StateSize - 1]uint32

	stateArena := make([]math.Monomial, 0, InitialArenaCapacity)
	prodSrcArena := make([]math.Monomial, 0, InitialArenaCapacity)
	prodDstArena := make([]math.Monomial, 0, InitialArenaCapacity)
	sumArena := make([]math.Monomial, 0, InitialArenaCapacity)
	sumScratchArena := make([]math.Monomial, 0, InitialArenaCapacity)

	for i := range StateSize {
		stateArena = stateArena[:0]
		stateArena = append(stateArena, polynomials.GetPolynomial(RulesCount-1, i)...)

		var subs [StateSize]uint8
		for j := RulesCount - 2; j >= 0; j-- {
			for _, monom := range stateArena {
				subsSlice := monom.Subscripts(subs[:0])
				degree := len(subsSlice)

				if degree == 0 {
					continue
				}

				firstSubscript := int(subsSlice[0])
				firstPoly := polynomials.GetPolynomial(j, firstSubscript)

				if degree == 1 {
					sumScratchArena = sumScratchArena[:0]
					sumScratchArena = math.AddPolynomials(sumScratchArena, sumArena, firstPoly)

					sumArena, sumScratchArena = sumScratchArena, sumArena
					continue
				}

				prodSrcArena = prodSrcArena[:0]
				prodSrcArena = append(prodSrcArena, firstPoly...)

				for _, sub := range subsSlice[1:] {
					nextPoly := polynomials.GetPolynomial(j, int(sub))

					prodDstArena = math.MultiplyPolynomials(prodDstArena, prodSrcArena, nextPoly)

					prodSrcArena, prodDstArena = prodDstArena, prodSrcArena
				}

				sumScratchArena = sumScratchArena[:0]
				sumScratchArena = math.AddPolynomials(sumScratchArena, sumArena, prodSrcArena)

				sumArena, sumScratchArena = sumScratchArena, sumArena
			}

			stateArena = stateArena[:0]
			stateArena = append(stateArena, sumArena...)

			sumArena = sumArena[:0]
		}

		masterArena = append(masterArena, stateArena...)

		if ca.shift.At(math.Subscript(i)) == 1 {
			masterArena = append(masterArena, math.IdentityMonomial)
		}

		if i < len(offsets) {
			offsets[i] = uint32(len(masterArena))
		}
	}

	return &general.CA{
		Arena:   masterArena,
		Offsets: offsets,
	}, nil
}
```

Файл lib/internal/automata/invertible/derivation_test.go

```
package invertible_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestDeriveGeneralCA(t *testing.T) {
	ca := invertible.NewCA()
	rng := rand.New(rand.NewSource(42))

	if err := ca.Generate(rng); err != nil {
		t.Fatalf("Failed to generate CA: %v", err)
	}

	gen, err := ca.DeriveGeneralCA()
	if err != nil {
		t.Fatalf("Failed to derive general CA: %v", err)
	}

	src := make([]byte, invertible.StateBytes)
	dstInvertible := make([]byte, invertible.StateBytes)
	dstGeneral := make([]byte, invertible.StateBytes)

	for i := range 5 {
		rng.Read(src)

		ca.Apply(dstInvertible, src)
		gen.Apply(dstGeneral, src)

		if !bytes.Equal(dstInvertible, dstGeneral) {
			t.Errorf("Inconsistency found at iteration %d:\nInvertible Apply: %x\nGeneral Apply:    %x", i, dstInvertible, dstGeneral)
		}
	}
}
```

Файл lib/internal/automata/invertible/registry.go

```
package invertible

import (
	"slices"

	"github.com/staleread/aquila/internal/automata/math"
)

const (
	FoldsCount             = StateSize / math.VectorSize
	SymbolicPolynomialSize = math.VectorSize + math.ConfusionDegree - 1
	MaxFoldMonomials       = (math.VectorSize + math.ConfusionDegree - 1) * math.VectorSize
	MaxLinearFoldMonomials = math.VectorSize * math.VectorSize
	MaxRuleMonomials       = MaxLinearFoldMonomials + MaxFoldMonomials*(FoldsCount-1)
	MaxCAMonomials         = MaxRuleMonomials * RulesCount
)

type SymbolicRegistry struct {
	arena    []math.Monomial
	offsets  [RulesCount][StateSize + 1]uint32
	invPerms [RulesCount][StateSize]uint8
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

		for k := range StateSize {
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
		e.offsets[r][StateSize] = uint32(len(e.arena))
	}
	return e
}
```

Файл lib/internal/automata/invertible/rule.go

```
package invertible

import (
	"io"

	"github.com/staleread/aquila/internal/automata/config"
	"github.com/staleread/aquila/internal/automata/math"
)

const (
	RulesCount      = config.CompositionCount + 1
	LinearFoldBytes = math.SLEBytes
	FoldBytes       = math.SLEBytes + math.ConfusionMapBytes
	RuleFoldsBytes  = LinearFoldBytes + FoldBytes*(FoldsCount-1)
)

type Rule struct {
	arena []byte
}

type Fold struct {
	sle       math.SLE
	confusion math.ConfusionMap
}

func (r *Rule) Generate(rnd io.Reader) error {
	perm := r.getPermutation()
	entropyBuf := r.arena[:math.PermutationBytes-1]

	if err := perm.Generate(rnd, entropyBuf); err != nil {
		return err
	}

	for i := range FoldsCount {
		fold := r.getFold(i)

		if err := fold.sle.Generate(rnd); err != nil {
			return err
		}

		if i == 0 {
			continue
		}

		maxSub := math.Subscript(i * math.VectorSize)
		if err := fold.confusion.Generate(rnd, maxSub, perm); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rule) Apply(s *State) {
	perm := r.getPermutation()
	var srcState State = *s

	for i := range FoldsCount {
		fold := r.getFold(i)

		in := perm.Gather(s, i)

		out := fold.sle.Eval(in)
		out ^= fold.confusion.Eval(srcState)

		perm.Scatter(s, i, out)
	}
}

func (r *Rule) Revert(s *State) {
	perm := r.getPermutation()

	for i := range FoldsCount {
		fold := r.getFold(i)

		in := perm.Gather(s, i)

		noise := fold.confusion.Eval(*s)
		out := fold.sle.Solve(in ^ noise)

		perm.Scatter(s, i, out)
	}
}

func (r *Rule) getFold(idx int) Fold {
	if idx == 0 {
		return Fold{
			sle:       math.NewSLE(r.arena[:LinearFoldBytes]),
			confusion: math.EmptyConfusionMap(),
		}
	}

	offset := LinearFoldBytes + (idx-1)*FoldBytes

	return Fold{
		sle:       math.NewSLE(r.arena[offset : offset+math.SLEBytes]),
		confusion: math.NewConfusionMap(r.arena[offset+math.SLEBytes : offset+FoldBytes]),
	}
}

func (r *Rule) getPermutation() *math.Permutation {
	const offset = RuleFoldsBytes

	view := r.arena[offset : offset+math.PermutationBytes]

	return math.NewPermutation(view)
}
```

Файл lib/internal/automata/config/composition0.go

```
//go:build comp0

package config

const CompositionCount = 0
```

Файл lib/internal/automata/config/composition1.go

```
//go:build comp1 || (!comp0 && !comp2 && !comp3)

package config

const CompositionCount = 1
```

Файл lib/internal/automata/config/composition2.go

```
//go:build comp2

package config

const CompositionCount = 2
```

Файл lib/internal/automata/config/composition3.go

```
//go:build comp3

package config

const CompositionCount = 3
```

Файл lib/internal/automata/config/config.go

```
package config

import (
	"fmt"
	"io"

	"github.com/staleread/aquila/internal/automata/math"
)

const CAConfigBytes = 4

type CAConfig struct {
	Block  byte
	Comp   byte
	Fold   byte
	Degree byte
}

var Current = CAConfig{
	Block:  byte(math.BitsetSize),
	Comp:   byte(CompositionCount),
	Fold:   byte(math.VectorSize),
	Degree: byte(math.ConfusionDegree),
}

func (c CAConfig) Write(w io.Writer) error {
	buf := [CAConfigBytes]byte{
		c.Block,
		c.Comp,
		c.Fold,
		c.Degree,
	}
	_, err := w.Write(buf[:])
	return err
}

func (c CAConfig) Check(r io.Reader) error {
	var buf [CAConfigBytes]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	block := int(buf[0])
	comp := int(buf[1])
	fold := int(buf[2])
	deg := int(buf[3])

	if block != int(c.Block) {
		panic(fmt.Sprintf("incompatible block size: got %d, want %d", block, c.Block))
	}
	if comp != int(c.Comp) {
		panic(fmt.Sprintf("incompatible composition count: got %d, want %d", comp, c.Comp))
	}
	if fold != int(c.Fold) {
		panic(fmt.Sprintf("incompatible fold size: got %d, want %d", fold, c.Fold))
	}
	if deg != int(c.Degree) {
		panic(fmt.Sprintf("incompatible confusion degree: got %d, want %d", deg, c.Degree))
	}
	return nil
}
```

Файл lib/internal/automata/general/ca.go

```
package general

import (
	"encoding/binary"
	"io"

	"github.com/staleread/aquila/internal/automata/config"
	"github.com/staleread/aquila/internal/automata/math"
)

const (
	StateSize  = math.BitsetSize
	StateBytes = math.BitsetBytes
)

type State = math.Bitset

type CA struct {
	Arena   []math.Monomial
	Offsets [StateSize - 1]uint32
}

func (ca *CA) GetPolynomial(idx int) []math.Monomial {
	start := uint32(0)
	if idx > 0 {
		start = ca.Offsets[idx-1]
	}
	end := uint32(len(ca.Arena))
	if idx < StateSize-1 {
		end = ca.Offsets[idx]
	}
	return ca.Arena[start:end]
}

func (ca *CA) GetMonomialCounts() []int {
	counts := make([]int, StateSize)
	for i := range StateSize {
		counts[i] = len(ca.GetPolynomial(i))
	}
	return counts
}

func (ca *CA) Apply(dst, src []byte) {
	var srcState State
	srcState.Read(src)
	dstState := ca.applyOnState(srcState)
	dstState.Write(dst)
}

func (ca *CA) Save(dst io.Writer) error {
	if err := config.Current.Write(dst); err != nil {
		return err
	}
	if err := binary.Write(dst, binary.LittleEndian, ca.Offsets); err != nil {
		return err
	}
	if err := binary.Write(dst, binary.LittleEndian, uint32(len(ca.Arena))); err != nil {
		return err
	}
	return binary.Write(dst, binary.LittleEndian, ca.Arena)
}

func LoadCA(src io.Reader) (*CA, error) {
	if err := config.Current.Check(src); err != nil {
		return nil, err
	}

	ca := &CA{}
	if err := binary.Read(src, binary.LittleEndian, &ca.Offsets); err != nil {
		return nil, err
	}
	var length uint32
	if err := binary.Read(src, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	ca.Arena = make([]math.Monomial, length)
	if err := binary.Read(src, binary.LittleEndian, &ca.Arena); err != nil {
		return nil, err
	}
	return ca, nil
}

func (ca *CA) applyOnState(srcState State) State {
	var dstState State
	for i := range StateSize {
		var res uint8
		for _, monom := range ca.GetPolynomial(i) {
			res ^= monom.Eval(srcState)
		}
		dstState.SetAt(math.Subscript(i), res)
	}
	return dstState
}
```

Файл lib/internal/automata/general/ca_test.go

```
package general_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestCACompilationCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping CA compilation test in short mode")
	}

	rng := rand.New(rand.NewSource(42))
	invCA := invertible.NewCA()
	if err := invCA.Generate(rng); err != nil {
		t.Fatalf("Failed to generate invertible CA: %v", err)
	}

	genCA, err := invCA.DeriveGeneralCA()
	if err != nil {
		t.Fatalf("Failed to compile CA: %v", err)
	}

	src := make([]byte, invertible.StateBytes)
	got := make([]byte, invertible.StateBytes)
	want := make([]byte, invertible.StateBytes)

	for range 5 {
		rng.Read(src)

		invCA.Apply(want, src)
		genCA.Apply(got, src)

		if !bytes.Equal(got, want) {
			t.Errorf("Compiled CA output mismatch!\ngot:  %x\nwant: %x", got, want)
		}
	}
}
```

Файл lib/internal/automata/general/export.go

```
package general

import (
	"fmt"
	"io"

	"github.com/staleread/aquila/internal/automata/math"
)

func (ca *CA) ExportToANF(w io.Writer, input State) error {
	dstState := ca.applyOnState(input)

	var subs [StateSize]uint8
	for idx := range StateSize {
		poly := ca.GetPolynomial(idx)
		bit := dstState.At(math.Subscript(idx))

		var monomialsToPrint []math.Monomial
		addOneTerm := false

		if len(poly) > 0 {
			lastMonom := poly[len(poly)-1]
			if lastMonom == math.IdentityMonomial {
				monomialsToPrint = poly[:len(poly)-1]
				if bit == 0 {
					addOneTerm = true
				}
			} else {
				monomialsToPrint = poly
				if bit == 1 {
					addOneTerm = true
				}
			}
		} else {
			if bit == 1 {
				addOneTerm = true
			}
		}

		printedAny := false
		for _, monom := range monomialsToPrint {
			if printedAny {
				if _, err := io.WriteString(w, " + "); err != nil {
					return err
				}
			}
			subsSlice := monom.Subscripts(subs[:0])
			for k, sub := range subsSlice {
				if k > 0 {
					if _, err := io.WriteString(w, "*"); err != nil {
						return err
					}
				}
				if _, err := fmt.Fprintf(w, "x%d", sub); err != nil {
					return err
				}
			}
			printedAny = true
		}

		if addOneTerm {
			if printedAny {
				if _, err := io.WriteString(w, " + "); err != nil {
					return err
				}
			}
			if _, err := io.WriteString(w, "1"); err != nil {
				return err
			}
			printedAny = true
		}

		if !printedAny {
			if _, err := io.WriteString(w, "0"); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}
```

Файл lib/asym/asym.go

```
package asym

import (
	"io"

	"github.com/staleread/aquila/internal/automata/config"
	"github.com/staleread/aquila/internal/automata/invertible"
	"github.com/staleread/aquila/internal/automata/math"
)

const (
	BlockSize    = invertible.StateBytes
	FoldSize     = math.VectorSize
	Degree       = math.ConfusionDegree
	Compositions = config.CompositionCount
)

func GenerateKeyPair(rnd io.Reader) (*PrivateKey, *PublicKey, error) {
	priv, err := GeneratePrivateKey(rnd)
	if err != nil {
		return nil, nil, err
	}

	pub, err := priv.PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return priv, pub, nil
}
```

Файл lib/asym/asym_test.go

```
package asym_test

import (
	"crypto/rand"
	"testing"

	"github.com/staleread/aquila/asym"
)

func TestEncryptDecrypt(t *testing.T) {
	priv, pub, err := asym.GenerateKeyPair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	payload := make([]byte, asym.BlockSize)
	rand.Read(payload)

	ciphertext, err := pub.Encrypt(rand.Reader, payload)
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := priv.Decrypt(rand.Reader, ciphertext, nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(plaintext) != string(payload) {
		t.Errorf("expected %x, got %x", payload, plaintext)
	}
}

func TestSignVerify(t *testing.T) {
	priv, pub, err := asym.GenerateKeyPair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	digest := make([]byte, asym.BlockSize)
	rand.Read(digest)

	sig, err := priv.Sign(rand.Reader, digest, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !pub.Verify(digest, sig) {
		t.Error("verification failed for valid signature")
	}

	// Corrupt signature
	sig[0] ^= 1
	if pub.Verify(digest, sig) {
		t.Error("verification succeeded for corrupt signature")
	}
}
```

Файл lib/asym/private.go

```
package asym

import (
	"crypto"
	"errors"
	"io"

	"github.com/staleread/aquila/internal/automata/invertible"
)

var _ crypto.Decrypter = (*PrivateKey)(nil)
var _ crypto.Signer = (*PrivateKey)(nil)

type PrivateKey struct {
	ca  *invertible.CA
	pub *PublicKey
}

func GeneratePrivateKey(rnd io.Reader) (*PrivateKey, error) {
	ca := invertible.NewCA()
	if err := ca.Generate(rnd); err != nil {
		return nil, err
	}
	return &PrivateKey{ca: ca}, nil
}

func (k *PrivateKey) PublicKey() (*PublicKey, error) {
	if k.pub == nil {
		gen, err := k.ca.DeriveGeneralCA()
		if err != nil {
			return nil, err
		}
		k.pub = &PublicKey{gen}
	}
	return k.pub, nil
}

func (k *PrivateKey) Public() crypto.PublicKey {
	if k.pub == nil {
		gen, err := k.ca.DeriveGeneralCA()
		if err != nil {
			return nil
		}
		k.pub = &PublicKey{gen}
	}
	return k.pub
}

func (k *PrivateKey) Decode(src io.Reader) error {
	ca := invertible.NewCA()
	if err := ca.Load(src); err != nil {
		return err
	}
	k.ca = ca
	return nil
}

func (k *PrivateKey) Encode(dst io.Writer) error {
	return k.ca.Save(dst)
}

func (k *PrivateKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	if len(msg)%invertible.StateBytes != 0 {
		return nil, errors.New("invalid ciphertext length")
	}

	plaintext = make([]byte, len(msg))
	for i := 0; i < len(msg); i += invertible.StateBytes {
		k.ca.Revert(plaintext[i:i+invertible.StateBytes], msg[i:i+invertible.StateBytes])
	}

	return plaintext, nil
}

func (k *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if len(digest)%invertible.StateBytes != 0 {
		return nil, errors.New("invalid digest length")
	}

	signature = make([]byte, len(digest))
	for i := 0; i < len(digest); i += invertible.StateBytes {
		k.ca.Revert(signature[i:i+invertible.StateBytes], digest[i:i+invertible.StateBytes])
	}

	return signature, nil
}
```

Файл lib/asym/public.go

```
package asym

import (
	"bytes"
	"errors"
	"io"

	"github.com/staleread/aquila/internal/automata/general"
)

type PublicKey struct {
	ca *general.CA
}

type PublicKeyInfo struct {
	BlockSizeBytes int
	BlockSizeBits  int
	FoldSize       int
	Degree         int
	Compositions   int
	MonomialCounts []int
}

func (k *PublicKey) Describe() PublicKeyInfo {
	return PublicKeyInfo{
		BlockSizeBytes: BlockSize,
		BlockSizeBits:  BlockSize * 8,
		FoldSize:       FoldSize,
		Degree:         Degree,
		Compositions:   Compositions,
		MonomialCounts: k.ca.GetMonomialCounts(),
	}
}

func (k *PublicKey) Decode(src io.Reader) error {
	ca, err := general.LoadCA(src)
	if err != nil {
		return err
	}
	k.ca = ca
	return nil
}

func (k *PublicKey) Encode(dst io.Writer) error {
	return k.ca.Save(dst)
}

func (k *PublicKey) Encrypt(rand io.Reader, msg []byte) (ciphertext []byte, err error) {
	if len(msg)%general.StateBytes != 0 {
		return nil, errors.New("invalid plaintext length")
	}

	ciphertext = make([]byte, len(msg))
	for i := 0; i < len(msg); i += general.StateBytes {
		k.ca.Apply(ciphertext[i:i+general.StateBytes], msg[i:i+general.StateBytes])
	}

	return ciphertext, nil
}

func (k *PublicKey) Verify(digest []byte, signature []byte) bool {
	if len(signature) != len(digest) || len(signature)%general.StateBytes != 0 {
		return false
	}

	temp := make([]byte, len(signature))
	for i := 0; i < len(signature); i += general.StateBytes {
		k.ca.Apply(temp[i:i+general.StateBytes], signature[i:i+general.StateBytes])
	}

	return bytes.Equal(temp, digest)
}

func (k *PublicKey) ExportToANF(w io.Writer, input []byte) error {
	if len(input) != BlockSize {
		return errors.New("invalid input block length")
	}
	var s general.State
	s.Read(input)
	return k.ca.ExportToANF(w, s)
}
```

Файл lib/sym/cipher.go

```
package sym

import (
	"crypto/cipher"
	"io"

	"github.com/staleread/aquila/internal/automata/invertible"
)

var _ cipher.Block = (*AquilaBlock)(nil)

type AquilaBlock struct {
	ca *invertible.CA
}

func New(rnd io.Reader) (*AquilaBlock, error) {
	ca := invertible.NewCA()

	if err := ca.Generate(rnd); err != nil {
		return nil, err
	}

	return &AquilaBlock{ca}, nil
}

func Decode(src io.Reader) (*AquilaBlock, error) {
	ca := invertible.NewCA()

	if err := ca.Load(src); err != nil {
		return nil, err
	}
	return &AquilaBlock{ca}, nil
}

func (b *AquilaBlock) Encode(dst io.Writer) error { return b.ca.Save(dst) }

func (b *AquilaBlock) BlockSize() int          { return invertible.StateBytes }
func (b *AquilaBlock) Encrypt(dst, src []byte) { b.ca.Apply(dst, src) }
func (b *AquilaBlock) Decrypt(dst, src []byte) { b.ca.Revert(dst, src) }
```

Файл cli/analyze.go

```
package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"math/bits"
	"os"
	"path/filepath"
	"strconv"

	"github.com/staleread/aquila/asym"
)

func runAnalyze(folder string, rndSamples int, correlation bool) error {
	priv, err := asym.GeneratePrivateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	blockSizeBytes := asym.BlockSize
	blockSizeBits := blockSizeBytes * 8

	changeCounts := make([]int, blockSizeBits)

	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", folder, err)
	}

	avFile, err := os.Create(filepath.Join(folder, "avalanche.csv"))
	if err != nil {
		return fmt.Errorf("failed to create avalanche.csv: %w", err)
	}
	defer avFile.Close()

	writer := csv.NewWriter(avFile)
	defer writer.Flush()

	if err := writer.Write([]string{"bit_index", "zeros_bg", "ones_bg", "random_bg_avg"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	hammingDistance := func(a, b []byte) int {
		dist := 0
		for i := range a {
			dist += bits.OnesCount8(a[i] ^ b[i])
		}
		return dist
	}

	zeroBaseSrc := make([]byte, blockSizeBytes)
	zeroBaseDst, err := priv.Sign(nil, zeroBaseSrc, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt zero base: %w", err)
	}

	oneBaseSrc := make([]byte, blockSizeBytes)
	for j := range oneBaseSrc {
		oneBaseSrc[j] = 0xFF
	}
	oneBaseDst, err := priv.Sign(nil, oneBaseSrc, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt one base: %w", err)
	}

	randSamplesSrc := make([][]byte, rndSamples)
	randSamplesDst := make([][]byte, rndSamples)
	for s := range rndSamples {
		randSamplesSrc[s] = make([]byte, blockSizeBytes)
		if _, err := io.ReadFull(rand.Reader, randSamplesSrc[s]); err != nil {
			return fmt.Errorf("failed to read random bytes: %w", err)
		}
		randSamplesDst[s], err = priv.Sign(nil, randSamplesSrc[s], nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt random sample: %w", err)
		}
	}

	srcBuf := make([]byte, blockSizeBytes)
	var dstBuf []byte

	for i := range blockSizeBits {
		// --- Zeros BG ---
		for j := range srcBuf {
			srcBuf[j] = 0
		}
		srcBuf[i/8] = 1 << (i % 8)
		dstBuf, err = priv.Sign(nil, srcBuf, nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt zeros bg: %w", err)
		}
		zerosBgVal := hammingDistance(dstBuf, zeroBaseDst)

		// --- Ones Background ---
		for j := range srcBuf {
			srcBuf[j] = 0xFF
		}
		srcBuf[i/8] &^= (1 << (i % 8))
		dstBuf, err = priv.Sign(nil, srcBuf, nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt ones bg: %w", err)
		}
		onesBgVal := hammingDistance(dstBuf, oneBaseDst)

		// --- Random Background Avg ---
		totalRandDiff := 0
		for s := range rndSamples {
			copy(srcBuf, randSamplesSrc[s])
			srcBuf[i/8] ^= (1 << (i % 8))
			dstBuf, err = priv.Sign(nil, srcBuf, nil)
			if err != nil {
				return fmt.Errorf("failed to encrypt random background: %w", err)
			}
			totalRandDiff += hammingDistance(dstBuf, randSamplesDst[s])

			for o := range blockSizeBits {
				if ((randSamplesDst[s][o/8] ^ dstBuf[o/8]) & (1 << (o % 8))) != 0 {
					changeCounts[o]++
				}
			}
		}
		randBgAvgVal := float64(totalRandDiff) / float64(rndSamples)

		row := []string{
			strconv.Itoa(i),
			strconv.Itoa(zerosBgVal),
			strconv.Itoa(onesBgVal),
			fmt.Sprintf("%.4f", randBgAvgVal),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row %d to avalanche.csv: %w", i, err)
		}
	}

	if correlation {
		pub, err := priv.PublicKey()
		if err != nil {
			return fmt.Errorf("failed to derive public key: %w", err)
		}
		desc := pub.Describe()

		corrFile, err := os.Create(filepath.Join(folder, "monom-correlation.csv"))
		if err != nil {
			return fmt.Errorf("failed to create monom-correlation.csv: %w", err)
		}
		defer corrFile.Close()

		corrWriter := csv.NewWriter(corrFile)
		defer corrWriter.Flush()

		if err := corrWriter.Write([]string{"output_bit", "monomial_count", "avalanche_prob"}); err != nil {
			return fmt.Errorf("failed to write CSV header to monom-correlation.csv: %w", err)
		}

		for o := range blockSizeBits {
			prob := float64(changeCounts[o]) / float64(blockSizeBits*rndSamples)
			row := []string{
				strconv.Itoa(o),
				strconv.Itoa(desc.MonomialCounts[o]),
				fmt.Sprintf("%.6f", prob),
			}
			if err := corrWriter.Write(row); err != nil {
				return fmt.Errorf("failed to write row %d to monom-correlation.csv: %w", o, err)
			}
		}
	}

	return nil
}
```

Файл cli/anf.go

```
package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func exportANF(keyPath string) error {
	keyF, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer keyF.Close()

	pub := &asym.PublicKey{}
	if err := pub.Decode(keyF); err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	randInput := make([]byte, asym.BlockSize)
	if _, err := rand.Read(randInput); err != nil {
		return fmt.Errorf("failed to generate random input block: %w", err)
	}

	writer := bufio.NewWriter(os.Stdout)
	if err := pub.ExportToANF(writer, randInput); err != nil {
		return fmt.Errorf("failed to export ANF: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}
```

Файл cli/build.go

```
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runBuild(configID string) error {
	var block, comp, fold, deg int
	n, err := fmt.Sscanf(configID, "%dc%df%dd%d", &block, &comp, &fold, &deg)
	if err != nil || n != 4 {
		return fmt.Errorf("invalid config ID format (must be <block>c<comp>f<fold>d<deg>): %w", err)
	}

	tags := fmt.Sprintf("block%d comp%d fold%d deg%d", block, comp, fold, deg)
	fmt.Printf("Building aquila-cli with tags: %s\n", tags)

	cmd := exec.Command("go", "build", "-tags", tags)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go build: %w", err)
	}

	fmt.Println("Build successful!")
	return nil
}
```

Файл cli/config.go

```
package main

import (
	"fmt"

	"github.com/staleread/aquila/asym"
)

func showConfig() {
	fmt.Printf("Block size: %d\n", asym.BlockSize*8)
	fmt.Printf("Compositions: %d\n", asym.Compositions)
	fmt.Printf("Fold size: %d\n", asym.FoldSize)
	fmt.Printf("Confusion Degree: %d\n", asym.Degree)
}

func getConfigID() string {
	return fmt.Sprintf("%dc%df%dd%d", asym.BlockSize*8, asym.Compositions, asym.FoldSize, asym.Degree)
}
```

Файл cli/decrypt.go

```
package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func decryptFile(inputPath, outputPath, keyPath string) error {
	keyF, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open private key file: %w", err)
	}
	defer keyF.Close()

	priv := &asym.PrivateKey{}
	if err := priv.Decode(keyF); err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if len(ciphertext)%asym.BlockSize != 0 {
		return fmt.Errorf("ciphertext length is not a multiple of the block size")
	}

	paddedPlaintext, err := priv.Decrypt(rand.Reader, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	plaintext, err := pkcs7Unpad(paddedPlaintext, asym.BlockSize)
	if err != nil {
		return fmt.Errorf("failed to unpad decrypted data: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully decrypted %s to %s\n", inputPath, outputPath)
	return nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}

	unpadding := int(data[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	padtext := data[length-unpadding:]
	for _, b := range padtext {
		if int(b) != unpadding {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:(length - unpadding)], nil
}
```

Файл cli/encrypt.go

```
package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func encryptFile(inputPath, outputPath, keyPath string) error {
	keyF, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer keyF.Close()

	pub := &asym.PublicKey{}
	if err := pub.Decode(keyF); err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	paddedData := pkcs7Pad(inputData, asym.BlockSize)

	ciphertext, err := pub.Encrypt(rand.Reader, paddedData)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully encrypted %s to %s\n", inputPath, outputPath)
	return nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}
```

Файл cli/generate.go

```
package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/staleread/aquila/asym"
)

func generateKeyPair(name string) error {
	cfgID := getConfigID()
	privFile := fmt.Sprintf("id_aquila%s", cfgID)
	pubFile := fmt.Sprintf("id_aquila%s.pub", cfgID)

	if name != "" {
		privFile = fmt.Sprintf("id_aquila%s_%s", cfgID, name)
		pubFile = fmt.Sprintf("id_aquila%s_%s.pub", cfgID, name)
	}

	priv, pub, err := asym.GenerateKeyPair(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	privF, err := os.OpenFile(privFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open private key file: %w", err)
	}
	defer privF.Close()

	if err := priv.Encode(privF); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	pubF, err := os.OpenFile(pubFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer pubF.Close()

	if err := pub.Encode(pubF); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	fmt.Printf("Successfully generated key pair:\n  Private key: %s\n  Public key:  %s\n", privFile, pubFile)
	return nil
}
```

Файл cli/main.go

```
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Gen struct {
		Name string `short:"n" help:"Optional name for the key pair."`
	} `cmd:"" help:"Generate a new key pair."`

	Enc struct {
		Input  string `short:"i" required:"" type:"existingfile" help:"Path to the input file."`
		Output string `short:"o" required:"" type:"path" help:"Path to the output file."`
		Key    string `short:"k" required:"" type:"existingfile" help:"Path to the public key file."`
	} `cmd:"" help:"Encrypt a file."`

	Dec struct {
		Input  string `short:"i" required:"" type:"existingfile" help:"Path to the input file."`
		Output string `short:"o" required:"" type:"path" help:"Path to the output file."`
		Key    string `short:"k" required:"" type:"existingfile" help:"Path to the private key file."`
	} `cmd:"" help:"Decrypt a file."`

	Config struct{} `cmd:"" help:"Show cipher configuration."`

	Build struct {
		ConfigID string `arg:"" help:"Configuration ID in format <block>c<comp>f<fold>d<deg>."`
	} `cmd:"" help:"Build the CLI with the specified cipher configuration."`

	Analyze struct {
		Folder      string `short:"o" required:"" type:"path" help:"Path to the output folder."`
		RndSamples  int    `long:"rnd-samples" required:"" help:"Number of random samples to analyze."`
		Correlation bool   `long:"correlation" help:"Generate monom-correlation.csv using public key (caution: deriving public key for large configurations like 96-bit key with 2 compositions can result in ~3 GB public key)."`
	} `cmd:"" name:"analyze" help:"Run cipher analysis experiments."`

	Anf struct {
		Key string `short:"k" required:"" type:"existingfile" help:"Path to the public key file."`
	} `cmd:"" help:"Export public key equations in Algebraic Normal Form (ANF)."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "gen":
		if err := generateKeyPair(CLI.Gen.Name); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating keys: %v\n", err)
			os.Exit(1)
		}
	case "enc":
		if err := encryptFile(CLI.Enc.Input, CLI.Enc.Output, CLI.Enc.Key); err != nil {
			fmt.Fprintf(os.Stderr, "Error encrypting file: %v\n", err)
			os.Exit(1)
		}
	case "dec":
		if err := decryptFile(CLI.Dec.Input, CLI.Dec.Output, CLI.Dec.Key); err != nil {
			fmt.Fprintf(os.Stderr, "Error decrypting file: %v\n", err)
			os.Exit(1)
		}
	case "config":
		showConfig()
	case "build <config-id>":
		if err := runBuild(CLI.Build.ConfigID); err != nil {
			fmt.Fprintf(os.Stderr, "Error building CLI: %v\n", err)
			os.Exit(1)
		}
	case "analyze":
		if err := runAnalyze(CLI.Analyze.Folder, CLI.Analyze.RndSamples, CLI.Analyze.Correlation); err != nil {
			fmt.Fprintf(os.Stderr, "Error running analysis: %v\n", err)
			os.Exit(1)
		}
	case "anf":
		if err := exportANF(CLI.Anf.Key); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting ANF: %v\n", err)
			os.Exit(1)
		}
	}
}
```
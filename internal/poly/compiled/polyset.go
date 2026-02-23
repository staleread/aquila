package compiled

import (
	"github.com/staleread/aquila/internal/field"
	"github.com/staleread/aquila/internal/poly"
)

type Polyset struct {
	Subscripts   []poly.Subscript
	PolyCount    int
	PolyOffsets  []int
	MonomOffsets []int
}

func EmptyPolyset() *Polyset {
	return &Polyset{
		Subscripts:   nil,
		PolyCount:    0,
		PolyOffsets:  nil,
		MonomOffsets: nil,
	}
}

func RandPolyset(n, degree int, maxSub poly.Subscript) *Polyset {
	if poly.Subscript(degree) > maxSub {
		panic("Polynomial degree exceeds subscript range")
	}

	totalSubs := degree * (degree + 1) / 2 * n

	subs := make([]poly.Subscript, totalSubs)
	polyOffsets := make([]int, 1)
	monomOffsets := make([]int, 1)

	rands := poly.RandSubscripts(totalSubs)

	var monomOffset int
	var polyOffset int

	for range n {
		for deg := poly.Subscript(degree); deg > 0; deg-- {
			added := make(map[poly.Subscript]struct{}, deg)
			randUp := maxSub - deg

			for k := range deg {
				sub := rands[monomOffset] % (randUp + k + 1)

				if _, ok := added[sub]; ok {
					sub = randUp + k
				}
				added[sub] = struct{}{}
				subs[monomOffset] = sub

				monomOffset++
			}
			monomOffsets = append(monomOffsets, monomOffset)
			polyOffset++
		}
		polyOffsets = append(polyOffsets, polyOffset)
	}

	// FOR DEV
	if monomOffset != totalSubs {
		panic("layout mismatch")
	}

	return &Polyset{
		Subscripts:   subs,
		PolyCount:    n,
		PolyOffsets:  polyOffsets,
		MonomOffsets: monomOffsets,
	}
}

func (set *Polyset) Eval(dst, src []field.Element) {
	if set.PolyCount == 0 {
		copy(dst, src)
		return
	}
	if set.PolyCount != len(dst) {
		panic("Element array does not match polyset size")
	}

	for i := range set.PolyCount {
		var sum field.Element

		pStart := set.PolyOffsets[i]
		pEnd := set.PolyOffsets[i+1]

		for j := pStart; j < pEnd; j++ {
			prod := field.Element(1)

			mStart := set.MonomOffsets[j]
			mEnd := set.MonomOffsets[j+1]

			for k := mStart; k < mEnd; k++ {
				val := src[set.Subscripts[k]]
				prod = field.Mul(prod, val)
			}
			sum = field.Add(sum, prod)
		}
		dst[i] = sum
	}
}

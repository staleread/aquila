package main

import (
	"fmt"
	"strings"
)

type snle struct {
	data []poly
}

func emptySnle(n int) *snle {
	return &snle{make([]poly, n)}
}

func randSnle(n, deg int, maxIdx idx) *snle {
	snle := emptySnle(n)

	for i := range n {
		snle.data[i] = randPoly(deg, maxIdx)
	}
	return snle
}

func (s *snle) eval(x, b vec) {
	for i, p := range s.data {
		var sum elem = 0

		for _, m := range p {
			var prod elem = 1

			for _, idx := range m {
				prod = mul(prod, x[idx])
			}
			sum = add(sum, prod)
		}
		b[i] = sum
	}
}

func (s *snle) String() string {
	sb := strings.Builder{}

	for i, p := range s.data {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		for j, m := range p {
			if j > 0 {
				sb.WriteString(" + ")
			}
			for k, id := range m {
				if k > 0 {
					sb.WriteRune('*')
				}
				fmt.Fprintf(&sb, "x%d", id+1)
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

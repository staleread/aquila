package field

import (
	"fmt"
	"strings"
)

type Polyset []Polynomial

func (set Polyset) Eval(dst, src []Element) {
	for i := range len(dst) {
		dst[i] = set[i].Eval(src)
	}
}

func (set Polyset) String() string {
	sb := strings.Builder{}

	for i, p := range set {
		fmt.Fprintf(&sb, "y%d = ", i+1)

		firstMonomial := true
		for m := range p.Monomials() {
			if !firstMonomial {
				sb.WriteString(" + ")
			}
			firstMonomial = false

			firstSubscript := true
			for s := range m.Subscripts() {
				if !firstSubscript {
					sb.WriteRune('*')
				}
				firstSubscript = false

				fmt.Fprintf(&sb, "x%d", s+1)
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

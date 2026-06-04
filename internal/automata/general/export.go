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

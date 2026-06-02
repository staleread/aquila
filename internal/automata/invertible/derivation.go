package invertible

import (
	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/math"
	"github.com/staleread/aquila/internal/automata/state"
)

func (ca *CA) DeriveGeneralCA() (*general.CA, error) {
	polynomials := CompileRegistry(ca)
	masterArena := make([]math.Monomial, 0, InitialArenaCapacity*state.StateSize)
	var offsets [state.StateSize - 1]uint32

	stateArena := make([]math.Monomial, 0, InitialArenaCapacity)
	prodSrcArena := make([]math.Monomial, 0, InitialArenaCapacity)
	prodDstArena := make([]math.Monomial, 0, InitialArenaCapacity)
	sumArena := make([]math.Monomial, 0, InitialArenaCapacity)
	sumScratchArena := make([]math.Monomial, 0, InitialArenaCapacity)

	for i := range state.StateSize {
		stateArena = stateArena[:0]
		stateArena = append(stateArena, polynomials.GetPolynomial(RulesCount-1, i)...)

		var subs [state.StateSize]uint8
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

		if ca.shift.At(state.Subscript(i)) == 1 {
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

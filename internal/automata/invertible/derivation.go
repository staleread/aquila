package invertible

import (
	stdmath "math"

	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/math"
)

var InitialArenaCapacity = calculateMaxMonomiasPerBit() / 2

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

func calculateMaxMonomiasPerBit() int {
	L := math.VectorSize / 2
	d0 := math.ConfusionDegree
	r := RulesCount - 1

	M := L + d0 - 1
	for i := 1; i <= r; i++ {
		if M <= 1 {
			M = M*M + L*M
			continue
		}
		mFloat := float64(M)
		p1 := stdmath.Pow(mFloat, float64(d0+1))
		p2 := stdmath.Pow(mFloat, 2.0)
		M = int((p1-p2)/float64(M-1)) + L*M
	}
	return M
}

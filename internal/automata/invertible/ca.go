package invertible

import (
	"fmt"
	"io"

	"github.com/staleread/aquila/internal/automata/core"
	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/math"
)

type CA struct {
	arena []byte
}

func NewCA() *CA {
	return &CA{
		arena: make([]byte, CABytes),
	}
}

func (ca *CA) Load(src io.Reader) error {
	_, err := io.ReadFull(src, ca.arena)
	return err
}

func (ca *CA) Save(dst io.Writer) error {
	_, err := dst.Write(ca.arena)
	return err
}

func (ca *CA) Generate(rnd io.Reader) error {
	for i := range RulesCount {
		rule := ca.getRule(i)

		if err := rule.Generate(rnd); err != nil {
			return fmt.Errorf("failed to generate rule %d: %w", i, err)
		}
	}
	return nil
}

func (ca *CA) Apply(dst, src []byte) {
	block := core.LoadBlock(src)

	for i := range RulesCount {
		rule := ca.getRule(i)
		rule.Apply(block)
	}

	block.WriteTo(dst)
}

func (ca *CA) Revert(dst, src []byte) {
	block := core.LoadBlock(src)

	for i := RulesCount - 1; i >= 0; i-- {
		rule := ca.getRule(i)
		rule.Revert(block)
	}

	block.WriteTo(dst)
}

func (ca *CA) getRule(idx int) Rule {
	offset := idx * RuleBytes

	return Rule{
		arena: ca.arena[offset : offset+RuleBytes],
	}
}

func (ca *CA) DeriveGeneralCA() (*general.CA, error) {
	polynomials := CompileRegistry(ca)
	masterArena := make([]math.Monomial, 0, EstimatedDensePolynomialSize)
	var offsets [core.BlockSize - 1]uint32

	stateArena := make([]math.Monomial, 0, EstimatedDensePolynomialSize)
	prodSrcArena := make([]math.Monomial, 0, EstimatedDensePolynomialSize)
	prodDstArena := make([]math.Monomial, 0, EstimatedDensePolynomialSize)
	sumArena := make([]math.Monomial, 0, EstimatedDensePolynomialSize)
	sumScratchArena := make([]math.Monomial, 0, EstimatedDensePolynomialSize)

	for i := range core.BlockSize {
		if i%8 == 0 {
			fmt.Printf("Compiling bit %d/%d (masterArena: %d monomials)\n", i, core.BlockSize, len(masterArena))
		}
		stateArena = stateArena[:0]
		stateArena = append(stateArena, polynomials.GetPolynomial(RulesCount-1, i)...)

		for j := RulesCount - 2; j >= 0; j-- {
			for _, monom := range stateArena {
				degree := monom.Degree()

				if degree == 0 {
					continue
				}

				firstSubscript := monom.FirstSubscript()
				firstPoly := polynomials.GetPolynomial(j, firstSubscript)

				if degree == 1 {
					sumScratchArena = sumScratchArena[:0]
					sumScratchArena = math.AddPolynomials(sumScratchArena, sumArena, firstPoly)

					sumArena, sumScratchArena = sumScratchArena, sumArena
					continue
				}

				prodSrcArena = prodSrcArena[:0]
				prodSrcArena = append(prodSrcArena, firstPoly...)

				currentSub := firstSubscript

				for {
					currentSub = monom.NextSubscript(currentSub + 1)

					if currentSub == core.BlockSize {
						break
					}

					nextPoly := polynomials.GetPolynomial(j, currentSub)

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

		if i < len(offsets) {
			offsets[i] = uint32(len(masterArena))
		}
	}

	return &general.CA{
		Arena:   masterArena,
		Offsets: offsets,
	}, nil
}

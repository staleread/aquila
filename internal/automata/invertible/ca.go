package invertible

import (
	"fmt"
	"io"

	"github.com/staleread/aquila/internal/automata/core"
	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/math"
)

const InitialArenaCapacity = 7_929

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
	var block core.Block
	block.Read(src)

	for i := range RulesCount {
		rule := ca.getRule(i)
		rule.Apply(&block)
	}

	block.Write(dst)
}

func (ca *CA) Revert(dst, src []byte) {
	var block core.Block
	block.Read(src)

	for i := RulesCount - 1; i >= 0; i-- {
		rule := ca.getRule(i)
		rule.Revert(&block)
	}

	block.Write(dst)
}

func (ca *CA) getRule(idx int) Rule {
	offset := idx * RuleBytes

	return Rule{
		arena: ca.arena[offset : offset+RuleBytes],
	}
}

func (ca *CA) DeriveGeneralCA() (*general.CA, error) {
	polynomials := CompileRegistry(ca)
	masterArena := make([]math.Monomial, 0, InitialArenaCapacity*core.BlockSize)
	var offsets [core.BlockSize - 1]uint32

	stateArena := make([]math.Monomial, 0, InitialArenaCapacity)
	prodSrcArena := make([]math.Monomial, 0, InitialArenaCapacity)
	prodDstArena := make([]math.Monomial, 0, InitialArenaCapacity)
	sumArena := make([]math.Monomial, 0, InitialArenaCapacity)
	sumScratchArena := make([]math.Monomial, 0, InitialArenaCapacity)

	for i := range core.BlockSize {
		stateArena = stateArena[:0]
		stateArena = append(stateArena, polynomials.GetPolynomial(RulesCount-1, i)...)

		var subs [core.BlockSize]uint8
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

		if i < len(offsets) {
			offsets[i] = uint32(len(masterArena))
		}
	}

	return &general.CA{
		Arena:   masterArena,
		Offsets: offsets,
	}, nil
}

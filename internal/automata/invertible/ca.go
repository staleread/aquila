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

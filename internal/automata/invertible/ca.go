package invertible

import (
	"fmt"
	"io"

	"github.com/staleread/aquila/internal/automata"
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

func (ca *CA) Generate(rnd io.Reader) error {
	for i := range RulesCount {
		rule := ca.getRule(i)

		if err := rule.Generate(rnd); err != nil {
			return fmt.Errorf("failed to generate rule %d: %w", i, err)
		}
	}
	return nil
}

func (ca *CA) Apply(block *automata.Block) {
	for i := range RulesCount {
		rule := ca.getRule(i)
		rule.Apply(block)
	}
}

func (ca *CA) Revert(block *automata.Block) {
	for i := RulesCount - 1; i >= 0; i-- {
		rule := ca.getRule(i)
		rule.Revert(block)
	}
}

func (ca *CA) getRule(idx int) Rule {
	offset := idx * RuleBytes

	return Rule{
		arena: ca.arena[offset : offset+RuleBytes],
	}
}

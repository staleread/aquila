package config

import (
	"fmt"
	"io"

	"github.com/staleread/aquila/internal/automata/math"
)

const CAConfigBytes = 4

type CAConfig struct {
	Block  byte
	Comp   byte
	Fold   byte
	Degree byte
}

var Current = CAConfig{
	Block:  byte(math.BitsetSize),
	Comp:   byte(CompositionCount),
	Fold:   byte(math.VectorSize),
	Degree: byte(math.ConfusionDegree),
}

func (c CAConfig) Write(w io.Writer) error {
	buf := [CAConfigBytes]byte{
		c.Block,
		c.Comp,
		c.Fold,
		c.Degree,
	}
	_, err := w.Write(buf[:])
	return err
}

func (c CAConfig) Check(r io.Reader) error {
	var buf [CAConfigBytes]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	block := int(buf[0])
	comp := int(buf[1])
	fold := int(buf[2])
	deg := int(buf[3])

	if block != int(c.Block) {
		panic(fmt.Sprintf("incompatible block size: got %d, want %d", block, c.Block))
	}
	if comp != int(c.Comp) {
		panic(fmt.Sprintf("incompatible composition count: got %d, want %d", comp, c.Comp))
	}
	if fold != int(c.Fold) {
		panic(fmt.Sprintf("incompatible fold size: got %d, want %d", fold, c.Fold))
	}
	if deg != int(c.Degree) {
		panic(fmt.Sprintf("incompatible confusion degree: got %d, want %d", deg, c.Degree))
	}
	return nil
}

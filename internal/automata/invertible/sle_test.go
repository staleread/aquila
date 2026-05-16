package invertible_test

import (
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestSLEEvalSolve(t *testing.T) {
	arena := make([]byte, invertible.SLEBytes)
	sle := invertible.NewSLE(arena)

	rng := rand.New(rand.NewSource(42))

	if err := sle.Generate(rng); err != nil {
		t.Fatalf("Failed to generate SLE: %v", err)
	}

	for range 100 {
		x := invertible.Vector(rng.Uint32())

		y := sle.Eval(x)
		xPrime := sle.Solve(y)

		if x != xPrime {
			t.Errorf("Eval/Solve mismatch: x=%04x, y=%04x, x'=%04x", x, y, xPrime)
		}
	}
}

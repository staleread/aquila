package invertible

// Test General() composition using the README example.
//
// Rule 1 (5 variables, folds: {x1,x2,x3}, {x4,x5}):
//
//	x1' = x2
//	x2' = x3
//	x3' = x1
//	x4' = x5 + x1*x2
//	x5' = x4      + x2*x3
//
// Rule 2 (folds: {x4',x5'}, {x1',x2',x3'}):
//
//	y1 = x4'
//	y2 = x5'
//	y3 = x1'         + x4'*x5'
//	y4 = x2'         + x1'*x4'
//	y5 = x3' + x1'
//
// Expected composed result (substitute Rule 1 into Rule 2):
//
//	y1 = x1*x2 + x5
//	y2 = x2*x3 + x4
//	y3 = x1*x2*x3 + x1*x2*x4 + x2*x3*x5 + x4*x5 + x2
//	y4 = x1*x2 + x2*x5 + x3
//	y5 = x1 + x2
import (
	"testing"

	"github.com/staleread/aquila/internal/canonical"
	"github.com/staleread/aquila/internal/field"
)

// buildPolyset constructs a Polyset from a slice of monomial groups.
// Each element is a list of subscript lists; each subscript list is one monomial.
// Subscripts are 0-indexed.
func buildPolyset(arena *canonical.AllocationArena, monomsPerPoly [][]field.Subscript) field.Polyset {
	ps := make(field.Polyset, len(monomsPerPoly))
	for i, subs := range monomsPerPoly {
		poly := make(field.Polynomial)
		poly.ToggleMonomial(arena.CreateMonomial(subs...))
		ps[i] = poly
	}
	return ps
}

// buildPolysetMulti constructs a Polyset where each polynomial may have multiple monomials.
func buildPolysetMulti(arena *canonical.AllocationArena, monomsPerPoly [][][]field.Subscript) field.Polyset {
	ps := make(field.Polyset, len(monomsPerPoly))
	for i, monoms := range monomsPerPoly {
		poly := make(field.Polynomial)
		for _, subs := range monoms {
			poly.ToggleMonomial(arena.CreateMonomial(subs...))
		}
		ps[i] = poly
	}
	return ps
}

// evalBool evaluates a polyset at a boolean input vector (GF(2)).
func evalBool(ps field.Polyset, xs []field.Element) []field.Element {
	out := make([]field.Element, len(ps))
	ps.Eval(out, xs)
	return out
}

// TestCompositionREADMEExample verifies General()-style composition using
// the 5-variable, 2-rule example from the README.
func TestCompositionREADMEExample(t *testing.T) {
	arena := canonical.NewAllocationArena()

	// Subscript aliases (0-indexed: x0..x4 = x1..x5 in README notation)
	x0, x1, x2, x3, x4 := field.Subscript(0), field.Subscript(1), field.Subscript(2), field.Subscript(3), field.Subscript(4)

	// Rule 1 as polyset (maps x -> x'):
	//   out[0] = x1        (x2 in README, 0-indexed = x1)
	//   out[1] = x2
	//   out[2] = x0
	//   out[3] = x4 + x0*x1
	//   out[4] = x3 + x1*x2
	rule1 := buildPolysetMulti(arena, [][][]field.Subscript{
		{{x1}},           // y0 = x1
		{{x2}},           // y1 = x2
		{{x0}},           // y2 = x0
		{{x4}, {x0, x1}}, // y3 = x4 + x0*x1
		{{x3}, {x1, x2}}, // y4 = x3 + x1*x2
	})

	// Rule 2 as polyset (maps x' -> y):
	//   out[0] = x3         (x4' in README = index 3)
	//   out[1] = x4
	//   out[2] = x0 + x3*x4
	//   out[3] = x1 + x0*x3
	//   out[4] = x2 + x0
	rule2 := buildPolysetMulti(arena, [][][]field.Subscript{
		{{x3}},           // y0 = x3
		{{x4}},           // y1 = x4
		{{x0}, {x3, x4}}, // y2 = x0 + x3*x4
		{{x1}, {x0, x3}}, // y3 = x1 + x0*x3
		{{x2}, {x0}},     // y4 = x2 + x0
	})

	// Expected composed result rule2(rule1(x)):
	//   y0 = x0*x1 + x4
	//   y1 = x1*x2 + x3
	//   y2 = x1*x2*x3 + ... (see README)
	//   y3 = x0*x1 + x1*x4 + x2
	//   y4 = x0 + x1
	expected := buildPolysetMulti(arena, [][][]field.Subscript{
		{{x0, x1}, {x4}}, // y0 = x0*x1 + x4
		{{x1, x2}, {x3}}, // y1 = x1*x2 + x3
		{{x0, x1, x2}, {x0, x1, x3}, {x1, x2, x4}, {x3, x4}, {x1}}, // y2
		{{x0, x1}, {x1, x4}, {x2}},                                 // y3 = x0*x1 + x1*x4 + x2
		{{x0}, {x1}},                                               // y4 = x0 + x1
	})

	// Compose: result = rule2(rule1(x)) by substituting rule1 into rule2.
	cmpArena := canonical.NewComputationArena()
	prodCache := make(map[*field.Monomial]field.Polynomial)

	result := make(field.Polyset, len(rule2))
	for j, p := range rule2 {
		sum := make(field.Polynomial)

		for m := range p.Monomials() {
			if cached, hit := prodCache[m]; hit {
				sum.AddTo(cached)
				continue
			}

			prod := make(field.Polynomial)
			first := true
			for s := range m.Subscripts() {
				if first {
					prod.AddTo(rule1[s])
					first = false
					continue
				}
				cmpArena.MulPolynomialBy(prod, rule1[s])
			}
			sum.AddTo(prod)
			prodCache[m] = prod
		}
		result[j] = sum
	}

	// Print what result[3] actually contains
	t.Logf("result[3] = %v", result.String())
	t.Logf("expected[3] via Polyset.String = %v", expected.String())

	// Check result matches expected on all 32 inputs.
	for bits := range 32 {
		xs := make([]field.Element, 5)
		for i := range 5 {
			xs[i] = field.Element((bits >> i) & 1)
		}

		got := evalBool(result, xs)
		want := evalBool(expected, xs)

		// Also verify via chained single-step eval.
		mid := evalBool(rule1, xs)
		chain := evalBool(rule2, mid)

		for i := range 5 {
			if got[i] != want[i] {
				t.Errorf("input=%v output[%d]: got %d, want %d (via chain: %d)",
					xs, i, got[i], want[i], chain[i])
			}
			if got[i] != chain[i] {
				t.Errorf("input=%v output[%d]: composed=%d but chained rule eval=%d",
					xs, i, got[i], chain[i])
			}
		}
	}
}

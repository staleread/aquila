= Implementation correspondence <sec-implementation>

The unit-LU generation, evaluation, and solution are implemented in
`internal/automata/math/sle.go`. The degree-$2$ through degree-$d$ monomial
terms and the predecessor-coordinate restriction are implemented in
`internal/automata/math/confusion.go`. Local-to-global coordinate transport is
implemented by `Gather` and `Scatter` in `internal/automata/math/permutation.go`.
Forward and inverse fold order are implemented in
`internal/automata/invertible/rule.go`; Rule composition and the final affine
shift are implemented in `internal/automata/invertible/ca.go`.

// Each correspondence above becomes a claim only once it is a hypothesis with
// a passing test (see hypotheses.yaml and AGENTS.md).

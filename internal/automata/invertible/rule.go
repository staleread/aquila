package invertible

import (
	"io"

	"github.com/staleread/aquila/internal/automata/config"
	"github.com/staleread/aquila/internal/automata/math"
)

const (
	RulesCount      = config.CompositionCount + 1
	LinearFoldBytes = math.SLEBytes
	FoldBytes       = math.SLEBytes + math.ConfusionMapBytes
	RuleFoldsBytes  = LinearFoldBytes + FoldBytes*(FoldsCount-1)
)

type Rule struct {
	arena []byte
}

type Fold struct {
	sle       math.SLE
	confusion math.ConfusionMap
}

func (r *Rule) Generate(rnd io.Reader) error {
	perm := r.getPermutation()
	entropyBuf := r.arena[:math.PermutationBytes-1]

	if err := perm.Generate(rnd, entropyBuf); err != nil {
		return err
	}

	for i := range FoldsCount {
		fold := r.getFold(i)

		if err := fold.sle.Generate(rnd); err != nil {
			return err
		}

		if i == 0 {
			continue
		}

		maxSub := math.Subscript(i * math.VectorSize)
		if err := fold.confusion.Generate(rnd, maxSub, perm); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rule) Apply(s *State) {
	perm := r.getPermutation()
	var srcState State = *s

	for i := range FoldsCount {
		fold := r.getFold(i)

		in := perm.Gather(s, i)

		out := fold.sle.Eval(in)
		out ^= fold.confusion.Eval(srcState)

		perm.Scatter(s, i, out)
	}
}

func (r *Rule) Revert(s *State) {
	perm := r.getPermutation()

	for i := range FoldsCount {
		fold := r.getFold(i)

		in := perm.Gather(s, i)

		noise := fold.confusion.Eval(*s)
		out := fold.sle.Solve(in ^ noise)

		perm.Scatter(s, i, out)
	}
}

func (r *Rule) getFold(idx int) Fold {
	if idx == 0 {
		return Fold{
			sle:       math.NewSLE(r.arena[:LinearFoldBytes]),
			confusion: math.EmptyConfusionMap(),
		}
	}

	offset := LinearFoldBytes + (idx-1)*FoldBytes

	return Fold{
		sle:       math.NewSLE(r.arena[offset : offset+math.SLEBytes]),
		confusion: math.NewConfusionMap(r.arena[offset+math.SLEBytes : offset+FoldBytes]),
	}
}

func (r *Rule) getPermutation() *math.Permutation {
	const offset = RuleFoldsBytes

	view := r.arena[offset : offset+math.PermutationBytes]

	return math.NewPermutation(view)
}

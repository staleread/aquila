package invertible

import "io"

type Rule struct {
	arena []byte
	perm  Permutation
}

type Fold struct {
	sle       *SLE
	confusion *ConfusionMap
}

func (r *Rule) Generate(rnd io.Reader) error {
	entropyBuf := r.arena[:PermutationBytes]

	if err := r.perm.Generate(rnd, entropyBuf); err != nil {
		return err
	}

	for i := range FoldsCount {
		fold := r.getFold(i)

		if err := fold.sle.Generate(rnd); err != nil {
			return err
		}
		if err := fold.confusion.Generate(rnd, i, r.perm); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rule) Apply(state Block) {
	for i := range FoldsCount {
		fold := r.getFold(i)

		in := r.perm.Gather(state, i)

		out := fold.sle.Eval(in)
		out ^= fold.confusion.Eval(state)

		r.perm.Scatter(state, i, out)
	}
}

func (r *Rule) Revert(state Block) {
	for i := range FoldsCount {
		fold := r.getFold(i)

		in := r.perm.Gather(state, i)

		noise := fold.confusion.Eval(state)
		out := fold.sle.Solve(in ^ noise)

		r.perm.Scatter(state, i, out)
	}
}

func (r *Rule) getFold(idx int) Fold {
	offset := idx * VectorSize

	sleArena := r.arena[offset : offset+SLEBytes]
	confusionArena := r.arena[offset+SLEBytes : offset+FoldBytes]

	return Fold{
		sle:       NewSLE(sleArena),
		confusion: NewConfusionMap(confusionArena),
	}
}

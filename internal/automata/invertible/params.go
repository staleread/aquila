package invertible

import "github.com/staleread/aquila/internal/automata"

const (
	VectorSize       = 16
	FoldsCount       = automata.BlockSize / VectorSize
	ConfusionDegree  = 3
	ConfusionMapSize = ConfusionDegree * (ConfusionDegree + 1) / 2 * VectorSize

	SLEBytes          = VectorSize * VectorSize / 8
	ConfusionMapBytes = ConfusionMapSize
	FoldBytes         = SLEBytes + ConfusionMapBytes
	PermutationBytes  = VectorSize
	RuleBytes         = FoldBytes*FoldsCount - ConfusionMapBytes + PermutationBytes
)

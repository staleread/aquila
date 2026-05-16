package invertible

import "github.com/staleread/aquila/internal/automata"

const (
	VectorSize       = 16
	ConfusionDegree  = 3
	ConfusionMapSize = ConfusionDegree * (ConfusionDegree + 1) / 2 * VectorSize
	PermutationSize  = automata.BlockSize
	FoldsCount       = automata.BlockSize / VectorSize
	RulesCount       = 2

	SLEBytes          = VectorSize * VectorSize / 8
	ConfusionMapBytes = ConfusionMapSize
	PermutationBytes  = PermutationSize
	LinearFoldBytes   = SLEBytes
	FoldBytes         = SLEBytes + ConfusionMapBytes
	RuleFoldsBytes    = LinearFoldBytes + FoldBytes*(FoldsCount-1)
	RuleBytes         = RuleFoldsBytes + PermutationBytes
	CABytes           = RuleBytes * RulesCount
)

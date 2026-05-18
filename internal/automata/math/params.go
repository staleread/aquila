package math

import "github.com/staleread/aquila/internal/automata/core"

const (
	VectorSize       = 16
	ConfusionDegree  = 3
	ConfusionMapSize = (ConfusionDegree*(ConfusionDegree+1)/2 - 1) * VectorSize
	PermutationSize  = core.BlockSize

	SLEBytes          = VectorSize * VectorSize / 8
	ConfusionMapBytes = ConfusionMapSize
	PermutationBytes  = PermutationSize

	SymbolicPolynomialSize       = VectorSize + ConfusionDegree - 1
	EstimatedDensePolynomialSize = 4096
)

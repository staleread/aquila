package invertible

import (
	"github.com/staleread/aquila/internal/automata/core"
	"github.com/staleread/aquila/internal/automata/math"
)

const (
	FoldsCount             = core.BlockSize / math.VectorSize
	RulesCount             = 2
	SymbolicPolynomialSize = math.VectorSize + math.ConfusionDegree - 1

	MaxFoldMonomials             = (math.VectorSize + math.ConfusionDegree - 1) * math.VectorSize
	MaxLinearFoldMonomials       = math.VectorSize * math.VectorSize
	MaxRuleMonomials             = MaxLinearFoldMonomials + MaxFoldMonomials*(FoldsCount-1)
	MaxCAMonomials               = MaxRuleMonomials * RulesCount
	EstimatedDensePolynomialSize = 4096

	LinearFoldBytes = math.SLEBytes
	FoldBytes       = math.SLEBytes + math.ConfusionMapBytes
	RuleFoldsBytes  = LinearFoldBytes + FoldBytes*(FoldsCount-1)
	RuleBytes       = RuleFoldsBytes + math.PermutationBytes
	CABytes         = RuleBytes * RulesCount
)

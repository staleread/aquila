package automata

const (
	FoldsCount          = 8
	FoldSize            = 16
	FoldLinearBytes     = FoldSize * FoldSize / 8
	ConfusionMapDegree  = 3
	ConfusionMapBytes   = ConfusionMapDegree * (ConfusionMapDegree + 1) / 2 * FoldSize
	PermutationSize     = FoldSize - 1
	PermutationBytes    = PermutationSize
	InvertibleRuleBytes = FoldLinearBytes*FoldsCount + ConfusionMapBytes*(FoldsCount-1) + PermutationBytes
)

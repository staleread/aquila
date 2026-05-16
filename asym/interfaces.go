package asym

type GeneralCA interface {
	Apply(dst, src []byte)
}

type InvertibleCA interface {
	GeneralCA
	Revert(dst, src []byte)
}

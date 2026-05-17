package asym

type BlockEncrypter interface {
	Apply(dst, src []byte)
}

type BlockDecrypter interface {
	BlockEncrypter
	Revert(dst, src []byte)
}

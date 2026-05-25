package core

type Word = uint32
type Block = Block96

func LoadBlock(src []byte) *Block {
	return (*Block)(LoadBlock96(src))
}

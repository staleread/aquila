package ca

import (
	f "github.com/staleread/aquila/internal/gf2"
	"github.com/staleread/aquila/internal/nla"
)

type GeneralCA struct {
	size int
	rule nla.SNLE
	tmp  []f.Elt
}

func newGeneralCA(size int, rule nla.SNLE) *GeneralCA {
	tmp := make([]f.Elt, size)
	return &GeneralCA{size, rule, tmp}
}

func (ca *GeneralCA) Apply(state []f.Elt) {
	if len(state) != ca.size {
		panic("Invalid CA state size")
	}

	copy(ca.tmp, state)
	ca.rule.Eval(state, ca.tmp)
}

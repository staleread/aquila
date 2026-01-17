package main

type vec []elem

func zeros(n int) vec {
	return make(vec, n)
}

func rands(n int) vec {
	vec := zeros(n)

	for i := range n {
		vec[i] = randElem()
	}
	return vec
}

func readVec(src []byte) vec {
	return vec(readElems(src))
}

func (a vec) add(b vec) {
	for i := range len(a) {
		a[i] = add(a[i], b[i])
	}
}

func (a vec) sub(b vec) {
	for i := range len(a) {
		a[i] = sub(a[i], b[i])
	}
}

func (v vec) write(dst []byte) {
	writeElems(dst, v)
}

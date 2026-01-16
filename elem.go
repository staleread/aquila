package main

type elem uint8

func randElem() elem {
	return elem(randByte() & 1)
}

func randNzElem() elem {
	return elem(1)
}

func add(a, b elem) elem {
	return a ^ b
}

func sub(a, b elem) elem {
	return a ^ b
}

func mul(a, b elem) elem {
	return a & b
}

func div(a, b elem) elem {
	if b == 0 {
		panic("Division by zero")
	}
	return a
}

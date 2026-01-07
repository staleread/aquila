package main

import "fmt"

type BitPack uint32

const BytesPerPack = 32 / 8

func main() {
	lu := RandLuPack(16)
	fmt.Println(lu)

	b := BitVec{0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 0, 1, 1}
	fmt.Println(b)

	SolveLu(lu, b, b)
	fmt.Println(b)
}

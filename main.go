package main

import "fmt"

func main() {
	lu := RandLUPack(3)
	fmt.Println(&lu)

	b := Vec{ Size: 3, Data: []uint8{ 0, 1, 1 } }
	fmt.Println(&b)

	SolveLU(&lu, &b)
	fmt.Println(&b)
}

package main

import "fmt"

func main() {
	lu := RandLUPack(16)
	mat := FromLUPack(lu)

	fmt.Println(&lu)
	fmt.Println(&mat)
}

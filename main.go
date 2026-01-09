package main

import "fmt"

func main() {
	sle := RandSLE(4)
	b := Rands(4)

	fmt.Println(sle.lu)
	fmt.Println(b)

	sle.Solve(b, b)
	fmt.Println(b)
}

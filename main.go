package main

import "fmt"

func main() {
	sle := RandSLE(4)
	fmt.Println(sle.mat)

	b := Vec{123, 17, 56, 255}
	fmt.Println(b)

	sle.Solve(b, b)
	fmt.Println(b)
}

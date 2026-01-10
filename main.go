package main

import "fmt"

func main() {
	const (
		n      = 4
		folds  = 3
		degree = 3
	)

	decRule := RandDecRule(folds, n, degree)
	fmt.Println(decRule)

	c := Rands(n * folds)
	fmt.Println("Before")
	fmt.Println(c)

	decRule.Apply(c)
	fmt.Println("After")
	fmt.Println(c)
}

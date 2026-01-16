package main

import "fmt"

func main() {
	const (
		n      = 4
		folds  = 2
		degree = 3
	)

	r := randRule(folds, n, degree)
	fmt.Println(r)

	pt := rands(n * folds)
	ct := zeros(n * folds)

	fmt.Println("Plain text")
	fmt.Println(pt)

	r.encrypt(pt, ct)
	fmt.Println("Cipher text")
	fmt.Println(ct)
}

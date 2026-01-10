package main

import "fmt"

func main() {
	const (
		n      = 4
		bBytes = 2
		folds  = bBytes * 8
		degree = 3
	)

	b := RandBVec(n, bBytes)
	fmt.Println(b)

	sle := RandSFoldSLE(n, bBytes)
	fmt.Println("Linear Part\n", sle)

	xLin := ZeroBVec(n, bBytes)
	sle.Solve(xLin, b)
	fmt.Println(xLin)

	snle := RandSFoldSNLE(n, folds, degree)
	fmt.Println("Non-Linear Part", snle)

	bLin := ZeroBVec(n, bBytes)
	snle.Eval(xLin, b)

	// sle.Solve(b, b)
	// fmt.Println(b)
}

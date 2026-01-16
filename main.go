package main

import "fmt"

func main() {
	const (
		bSize  = 6
		folds  = 2
		degree = 3
	)

	r := randRule(folds, bSize/folds, degree)
	sn := r.toSnle()
	fmt.Println(sn)

	pt := rands(bSize)
	fmt.Println("Plain text    ", pt)

	exp := zeros(bSize)
	sn.eval(pt, exp)

	fmt.Println("Cipher text (expected)", exp)

	act := zeros(bSize)
	r.encrypt(pt, act)
	fmt.Println("Cipher text   (actual)", act)

	dec := zeros(bSize)
	r.decrypt(dec, act)
	fmt.Println("Decrypted text", dec)
}

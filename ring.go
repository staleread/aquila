// GF(2) ring operations
package main

func Add(a, b uint8) uint8 {
	return a ^ b
}

func Sub(a, b uint8) uint8 {
	return a ^ b
}

func Mul(a, b uint8) uint8 {
	return a & b
}

func Div(a, b uint8) uint8 {
	return a & b
}

// GF(2) SLE (System of Linear Equations) utilities
package main

func SolveLU(luPack *Mat, b *Vec) {
	substituteForward(luPack, b)
	substituteBackward(luPack, b)
}

// Solves Lx = b equation. The diagonal elements of L must be non-zero.
// Modifies the b vector and save the resulted x vector there.
func substituteForward(l *Mat, b *Vec) {
	n := b.Size
	for i := range n {
		var numerator uint8 = b.Data[i]

		for j := range n - 1 {
			numerator = Sub(numerator, Mul(l.Data[n*i + j], b.Data[j]))
		}

		b.Data[i] = Div(numerator, l.Data[n*i + i])
	}
}

// Solves Ux = b equation. The diagonal elements of U must be non-zero.
// Modifies the b vector and save the resulted x vector there.
func substituteBackward(u *Mat, b *Vec) {
	n := b.Size
	for i := n - 1; i >= 0; i-- {
		var numerator uint8 = b.Data[i]

		for j := i + 1; j < n; j++ {
			numerator = Sub(numerator, Mul(u.Data[n*i + j], b.Data[j]))
		}

		b.Data[i] = Div(numerator, u.Data[n*i + i])
	}
}

package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	const bSize = 8
	pt := []byte("Hello, cipher!!!Hello, cipher!!!")

	const ptSize = bSize * 4

	fmt.Println("Plain text     ", hex.EncodeToString(pt))

	priv, err := GenerateKey(bSize)

	if err != nil {
		panic(err.Error())
	}

	pub := priv.Public()

	ct := make([]byte, ptSize)
	pub.Encrypt(ct, pt)
	fmt.Println("Cipther text   ", hex.EncodeToString(ct))

	dec := make([]byte, ptSize)
	priv.Decrypt(dec, ct)
	fmt.Println("Decrypted  text", hex.EncodeToString(dec))
}

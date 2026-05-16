package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	pt := []byte("Hello, world!!!!")

	fmt.Println("Plain text     ", hex.EncodeToString(pt))

	priv, err := GenerateKey()

	if err != nil {
		panic(err.Error())
	}

	ct := make([]byte, len(pt))
	priv.encrypt(ct, pt)
	fmt.Println("Cipher text    ", hex.EncodeToString(ct))

	dec := make([]byte, len(pt))
	priv.Decrypt(dec, ct)
	fmt.Println("Decrypted  text", hex.EncodeToString(dec))
}

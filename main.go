package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	const bSize = 16
	const ptSize = bSize * 2

	pt := []byte("Hello, cipher!!!Hello, cipher!!!")
	fmt.Println("Plain text     ", hex.EncodeToString(pt))

	k, err := GenerateKey(bSize)

	if err != nil {
		panic(err.Error())
	}

	ct := make([]byte, ptSize)
	k.encryptTest(ct, pt)
	fmt.Println("Cipther text   ", hex.EncodeToString(ct))

	dec := make([]byte, ptSize)
	k.Decrypt(dec, ct)
	fmt.Println("Decrypted  text", hex.EncodeToString(dec))
}

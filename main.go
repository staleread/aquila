package main

import (
	"encoding/hex"
	"fmt"

	"github.com/staleread/aquila/asym"
)

func main() {
	pt := []byte("Hello, world!!!!")

	fmt.Println("Plain text     ", hex.EncodeToString(pt))

	priv, err := asym.GenerateKey()

	if err != nil {
		panic(err.Error())
	}

	ct := make([]byte, len(pt))
	priv.Encrypt(ct, pt)
	fmt.Println("Cipher text    ", hex.EncodeToString(ct))

	dec := make([]byte, len(pt))
	priv.Decrypt(dec, ct)
	fmt.Println("Decrypted  text", hex.EncodeToString(dec))
}

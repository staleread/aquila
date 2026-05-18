package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/staleread/aquila/sym"
)

const keyFile = "aquila.key"

func main() {
	pt := []byte("Hello world!")

	fmt.Println("Plain text    ", hex.EncodeToString(pt))

	var block *sym.AquilaBlock
	if f, err := os.Open(keyFile); err == nil {
		block, _ = sym.Decode(f)
		f.Close()
	} else {
		block, _ = sym.New(rand.Reader)
		f, _ := os.Create(keyFile)
		block.Encode(f)
		f.Close()
	}

	ct := make([]byte, len(pt))
	block.Encrypt(ct, pt)
	fmt.Println("Cipher text    ", hex.EncodeToString(ct))

	dec := make([]byte, len(pt))
	block.Decrypt(dec, ct)
	fmt.Println("Decrypted  text", hex.EncodeToString(dec))
}

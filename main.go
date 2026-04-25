package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	// "runtime"
	"runtime/pprof"
)

var memprofile = flag.String("memprofile", "", "write memory profile to file")

// func init() {
// 	runtime.MemProfileRate = 1
// }

func main() {
	flag.Parse()

	const bSize = 2
	pt := []byte("MarsMars")

	const ptSize = bSize * 4

	fmt.Println("Plain text     ", hex.EncodeToString(pt))

	priv, err := GenerateKey(bSize)

	if err != nil {
		panic(err.Error())
	}

	pub := priv.Public()

	ct := make([]byte, ptSize)
	pub.Encrypt(ct, pt)
	fmt.Println("Cipher text    ", hex.EncodeToString(ct))

	dec := make([]byte, ptSize)
	priv.Decrypt(dec, ct)
	fmt.Println("Decrypted  text", hex.EncodeToString(dec))

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.WriteHeapProfile(f)
	}
}

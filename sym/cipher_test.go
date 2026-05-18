package sym_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/staleread/aquila/sym"
)

func BenchmarkRawBlock(b *testing.B) {
	aquilaBlock, err := sym.New(rand.Reader)
	if err != nil {
		b.Fatalf("failed to create Aquila cipher: %v", err)
	}
	aqSize := aquilaBlock.BlockSize()
	aqSrc, aqDst := make([]byte, aqSize), make([]byte, aqSize)

	aesKey := make([]byte, 32)
	_, _ = rand.Read(aesKey)
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		b.Fatalf("failed to create AES cipher: %v", err)
	}
	aesSize := aesBlock.BlockSize()
	aesSrc, aesDst := make([]byte, aesSize), make([]byte, aesSize)

	b.Run("Aquila/EncryptBlock", func(b *testing.B) {
		for b.Loop() {
			aquilaBlock.Encrypt(aqDst, aqSrc)
		}
	})
	b.Run("Aquila/DecryptBlock", func(b *testing.B) {
		for b.Loop() {
			aquilaBlock.Decrypt(aqDst, aqSrc)
		}
	})

	b.Run("AES/EncryptBlock", func(b *testing.B) {
		for b.Loop() {
			aesBlock.Encrypt(aesDst, aesSrc)
		}
	})
	b.Run("AES/DecryptBlock", func(b *testing.B) {
		for b.Loop() {
			aesBlock.Decrypt(aesDst, aesSrc)
		}
	})
}

func runCBCBenchmark(b *testing.B, block cipher.Block, size int, decrypt bool) {
	src := make([]byte, size)
	dst := make([]byte, size)
	_, _ = rand.Read(src)

	iv := make([]byte, block.BlockSize())
	_, _ = rand.Read(iv)

	b.SetBytes(int64(size))
	b.ResetTimer()

	if decrypt {
		for b.Loop() {
			mode := cipher.NewCBCDecrypter(block, iv)
			mode.CryptBlocks(dst, src)
		}
	} else {
		for b.Loop() {
			mode := cipher.NewCBCEncrypter(block, iv)
			mode.CryptBlocks(dst, src)
		}
	}
}

func BenchmarkMode_Payloads(b *testing.B) {
	aquilaBlock, err := sym.New(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}

	aesKey := make([]byte, 16)
	_, _ = rand.Read(aesKey)
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		b.Fatal(err)
	}

	// The mutliples of 12 (Aquila) and 16 (AES) -byte blocks
	sizes := []int{1008, 65520}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Aquila/CBC-Encrypt-%dBytes", size), func(b *testing.B) {
			runCBCBenchmark(b, aquilaBlock, size, false)
		})
		b.Run(fmt.Sprintf("Aquila/CBC-Decrypt-%dBytes", size), func(b *testing.B) {
			runCBCBenchmark(b, aquilaBlock, size, true)
		})

		b.Run(fmt.Sprintf("AES/CBC-Encrypt-%dBytes", size), func(b *testing.B) {
			runCBCBenchmark(b, aesBlock, size, false)
		})
		b.Run(fmt.Sprintf("AES/CBC-Decrypt-%dBytes", size), func(b *testing.B) {
			runCBCBenchmark(b, aesBlock, size, true)
		})
	}
}

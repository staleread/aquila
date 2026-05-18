package sym_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"testing"

	"github.com/staleread/aquila/sym"
)

func runBlockBenchmark(b *testing.B, block cipher.Block) {
	size := block.BlockSize()
	src := make([]byte, size)
	dst := make([]byte, size)

	_, _ = rand.Read(src)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		block.Encrypt(dst, src)
	}
}

func BenchmarkAquila(b *testing.B) {
	aquilaBlock, err := sym.New(rand.Reader)
	if err != nil {
		b.Fatalf("failed to create Aquila cipher: %v", err)
	}

	b.Run("Encrypt", func(b *testing.B) {
		runBlockBenchmark(b, aquilaBlock)
	})
}

func BenchmarkAES(b *testing.B) {
	aesKey := make([]byte, 32)
	_, _ = rand.Read(aesKey)

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		b.Fatalf("failed to create AES cipher: %v", err)
	}

	b.Run("Encrypt", func(b *testing.B) {
		runBlockBenchmark(b, aesBlock)
	})
}

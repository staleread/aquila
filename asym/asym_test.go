package asym

import (
	"crypto/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/core"
)

func TestEncryptDecrypt(t *testing.T) {
	priv, pub, err := GenerateKeyPair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	payload := []byte("hello world!") // One block (12 bytes)
	if len(payload) != core.BlockBytes {
		payload = make([]byte, core.BlockBytes)
		rand.Read(payload)
	}

	ciphertext, err := pub.Encrypt(rand.Reader, payload)
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := priv.Decrypt(rand.Reader, ciphertext, nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(plaintext) != string(payload) {
		t.Errorf("expected %x, got %x", payload, plaintext)
	}
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for b.Loop() {
		_, _, err := GenerateKeyPair(rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncrypt(b *testing.B) {
	_, pub, err := GenerateKeyPair(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}

	payload := make([]byte, 1024*core.BlockBytes)
	if _, err := rand.Read(payload); err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		_, err := pub.Encrypt(rand.Reader, payload)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(payload)))
}

func BenchmarkDecrypt(b *testing.B) {
	priv, pub, err := GenerateKeyPair(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}

	payload := make([]byte, 1024*core.BlockBytes)
	if _, err := rand.Read(payload); err != nil {
		b.Fatal(err)
	}

	ciphertext, err := pub.Encrypt(rand.Reader, payload)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		_, err := priv.Decrypt(rand.Reader, ciphertext, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(payload)))
}

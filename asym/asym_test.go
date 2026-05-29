package asym_test

import (
	"crypto/rand"
	"testing"

	"github.com/staleread/aquila/asym"
)

func TestEncryptDecrypt(t *testing.T) {
	priv, pub, err := asym.GenerateKeyPair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	payload := make([]byte, asym.BlockSize)
	rand.Read(payload)

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

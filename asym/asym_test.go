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

func TestSignVerify(t *testing.T) {
	priv, pub, err := asym.GenerateKeyPair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	digest := make([]byte, asym.BlockSize)
	rand.Read(digest)

	sig, err := priv.Sign(rand.Reader, digest, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !pub.Verify(digest, sig) {
		t.Error("verification failed for valid signature")
	}

	// Corrupt signature
	sig[0] ^= 1
	if pub.Verify(digest, sig) {
		t.Error("verification succeeded for corrupt signature")
	}
}

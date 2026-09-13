package invertible_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/staleread/aquila/internal/automata/general"
	"github.com/staleread/aquila/internal/automata/invertible"
)

func TestPrivatePublicAndSerializedCAConformance(t *testing.T) {
	private := invertible.NewCA()
	if err := private.Generate(rand.New(rand.NewSource(42))); err != nil {
		t.Fatalf("private generation failed: %v", err)
	}

	public, err := private.DeriveGeneralCA()
	if err != nil {
		t.Fatalf("public derivation failed: %v", err)
	}

	var privateEncoding bytes.Buffer
	if err := private.Save(&privateEncoding); err != nil {
		t.Fatalf("private save failed: %v", err)
	}
	loadedPrivate := invertible.NewCA()
	if err := loadedPrivate.Load(bytes.NewReader(privateEncoding.Bytes())); err != nil {
		t.Fatalf("private load failed: %v", err)
	}

	var publicEncoding bytes.Buffer
	if err := public.Save(&publicEncoding); err != nil {
		t.Fatalf("public save failed: %v", err)
	}
	loadedPublic, err := general.LoadCA(bytes.NewReader(publicEncoding.Bytes()))
	if err != nil {
		t.Fatalf("public load failed: %v", err)
	}

	inputs := [][]byte{
		make([]byte, invertible.StateBytes),
		basisInput(0),
		basisInput(invertible.StateSize - 1),
		bytes.Repeat([]byte{0xff}, invertible.StateBytes),
	}
	for _, input := range inputs {
		privateOutput := make([]byte, invertible.StateBytes)
		loadedPrivateOutput := make([]byte, invertible.StateBytes)
		publicOutput := make([]byte, invertible.StateBytes)
		loadedPublicOutput := make([]byte, invertible.StateBytes)

		private.Apply(privateOutput, input)
		loadedPrivate.Apply(loadedPrivateOutput, input)
		public.Apply(publicOutput, input)
		loadedPublic.Apply(loadedPublicOutput, input)

		if !bytes.Equal(privateOutput, publicOutput) {
			t.Errorf("private/public output mismatch for input %x", input)
		}
		if !bytes.Equal(privateOutput, loadedPrivateOutput) {
			t.Errorf("private save/load changed output for input %x", input)
		}
		if !bytes.Equal(publicOutput, loadedPublicOutput) {
			t.Errorf("public save/load changed output for input %x", input)
		}

		reverted := make([]byte, invertible.StateBytes)
		loadedPrivate.Revert(reverted, privateOutput)
		if !bytes.Equal(reverted, input) {
			t.Errorf("loaded private Revert(Apply(%x)) = %x", input, reverted)
		}
	}
}

func basisInput(bit int) []byte {
	input := make([]byte, invertible.StateBytes)
	input[bit/8] = 1 << (bit % 8)
	return input
}

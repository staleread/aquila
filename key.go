package main

import "errors"

type AquilaConfig struct {
	bSize   int
	folds   int
	polyDeg int
	rules   int
}

type PrivateKey struct {
	config *AquilaConfig
	rules  []*rule
}

var configs = map[int]AquilaConfig{
	16: {
		bSize:   16,
		folds:   8,
		polyDeg: 2,
		rules:   16,
	},
	24: {
		bSize:   24,
		folds:   12,
		polyDeg: 3,
		rules:   24,
	},
	32: {
		bSize:   32,
		folds:   16,
		polyDeg: 3,
		rules:   32,
	},
}

func GenerateKey(bSize int) (*PrivateKey, error) {
	config, ok := configs[bSize]

	if !ok {
		return nil, errors.New("Unsupported block size")
	}

	kRules := make([]*rule, config.rules)
	rSize := bytesToElemCnt(config.bSize)

	for i := range config.rules {
		kRules[i] = randRule(rSize, config.folds, config.polyDeg)
	}
	return &PrivateKey{&config, kRules}, nil
}

func (k *PrivateKey) Decrypt(dst, src []byte) {
	bSize := k.config.bSize

	if len(src)%bSize != 0 {
		panic("Size of cipher text must be a multiple of cipher block size")
	}

	for i := range len(src) / bSize {
		from := bSize * i
		to := bSize * (i + 1)

		sv := readVec(src[from:to])
		dv := make(vec, len(sv))

		for j := len(k.rules) - 1; j >= 0; j-- {
			r := k.rules[j]

			if j%2 == 1 {
				r.decrypt(dv, sv)
			} else {
				r.decrypt(sv, dv)
			}
		}
		sv.write(dst[from:to])
	}
}

func (k *PrivateKey) encryptTest(dst, src []byte) {
	bSize := k.config.bSize

	if len(src)%bSize != 0 {
		panic("Size of plain text must be a multiple of cipher block size")
	}

	for i := range len(src) / bSize {
		from := bSize * i
		to := bSize * (i + 1)

		sv := readVec(src[from:to])
		dv := make(vec, len(sv))

		for j, r := range k.rules {
			if j%2 == 0 {
				r.encrypt(dv, sv)
			} else {
				r.encrypt(sv, dv)
			}
		}
		sv.write(dst[from:to])
	}
}

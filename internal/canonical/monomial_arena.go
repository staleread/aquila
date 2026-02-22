package canonical

import (
	"fmt"
	"github.com/staleread/aquila/internal/field"
)

type MonomialArena struct {
	monomCache map[field.MonomialHash]*field.Monomial
}

func newMonomialArena() *MonomialArena {
	cache := make(map[field.MonomialHash]*field.Monomial)
	return &MonomialArena{cache}
}

func (arena *MonomialArena) GetOrInsertMonomial(m field.Monomial) *field.Monomial {
	hash := m.Hash
	cached, cacheHit := arena.monomCache[hash]

	if !cacheHit {
		arena.monomCache[hash] = &m
		return &m
	}

	if cached.Equals(m) {
		return cached
	}

	fmt.Println("[INFO] Hash colision detected", m, cached)

	cacheKey := hash + 1
	for {
		if cacheKey == hash {
			panic("Intern Arena monom cache is out of hash space!")
		}
		if existing, ok := arena.monomCache[cacheKey]; !ok {
			mPtr := &m
			arena.monomCache[cacheKey] = mPtr
			return mPtr
		} else if existing.Equals(m) {
			return existing
		}
		cacheKey++
	}
}

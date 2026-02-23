package sparse

import (
	"fmt"
	"github.com/staleread/aquila/internal/poly"
)

type MonomialInternPool struct {
	monomCache map[MonomialHash]*Monomial
}

func NewMonomialInternPool() *MonomialInternPool {
	cache := make(map[MonomialHash]*Monomial)
	return &MonomialInternPool{cache}
}

func (pool *MonomialInternPool) CreateMonomial(subs ...poly.Subscript) *Monomial {
	m := NewMonomial(subs...)
	return pool.getOrInsertMonomial(m)
}

func (pool *MonomialInternPool) MulPolynomialBy(a, b Polynomial) {
	aMonoms := make([]*Monomial, 0, len(a))
	for ma := range a.Monomials() {
		aMonoms = append(aMonoms, ma)
	}

	for _, ma := range aMonoms {
		a.ToggleMonomial(ma)

		for mb := range b.Monomials() {
			mPtr := pool.getOrInsertMonomial(ma.Mul(*mb))
			a.ToggleMonomial(mPtr)
		}
	}
}

func (pool *MonomialInternPool) getOrInsertMonomial(m Monomial) *Monomial {
	hash := m.Hash
	cached, cacheHit := pool.monomCache[hash]

	if !cacheHit {
		pool.monomCache[hash] = &m
		return &m
	}

	if cached.Equals(m) {
		return cached
	}

	fmt.Println("[INFO] Hash colision detected", m, cached)

	cacheKey := hash + 1
	for {
		if cacheKey == hash {
			panic("Intern monom pool is out of hash space!")
		}
		if existing, ok := pool.monomCache[cacheKey]; !ok {
			mPtr := &m
			pool.monomCache[cacheKey] = mPtr
			return mPtr
		} else if existing.Equals(m) {
			return existing
		}
		cacheKey++
	}
}

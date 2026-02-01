# aquila 🦅

Implementation of the idea proposed by Puhua Guan in "Cellular Automaton
Public-Key Cryptosystem" paper.

## How it works?

*"Let the ground set be a commutative ring. The enciphering key `E` is a
composition of several time-varying inhomogeneous multifold linear invertible
rules, which is made public. The deciphering key `D`, which is kept private by
the designer, is the set of the individual rules in the composite enciphering
function."*

For fast lookup of monomials inside a polinomial, `MSet-XOR-Hash` incremental
hash was used, which is proved to be set-collision resistant.

## References
* [Puhua Guan "Cellular Automaton Public-Key Cryptosystem"](https://www.complex-systems.com/abstracts/v01_i01_a04/)
* [Incremental Multiset Hash Functions and Their Application to Memory Integrity Checking](https://www.researchgate.net/publication/2481704_Incremental_Multiset_Hash_Functions_and_Their_Application_to_Memory_Integrity_Checking)

package math

func (s Bitset) IsZero() bool {
	var zero Bitset
	return s == zero
}

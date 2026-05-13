package gf2

type SLE struct {
	// L and U packed together sharing the diagram of ones
	m *SquareMatrix
}

func NewSLE(n int, data Vector) (*SLE, error) {
	m, err := NewSquareMatrix(n, data)
	if err != nil {
		return nil, err
	}
	return &SLE{m: m}, nil
}

func (s *SLE) FillRand() {
	s.m.data.FillRand()

	for i := range s.m.n {
		s.m.Set(i, i, One)
	}
}

func (s *SLE) Solve(dst, src Vector) {
	s.m.SubstituteForward(dst, src)
	s.m.SubstituteBackward(dst, dst)
}

func (s *SLE) Eval(dst, src, tmp Vector) {
	s.m.MultiplyUpper(tmp, src)
	s.m.MultiplyLower(dst, tmp)
}

// TODO Optimize me
func (s *SLE) Coefs(coefs *SquareMatrix) {
	n := s.m.n

	for i := range n {
		for j := range n {
			sum := Zero

			for k := range min(i, j) + 1 {
				sum ^= s.m.At(i, k) & s.m.At(k, j)
			}

			coefs.Set(i, j, sum)
		}
	}
}

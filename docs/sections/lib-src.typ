#[
  = ДОДАТОК А

  #v(-1.5em)
  #align(center)[*Програмний код криптографічного ядра системи*]
  #v(1em)
  
  #set par(
    first-line-indent: 0em,
    leading: 0.7em,
    justify: false
  )

  #columns(2)[
  Файл `math/adder.go`

  ```
  package math

  // Adds two pre-sorted polinomials using two-pointer merge
  // Appends the resulting monomials to the 'dst' buffer.
  func AddPolynomials(dst, p1, p2 []Monomial) []Monomial {
    var ptr1, ptr2 int

    for ptr1 < len(p1) && ptr2 < len(p2) {
      m1 := p1[ptr1]
      m2 := p2[ptr2]

      cmp := CompareMonomials(m1, m2)

      switch {
      case cmp > 0:
        dst = append(dst, m1)
        ptr1++

      case cmp < 0:
        dst = append(dst, m2)
        ptr2++

      default:
        // Identical terms cancel each other (a ^ a = 0), so skip
        ptr1++
        ptr2++
      }
    }

    if ptr1 < len(p1) {
      dst = append(dst, p1[ptr1:]...)
    }

    if ptr2 < len(p2) {
      dst = append(dst, p2[ptr2:]...)
    }
    return dst
  }
  ```

  Файл `math/bitset96.go`

  ```
  //go:build block96 || (!block8 && !block16 && !block32 && !block64 && !block128)

  package math

  import (
    "encoding/binary"
    "math/bits"
  )

  const (
    BitsetSize  = 96
    BitsetBytes = BitsetSize / 8
  )

  type Subscript = uint8
  type Bitset [3]uint32

  func NewBitset(subs ...Subscript) Bitset {
    var s Bitset
    for _, i := range subs {
      s.SetAt(i, 1)
    }
    return s
  }

  func (s *Bitset) Read(src []byte) {
    _ = src[11] // BCE check

    s[0] = binary.LittleEndian.Uint32(src[0:4])
    s[1] = binary.LittleEndian.Uint32(src[4:8])
    s[2] = binary.LittleEndian.Uint32(src[8:12])
  }

  func (s *Bitset) Write(dst []byte) {
    _ = dst[11] // BCE check

    binary.LittleEndian.PutUint32(dst[0:4], s[0])
    binary.LittleEndian.PutUint32(dst[4:8], s[1])
    binary.LittleEndian.PutUint32(dst[8:12], s[2])
  }

  func (s Bitset) At(idx Subscript) uint8 {
    return uint8((s[idx/32] >> (idx % 32)) & 1)
  }

  func (s *Bitset) SetAt(idx Subscript, bit uint8) {
    wordIdx := idx / 32
    shift := idx % 32

    s[wordIdx] = (s[wordIdx] &^ (uint32(1) << shift)) | uint32(bit&1)<<shift
  }

  func (s Bitset) Compare(other Bitset) int {
    for i := range len(s) {
      if s[i] == other[i] {
        continue
      }
      if s[i] < other[i] {
        return -1
      }
      return 1
    }
    return 0
  }

  func (s Bitset) Or(other Bitset) Bitset {
    return Bitset{
      s[0] | other[0],
      s[1] | other[1],
      s[2] | other[2],
    }
  }

  func (s *Bitset) XorWith(other Bitset) {
    s[0] ^= other[0]
    s[1] ^= other[1]
    s[2] ^= other[2]
  }

  func (s Bitset) Contains(other Bitset) bool {
    return s[0]&other[0] == other[0] && s[1]&other[1] == other[1] && s[2]&other[2] == other[2]
  }

  func (s Bitset) Subscripts(dst []Subscript) []Subscript {
    for i, w := range s {
      wordOffset := i * 32
      for w != 0 {
        tz := bits.TrailingZeros32(w)
        dst = append(dst, Subscript(wordOffset+tz))
        w &= w - 1
      }
    }
    return dst
  }
  ```

  Файл `math/confusion.go`

  ```
  package math

  import (
    "io"
    "slices"
  )

  const ConfusionMapBytes = (ConfusionDegree*(ConfusionDegree+1)/2 - 1) * VectorSize

  type ConfusionMap struct {
    Data []Subscript
  }

  func NewConfusionMap(arena []byte) ConfusionMap {
    return ConfusionMap{Data: []Subscript(arena)}
  }

  func EmptyConfusionMap() ConfusionMap {
    return ConfusionMap{Data: nil}
  }

  func (m *ConfusionMap) Generate(rnd io.Reader, maxSub Subscript, perm *Permutation) error {
    if _, err := io.ReadFull(rnd, m.Data); err != nil {
      return err
    }

    for i := range m.Data {
      m.Data[i] %= maxSub
    }

    subIdx := 0
    for range VectorSize {
      for j := ConfusionDegree; j > 1; j-- {
        upperBound := int(maxSub) - j

        for k := range j {
          candidateIdx := int(m.Data[subIdx]) % (upperBound + k + 1)
          candidate := perm.Data[candidateIdx]

          duplicate := slices.Contains(m.Data[subIdx-k:subIdx], candidate)

          if duplicate {
            m.Data[subIdx] = perm.Data[upperBound+k]
          } else {
            m.Data[subIdx] = candidate
          }
          subIdx++
        }
      }
    }
    return nil
  }

  func (m *ConfusionMap) Eval(s Bitset) Vector {
    if m == nil || len(m.Data) == 0 {
      return Vector(0)
    }

    var res Vector
    subIdx := 0

    for i := range VectorSize {
      var sum Vector

      for j := range ConfusionDegree - 1 {
        prod := Vector(1)
        subCnt := ConfusionDegree - j

        for range subCnt {
          bit := s.At(m.Data[subIdx])
          prod &= Vector(bit)

          subIdx++
        }
        sum ^= prod
      }
      res |= sum << i
    }
    return res
  }
  ```

  Файл `math/monomial.go`

  ```
  package math

  type Monomial Bitset

  var IdentityMonomial Monomial

  func NewMonomial(subs ...Subscript) Monomial {
    var b Bitset
    for _, sub := range subs {
      b.SetAt(sub, 1)
    }
    return Monomial(b)
  }

  func CompareMonomials(a, b Monomial) int {
    return Bitset(a).Compare(Bitset(b))
  }

  func (m Monomial) Mul(other Monomial) Monomial {
    return Monomial(Bitset(m).Or(Bitset(other)))
  }

  func (m Monomial) Subscripts(dst []uint8) []uint8 {
    return Bitset(m).Subscripts(dst)
  }

  func (m *Monomial) SetAt(sub Subscript, bit uint8) {
    (*Bitset)(m).SetAt(sub, bit)
  }

  func (m Monomial) At(sub Subscript) uint8 {
    return Bitset(m).At(sub)
  }

  func (m Monomial) Eval(b Bitset) uint8 {
    if m == IdentityMonomial || b.Contains(Bitset(m)) {
      return 1
    }
    return 0
  }
  ```

  Файл `math/multiplier.go`

  ```
  package math

  import "slices"

  // Multiplies two polynomials and overwrites the 'dst' buffer with the resulting
  // polynomial. The resulting monomials are sorted and deduplicated.
  func MultiplyPolynomials(dst, p1, p2 []Monomial) []Monomial {
    if len(p1) == 0 || len(p2) == 0 {
      return dst[:0]
    }

    total := len(p1) * len(p2)

    if cap(dst) < total {
      dst = make([]Monomial, total)
    } else {
      dst = dst[:total]
    }

    for i, m1 := range p1 {
      for j, m2 := range p2 {
        dst[i*len(p2)+j] = m1.Mul(m2)
      }
    }

    slices.SortFunc(dst, func(a, b Monomial) int {
      return CompareMonomials(b, a)
    })

    writeIdx := 0
    for readIdx := 0; readIdx < len(dst); {
      curr := dst[readIdx]
      count := 1
      readIdx++

      for readIdx < len(dst) && dst[readIdx] == curr {
        count++
        readIdx++
      }

      if count%2 != 0 {
        dst[writeIdx] = curr
        writeIdx++
      }
    }

    return dst[:writeIdx]
  }
  ```

  Файл `math/permutation.go`

  ```
  package math

  import "io"

  const PermutationSize = BitsetSize
  const PermutationBytes = PermutationSize

  type Permutation struct {
    Data []Subscript
  }

  func NewPermutation(arena []byte) *Permutation {
    return &Permutation{
      Data: arena,
    }
  }

  func (p *Permutation) Generate(rnd io.Reader, tmp []byte) error {
    const n Subscript = PermutationSize

    for i := range n {
      p.Data[i] = i
    }

    if _, err := io.ReadFull(rnd, tmp); err != nil {
      return err
    }

    for i := range n - 1 {
      j := tmp[i]%(n-i) + i
      p.Data[i], p.Data[j] = p.Data[j], p.Data[i]
    }
    return nil
  }

  func (p *Permutation) Gather(b *Bitset, foldIdx int) Vector {
    var res Vector
    offset := foldIdx * VectorSize

    for i := range VectorSize {
      idx := p.Data[offset+i]
      res |= Vector(b.At(idx)) << i
    }
    return res
  }

  func (p *Permutation) Scatter(b *Bitset, foldIdx int, v Vector) {
    offset := foldIdx * VectorSize

    for i := range VectorSize {
      idx := p.Data[offset+i]
      bit := uint8((v >> i) & 1)
      b.SetAt(idx, bit)
    }
  }
  ```

  Файл `math/sle.go`

  ```
  package math

  import (
    "io"
    "unsafe"
  )

  const (
    SLESize  = VectorSize * VectorSize
    SLEBytes = VectorBytes * VectorSize
  )

  type SLE struct {
    data []Vector
  }

  func NewSLE(arena []byte) SLE {
    data := unsafe.Slice((*Vector)(unsafe.Pointer(&arena[0])), len(arena)/VectorBytes)

    return SLE{data}
  }

  func (s *SLE) Generate(rnd io.Reader) error {
    byteView := unsafe.Slice((*byte)(unsafe.Pointer(&s.data[0])), len(s.data)*VectorBytes)
    if _, err := io.ReadFull(rnd, byteView); err != nil {
      return err
    }

    for i := range VectorSize {
      s.data[i] |= 1 << i
    }
    return nil
  }

  func (s *SLE) Solve(b Vector) Vector {
    return s.substituteBackward(s.substituteForward(b))
  }

  func (s *SLE) Eval(x Vector) Vector {
    return s.multiplyLower(s.multiplyUpper(x))
  }

  func (s *SLE) Coefs(dst []Vector) {
    for i := range VectorSize {
      iMask := ^Vector((1 << i) - 1)
      dstRow := s.data[i] & iMask

      for k := range i {
        if (s.data[i]>>k)&1 != 1 {
          continue
        }
        kMask := ^Vector((1 << k) - 1)
        dstRow ^= s.data[k] & kMask
      }
      dst[i] = dstRow
    }
  }

  func (s *SLE) substituteForward(v Vector) Vector {
    var res Vector

    for i, row := range s.data {
      res |= ((v >> i) ^ row.Dot(res)) & 1 << i
    }
    return res
  }

  func (s *SLE) substituteBackward(v Vector) Vector {
    var res Vector

    for i := VectorSize - 1; i >= 0; i-- {
      row := s.data[i]
      res |= ((v >> i) ^ row.Dot(res)) & 1 << i
    }
    return res
  }

  func (s *SLE) multiplyLower(v Vector) Vector {
    var res Vector

    for i, row := range s.data {
      mask := Vector((1 << (i + 1)) - 1)
      res |= row.Dot(v&mask) << i
    }
    return res
  }

  func (s *SLE) multiplyUpper(v Vector) Vector {
    var res Vector

    for i, row := range s.data {
      mask := ^Vector((1 << i) - 1)
      res |= row.Dot(v&mask) << i
    }
    return res
  }
  ```

  Файл `invertible/ca.go`

  ```
  package invertible

  import (
    "fmt"
    "io"

    "github.com/staleread/aquilaconfig"
    "github.com/staleread/aquilamath"
  )

  const (
    StateSize  = math.BitsetSize
    StateBytes = math.BitsetBytes
    RuleBytes  = RuleFoldsBytes + math.PermutationBytes
    CABytes    = StateBytes + RuleBytes*RulesCount
  )

  type State = math.Bitset

  type CA struct {
    arena []byte
    shift State
  }

  func NewCA() *CA {
    return &CA{
      arena: make([]byte, CABytes),
    }
  }

  func (ca *CA) Load(src io.Reader) error {
    if err := config.Current.Check(src); err != nil {
      return err
    }
    if _, err := io.ReadFull(src, ca.arena); err != nil {
      return err
    }
    ca.shift.Read(ca.arena[:StateBytes])
    return nil
  }

  func (ca *CA) Save(dst io.Writer) error {
    if err := config.Current.Write(dst); err != nil {
      return err
    }
    _, err := dst.Write(ca.arena)
    return err
  }

  func (ca *CA) Generate(rnd io.Reader) error {
    if _, err := io.ReadFull(rnd, ca.arena[:StateBytes]); err != nil {
      return fmt.Errorf("failed to generate affine shift: %w", err)
    }
    ca.shift.Read(ca.arena[:StateBytes])

    for i := range RulesCount {
      rule := ca.getRule(i)

      if err := rule.Generate(rnd); err != nil {
        return fmt.Errorf("failed to generate rule %d: %w", i, err)
      }
    }
    return nil
  }

  func (ca *CA) Apply(dst, src []byte) {
    var block State
    block.Read(src)

    for i := range RulesCount {
      rule := ca.getRule(i)
      rule.Apply(&block)
    }

    block.XorWith(ca.shift)
    block.Write(dst)
  }

  func (ca *CA) Revert(dst, src []byte) {
    var block State
    block.Read(src)

    block.XorWith(ca.shift)

    for i := RulesCount - 1; i >= 0; i-- {
      rule := ca.getRule(i)
      rule.Revert(&block)
    }

    block.Write(dst)
  }

  func (ca *CA) getRule(idx int) Rule {
    offset := StateBytes + idx*RuleBytes

    return Rule{
      arena: ca.arena[offset : offset+RuleBytes],
    }
  }
  ```

  Файл `invertible/ca_test.go`

  ```
  package invertible_test

  import (
    "bytes"
    "math/rand"
    "testing"

    "github.com/staleread/aquilainvertible"
  )

  func TestCAInvertibility(t *testing.T) {
    ca := invertible.NewCA()
    rng := rand.New(rand.NewSource(42))

    if err := ca.Generate(rng); err != nil {
      t.Fatalf("Failed to generate CA: %v", err)
    }

    src := make([]byte, invertible.StateBytes)
    intermediate := make([]byte, invertible.StateBytes)
    reverted := make([]byte, invertible.StateBytes)

    for range 10 {
      rng.Read(src)

      original := make([]byte, invertible.StateBytes)
      copy(original, src)

      ca.Apply(intermediate, src)
      if bytes.Equal(intermediate, src) {
        t.Errorf("Apply did not change the block (highly unlikely)")
      }

      ca.Revert(reverted, intermediate)
      if !bytes.Equal(reverted, original) {
        t.Errorf("Revert failed: original=%v, reverted=%v", original, reverted)
      }
    }
  }
  ```

  Файл `invertible/derivation.go`

  ```
  package invertible

  import (
    stdmath "math"

    "github.com/staleread/aquilageneral"
    "github.com/staleread/aquilamath"
  )

  var InitialArenaCapacity = calculateMaxMonomiasPerBit() / 2

  func (ca *CA) DeriveGeneralCA() (*general.CA, error) {
    polynomials := CompileRegistry(ca)
    masterArena := make([]math.Monomial, 0, InitialArenaCapacity*StateSize)
    var offsets [StateSize - 1]uint32

    stateArena := make([]math.Monomial, 0, InitialArenaCapacity)
    prodSrcArena := make([]math.Monomial, 0, InitialArenaCapacity)
    prodDstArena := make([]math.Monomial, 0, InitialArenaCapacity)
    sumArena := make([]math.Monomial, 0, InitialArenaCapacity)
    sumScratchArena := make([]math.Monomial, 0, InitialArenaCapacity)

    for i := range StateSize {
      stateArena = stateArena[:0]
      stateArena = append(stateArena, polynomials.GetPolynomial(RulesCount-1, i)...)

      var subs [StateSize]uint8
      for j := RulesCount - 2; j >= 0; j-- {
        for _, monom := range stateArena {
          subsSlice := monom.Subscripts(subs[:0])
          degree := len(subsSlice)

          if degree == 0 {
            continue
          }

          firstSubscript := int(subsSlice[0])
          firstPoly := polynomials.GetPolynomial(j, firstSubscript)

          if degree == 1 {
            sumScratchArena = sumScratchArena[:0]
            sumScratchArena = math.AddPolynomials(sumScratchArena, sumArena, firstPoly)

            sumArena, sumScratchArena = sumScratchArena, sumArena
            continue
          }

          prodSrcArena = prodSrcArena[:0]
          prodSrcArena = append(prodSrcArena, firstPoly...)

          for _, sub := range subsSlice[1:] {
            nextPoly := polynomials.GetPolynomial(j, int(sub))

            prodDstArena = math.MultiplyPolynomials(prodDstArena, prodSrcArena, nextPoly)

            prodSrcArena, prodDstArena = prodDstArena, prodSrcArena
          }

          sumScratchArena = sumScratchArena[:0]
          sumScratchArena = math.AddPolynomials(sumScratchArena, sumArena, prodSrcArena)

          sumArena, sumScratchArena = sumScratchArena, sumArena
        }

        stateArena = stateArena[:0]
        stateArena = append(stateArena, sumArena...)

        sumArena = sumArena[:0]
      }

      masterArena = append(masterArena, stateArena...)

      if ca.shift.At(math.Subscript(i)) == 1 {
        masterArena = append(masterArena, math.IdentityMonomial)
      }

      if i < len(offsets) {
        offsets[i] = uint32(len(masterArena))
      }
    }

    return &general.CA{
      Arena:   masterArena,
      Offsets: offsets,
    }, nil
  }

  func calculateMaxMonomiasPerBit() int {
    L := math.VectorSize / 2
    d0 := math.ConfusionDegree
    r := RulesCount - 1

    M := L + d0 - 1
    for i := 1; i <= r; i++ {
      if M <= 1 {
        M = M*M + L*M
        continue
      }
      mFloat := float64(M)
      p1 := stdmath.Pow(mFloat, float64(d0+1))
      p2 := stdmath.Pow(mFloat, 2.0)
      M = int((p1-p2)/float64(M-1)) + L*M
    }
    return M
  }
  ```

  Файл `invertible/derivation_test.go`

  ```
  package invertible_test

  import (
    "bytes"
    "math/rand"
    "testing"

    "github.com/staleread/aquilainvertible"
  )

  func TestDeriveGeneralCA(t *testing.T) {
    ca := invertible.NewCA()
    rng := rand.New(rand.NewSource(42))

    if err := ca.Generate(rng); err != nil {
      t.Fatalf("Failed to generate CA: %v", err)
    }

    gen, err := ca.DeriveGeneralCA()
    if err != nil {
      t.Fatalf("Failed to derive general CA: %v", err)
    }

    src := make([]byte, invertible.StateBytes)
    dstInvertible := make([]byte, invertible.StateBytes)
    dstGeneral := make([]byte, invertible.StateBytes)

    for i := range 5 {
      rng.Read(src)

      ca.Apply(dstInvertible, src)
      gen.Apply(dstGeneral, src)

      if !bytes.Equal(dstInvertible, dstGeneral) {
        t.Errorf("Inconsistency found at iteration %d:\nInvertible Apply: %x\nGeneral Apply:    %x", i, dstInvertible, dstGeneral)
      }
    }
  }
  ```

  Файл `invertible/registry.go`

  ```
  package invertible

  import (
    "slices"

    "github.com/staleread/aquilamath"
  )

  const (
    FoldsCount             = StateSize / math.VectorSize
    SymbolicPolynomialSize = math.VectorSize + math.ConfusionDegree - 1
    MaxFoldMonomials       = (math.VectorSize + math.ConfusionDegree - 1) * math.VectorSize
    MaxLinearFoldMonomials = math.VectorSize * math.VectorSize
    MaxRuleMonomials       = MaxLinearFoldMonomials + MaxFoldMonomials*(FoldsCount-1)
    MaxCAMonomials         = MaxRuleMonomials * RulesCount
  )

  type SymbolicRegistry struct {
    arena    []math.Monomial
    offsets  [RulesCount][StateSize + 1]uint32
    invPerms [RulesCount][StateSize]uint8
  }

  func (e *SymbolicRegistry) GetPolynomial(ruleIdx, bitIdx int) []math.Monomial {
    k := e.invPerms[ruleIdx][bitIdx]
    start := e.offsets[ruleIdx][k]
    end := e.offsets[ruleIdx][k+1]
    return e.arena[start:end]
  }

  func CompileRegistry(ca *CA) *SymbolicRegistry {
    e := &SymbolicRegistry{}
    e.arena = make([]math.Monomial, 0, MaxCAMonomials)
    monomials := make([]math.Monomial, 0, SymbolicPolynomialSize)

    for r := range RulesCount {
      rule := ca.getRule(r)
      perm := rule.getPermutation()

      var sleCoefs [math.VectorSize]math.Vector

      for k := range StateSize {
        f := k / math.VectorSize
        i := k % math.VectorSize
        b := perm.Data[k]

        e.invPerms[r][b] = uint8(k)
        e.offsets[r][k] = uint32(len(e.arena))

        fold := rule.getFold(f)

        fold.sle.Coefs(sleCoefs[:])
        row := sleCoefs[i]

        monomials = monomials[:0]

        // SLE part
        foldOffset := f * math.VectorSize
        for j := range math.VectorSize {
          if (row>>j)&1 == 1 {
            var m math.Monomial
            m.SetAt(perm.Data[foldOffset+j], 1)
            monomials = append(monomials, m)
          }
        }

        // Confusion part
        if f > 0 && len(fold.confusion.Data) > 0 {
          const SubsPerBit = math.ConfusionMapBytes / math.VectorSize
          cursor := i * SubsPerBit

          for j := range math.ConfusionDegree - 1 {
            degree := math.ConfusionDegree - j
            var m math.Monomial
            for range degree {
              m.SetAt(fold.confusion.Data[cursor], 1)
              cursor++
            }
            monomials = append(monomials, m)
          }
        }

        slices.SortFunc(monomials, func(a, b math.Monomial) int {
          return math.CompareMonomials(b, a)
        })

        e.arena = append(e.arena, monomials...)
      }
      e.offsets[r][StateSize] = uint32(len(e.arena))
    }
    return e
  }
  ```

  Файл `invertible/rule.go`

  ```
  package invertible

  import (
    "io"

    "github.com/staleread/aquilaconfig"
    "github.com/staleread/aquilamath"
  )

  const (
    RulesCount      = config.CompositionCount + 1
    LinearFoldBytes = math.SLEBytes
    FoldBytes       = math.SLEBytes + math.ConfusionMapBytes
    RuleFoldsBytes  = LinearFoldBytes + FoldBytes*(FoldsCount-1)
  )

  type Rule struct {
    arena []byte
  }

  type Fold struct {
    sle       math.SLE
    confusion math.ConfusionMap
  }

  func (r *Rule) Generate(rnd io.Reader) error {
    perm := r.getPermutation()
    entropyBuf := r.arena[:math.PermutationBytes-1]

    if err := perm.Generate(rnd, entropyBuf); err != nil {
      return err
    }

    for i := range FoldsCount {
      fold := r.getFold(i)

      if err := fold.sle.Generate(rnd); err != nil {
        return err
      }

      if i == 0 {
        continue
      }

      maxSub := math.Subscript(i * math.VectorSize)
      if err := fold.confusion.Generate(rnd, maxSub, perm); err != nil {
        return err
      }
    }
    return nil
  }

  func (r *Rule) Apply(s *State) {
    perm := r.getPermutation()
    var srcState State = *s

    for i := range FoldsCount {
      fold := r.getFold(i)

      in := perm.Gather(s, i)

      out := fold.sle.Eval(in)
      out ^= fold.confusion.Eval(srcState)

      perm.Scatter(s, i, out)
    }
  }

  func (r *Rule) Revert(s *State) {
    perm := r.getPermutation()

    for i := range FoldsCount {
      fold := r.getFold(i)

      in := perm.Gather(s, i)

      noise := fold.confusion.Eval(*s)
      out := fold.sle.Solve(in ^ noise)

      perm.Scatter(s, i, out)
    }
  }

  func (r *Rule) getFold(idx int) Fold {
    if idx == 0 {
      return Fold{
        sle:       math.NewSLE(r.arena[:LinearFoldBytes]),
        confusion: math.EmptyConfusionMap(),
      }
    }

    offset := LinearFoldBytes + (idx-1)*FoldBytes

    return Fold{
      sle:       math.NewSLE(r.arena[offset : offset+math.SLEBytes]),
      confusion: math.NewConfusionMap(r.arena[offset+math.SLEBytes : offset+FoldBytes]),
    }
  }

  func (r *Rule) getPermutation() *math.Permutation {
    const offset = RuleFoldsBytes

    view := r.arena[offset : offset+math.PermutationBytes]

    return math.NewPermutation(view)
  }
  ```

  Файл `config/config.go`

  ```
  package config

  import (
    "fmt"
    "io"

    "github.com/staleread/aquilamath"
  )

  const CAConfigBytes = 4

  type CAConfig struct {
    Block  byte
    Comp   byte
    Fold   byte
    Degree byte
  }

  var Current = CAConfig{
    Block:  byte(math.BitsetSize),
    Comp:   byte(CompositionCount),
    Fold:   byte(math.VectorSize),
    Degree: byte(math.ConfusionDegree),
  }

  func (c CAConfig) Write(w io.Writer) error {
    buf := [CAConfigBytes]byte{
      c.Block,
      c.Comp,
      c.Fold,
      c.Degree,
    }
    _, err := w.Write(buf[:])
    return err
  }

  func (c CAConfig) Check(r io.Reader) error {
    var buf [CAConfigBytes]byte
    if _, err := io.ReadFull(r, buf[:]); err != nil {
      return err
    }
    block := int(buf[0])
    comp := int(buf[1])
    fold := int(buf[2])
    deg := int(buf[3])

    if block != int(c.Block) {
      panic(fmt.Sprintf("incompatible block size: got %d, want %d", block, c.Block))
    }
    if comp != int(c.Comp) {
      panic(fmt.Sprintf("incompatible composition count: got %d, want %d", comp, c.Comp))
    }
    if fold != int(c.Fold) {
      panic(fmt.Sprintf("incompatible fold size: got %d, want %d", fold, c.Fold))
    }
    if deg != int(c.Degree) {
      panic(fmt.Sprintf("incompatible confusion degree: got %d, want %d", deg, c.Degree))
    }
    return nil
  }
  ```

  Файл `general/ca.go`

  ```
  package general

  import (
    "encoding/binary"
    "io"

    "github.com/staleread/aquilaconfig"
    "github.com/staleread/aquilamath"
  )

  const (
    StateSize  = math.BitsetSize
    StateBytes = math.BitsetBytes
  )

  type State = math.Bitset

  type CA struct {
    Arena   []math.Monomial
    Offsets [StateSize - 1]uint32
  }

  func (ca *CA) GetPolynomial(idx int) []math.Monomial {
    start := uint32(0)
    if idx > 0 {
      start = ca.Offsets[idx-1]
    }
    end := uint32(len(ca.Arena))
    if idx < StateSize-1 {
      end = ca.Offsets[idx]
    }
    return ca.Arena[start:end]
  }

  func (ca *CA) GetMonomialCounts() []int {
    counts := make([]int, StateSize)
    for i := range StateSize {
      counts[i] = len(ca.GetPolynomial(i))
    }
    return counts
  }

  func (ca *CA) Apply(dst, src []byte) {
    var srcState State
    srcState.Read(src)
    dstState := ca.applyOnState(srcState)
    dstState.Write(dst)
  }

  func (ca *CA) Save(dst io.Writer) error {
    if err := config.Current.Write(dst); err != nil {
      return err
    }
    if err := binary.Write(dst, binary.LittleEndian, ca.Offsets); err != nil {
      return err
    }
    if err := binary.Write(dst, binary.LittleEndian, uint32(len(ca.Arena))); err != nil {
      return err
    }
    return binary.Write(dst, binary.LittleEndian, ca.Arena)
  }

  func LoadCA(src io.Reader) (*CA, error) {
    if err := config.Current.Check(src); err != nil {
      return nil, err
    }

    ca := &CA{}
    if err := binary.Read(src, binary.LittleEndian, &ca.Offsets); err != nil {
      return nil, err
    }
    var length uint32
    if err := binary.Read(src, binary.LittleEndian, &length); err != nil {
      return nil, err
    }
    ca.Arena = make([]math.Monomial, length)
    if err := binary.Read(src, binary.LittleEndian, &ca.Arena); err != nil {
      return nil, err
    }
    return ca, nil
  }

  func (ca *CA) applyOnState(srcState State) State {
    var dstState State
    for i := range StateSize {
      var res uint8
      for _, monom := range ca.GetPolynomial(i) {
        res ^= monom.Eval(srcState)
      }
      dstState.SetAt(math.Subscript(i), res)
    }
    return dstState
  }
  ```

  Файл `general/ca_test.go`

  ```
  package general_test

  import (
    "bytes"
    "math/rand"
    "testing"

    "github.com/staleread/aquilainvertible"
  )

  func TestCACompilationCorrectness(t *testing.T) {
    if testing.Short() {
      t.Skip("skipping CA compilation test in short mode")
    }

    rng := rand.New(rand.NewSource(42))
    invCA := invertible.NewCA()
    if err := invCA.Generate(rng); err != nil {
      t.Fatalf("Failed to generate invertible CA: %v", err)
    }

    genCA, err := invCA.DeriveGeneralCA()
    if err != nil {
      t.Fatalf("Failed to compile CA: %v", err)
    }

    src := make([]byte, invertible.StateBytes)
    got := make([]byte, invertible.StateBytes)
    want := make([]byte, invertible.StateBytes)

    for range 5 {
      rng.Read(src)

      invCA.Apply(want, src)
      genCA.Apply(got, src)

      if !bytes.Equal(got, want) {
        t.Errorf("Compiled CA output mismatch!\ngot:  %x\nwant: %x", got, want)
      }
    }
  }
  ```

  Файл `asym/asym.go`

  ```
  package asym

  import (
    "io"

    "github.com/staleread/aquilaconfig"
    "github.com/staleread/aquilainvertible"
    "github.com/staleread/aquilamath"
  )

  const (
    BlockSize    = invertible.StateBytes
    FoldSize     = math.VectorSize
    Degree       = math.ConfusionDegree
    Compositions = config.CompositionCount
  )

  func GenerateKeyPair(rnd io.Reader) (*PrivateKey, *PublicKey, error) {
    priv, err := GeneratePrivateKey(rnd)
    if err != nil {
      return nil, nil, err
    }

    pub, err := priv.PublicKey()
    if err != nil {
      return nil, nil, err
    }

    return priv, pub, nil
  }
  ```

  Файл `asym/asym_test.go`

  ```
  package asym_test

  import (
    "crypto/rand"
    "testing"

    "github.com/staleread/aquilaasym"
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
  ```

  Файл `asym/private.go`

  ```
  package asym

  import (
    "crypto"
    "errors"
    "io"

    "github.com/staleread/aquilainvertible"
  )

  var _ crypto.Decrypter = (*PrivateKey)(nil)
  var _ crypto.Signer = (*PrivateKey)(nil)

  type PrivateKey struct {
    ca  *invertible.CA
    pub *PublicKey
  }

  func GeneratePrivateKey(rnd io.Reader) (*PrivateKey, error) {
    ca := invertible.NewCA()
    if err := ca.Generate(rnd); err != nil {
      return nil, err
    }
    return &PrivateKey{ca: ca}, nil
  }

  func (k *PrivateKey) PublicKey() (*PublicKey, error) {
    if k.pub == nil {
      gen, err := k.ca.DeriveGeneralCA()
      if err != nil {
        return nil, err
      }
      k.pub = &PublicKey{gen}
    }
    return k.pub, nil
  }

  func (k *PrivateKey) Public() crypto.PublicKey {
    if k.pub == nil {
      gen, err := k.ca.DeriveGeneralCA()
      if err != nil {
        return nil
      }
      k.pub = &PublicKey{gen}
    }
    return k.pub
  }

  func (k *PrivateKey) Decode(src io.Reader) error {
    ca := invertible.NewCA()
    if err := ca.Load(src); err != nil {
      return err
    }
    k.ca = ca
    return nil
  }

  func (k *PrivateKey) Encode(dst io.Writer) error {
    return k.ca.Save(dst)
  }

  func (k *PrivateKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
    if len(msg)%invertible.StateBytes != 0 {
      return nil, errors.New("invalid ciphertext length")
    }

    plaintext = make([]byte, len(msg))
    for i := 0; i < len(msg); i += invertible.StateBytes {
      k.ca.Revert(plaintext[i:i+invertible.StateBytes], msg[i:i+invertible.StateBytes])
    }

    return plaintext, nil
  }

  func (k *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
    if len(digest)%invertible.StateBytes != 0 {
      return nil, errors.New("invalid digest length")
    }

    signature = make([]byte, len(digest))
    for i := 0; i < len(digest); i += invertible.StateBytes {
      k.ca.Revert(signature[i:i+invertible.StateBytes], digest[i:i+invertible.StateBytes])
    }

    return signature, nil
  }
  ```

  #colbreak()
  
  Файл `asym/public.go`

  ```
  package asym

  import (
    "bytes"
    "errors"
    "io"

    "github.com/staleread/aquilageneral"
  )

  type PublicKey struct {
    ca *general.CA
  }

  type PublicKeyInfo struct {
    BlockSizeBytes int
    BlockSizeBits  int
    FoldSize       int
    Degree         int
    Compositions   int
    MonomialCounts []int
  }

  func (k *PublicKey) Describe() PublicKeyInfo {
    return PublicKeyInfo{
      BlockSizeBytes: BlockSize,
      BlockSizeBits:  BlockSize * 8,
      FoldSize:       FoldSize,
      Degree:         Degree,
      Compositions:   Compositions,
      MonomialCounts: k.ca.GetMonomialCounts(),
    }
  }

  func (k *PublicKey) Decode(src io.Reader) error {
    ca, err := general.LoadCA(src)
    if err != nil {
      return err
    }
    k.ca = ca
    return nil
  }

  func (k *PublicKey) Encode(dst io.Writer) error {
    return k.ca.Save(dst)
  }

  func (k *PublicKey) Encrypt(rand io.Reader, msg []byte) (ciphertext []byte, err error) {
    if len(msg)%general.StateBytes != 0 {
      return nil, errors.New("invalid plaintext length")
    }

    ciphertext = make([]byte, len(msg))
    for i := 0; i < len(msg); i += general.StateBytes {
      k.ca.Apply(ciphertext[i:i+general.StateBytes], msg[i:i+general.StateBytes])
    }

    return ciphertext, nil
  }

  func (k *PublicKey) Verify(digest []byte, signature []byte) bool {
    if len(signature) != len(digest) || len(signature)%general.StateBytes != 0 {
      return false
    }

    temp := make([]byte, len(signature))
    for i := 0; i < len(signature); i += general.StateBytes {
      k.ca.Apply(temp[i:i+general.StateBytes], signature[i:i+general.StateBytes])
    }

    return bytes.Equal(temp, digest)
  }

  func (k *PublicKey) ExportToANF(w io.Writer, input []byte) error {
    if len(input) != BlockSize {
      return errors.New("invalid input block length")
    }
    var s general.State
    s.Read(input)
    return k.ca.ExportToANF(w, s)
  }
  ```
  ]
]
# Glossary

Code, article and registry use these terms. A new term goes here first. A term
that differs between code and article is listed with both names until the code
is renamed.

| Term | Symbol | Code | Meaning |
| --- | --- | --- | --- |
| State | $X \in \mathbb{F}_2^m$ | `Bitset`, `StateSize` | The cellular automaton state; $m$ is the block size. |
| Block size | $m = s \cdot f$ | `block` in a config ID | Number of state bits. |
| Fold | $F_i$, width $f$ | `fold` in a config ID | One of $s$ disjoint groups of $f$ local coordinates of a Rule. |
| Rule | $R_t$, $0 \le t < q$ | `Rule` | One invertible round map $G_{\pi_t}^{-1} \circ T_t \circ G_{\pi_t}$. |
| Number of Rules | $q$ | `RulesCount` | Rules executed by the map. |
| Composition count | $c = q - 1$ | `CompositionCount`, `comp` / `c` in a config ID | Legacy parameter; `RulesCount = CompositionCount + 1`. The article uses $q$. |
| Nonlinear degree | $d$, $2 \le d \le f$ | `deg` / `d` in a config ID | Highest monomial degree in a confusion map coordinate. |
| Config ID | `<m>c<c>f<f>d<d>` | `id_aquila16c1f8d2` | Build-tag configuration, e.g. `16c1f8d2`: $m=16$, $c=1$ ($q=2$), $f=8$, $d=2$. |
| Permutation | $\pi_t$ | `Permutation` | Maps local coordinate $k$ to global coordinate $\pi_t(k)$. |
| Gather / Scatter | $G_{\pi_t}$, $G_{\pi_t}^{-1}$ | `Gather`, `Scatter` | Coordinate transport global → local and back. |
| Fold linear map | $A_{t,i} = L_{t,i} U_{t,i}$ | `SLE` | Invertible linear map stored as unit lower/upper triangular factors. |
| Confusion map | $\Phi_{t,i}^{(d)}$ | `ConfusionMap` | Nonlinear term of fold $i>0$, in predecessor-fold coordinates only. |
| Affine shift | $V \ne 0$ | affine shift in `ca.go` | Vector XOR-ed after the last Rule. |
| Unshifted map | $Q = R_{q-1} \circ \dots \circ R_0$ | | The composition before $V$. |
| Public map | $P = V \oplus Q$ | `asym.PublicKey`, `general.CA` | The map whose ANF is the public key. |
| Private representation | | `asym.PrivateKey`, `invertible.CA` | Rules and $V$; evaluates $P$ and $P^{-1}$. |
| ANF | | `general.CA`, `PublicKey.ExportToANF` | Algebraic normal form of $P$, one Boolean polynomial per output bit. |

## Evidence terms

| Term | Meaning |
| --- | --- |
| Hypothesis | A testable statement in `hypotheses.yaml`, id `Hnn`. |
| Research hypothesis | `area: research`: a top-level hypothesis of the research program; finer hypotheses point to it with `parent`. |
| Derivation | Evidence kind: a proof in the article, cited by Typst label. |
| Test | Evidence kind: a Go test whose name carries the id. A finite check, never a general or security claim. |
| Experiment | Evidence kind: a bounded protocol with recorded outcomes; its conclusion holds only under the stated protocol. |
| Literature | Evidence kind: a cited source (bib key); states what the source says, not what we verified. |
| Open question | A question with no testable candidate answer yet (README); it becomes a hypothesis once it has one. |

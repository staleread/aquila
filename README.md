# aquila 🦅

Implementation and article for the formalization and evaluation of Aquila, a
Guan-style multivariate public-key construction over $\mathbb{F}_2$ built from
inhomogeneous, time-varying cellular automata with per-Rule coordinate
permutations and an affine shift.

- Claims: [`hypotheses.yaml`](hypotheses.yaml), the only claim registry. Rules: [AGENTS.md](AGENTS.md).
- Terms: [GLOSSARY.md](GLOSSARY.md).
- Article: [`article/main.typ`](article/main.typ), with the formal model in [`article/sections/model.typ`](article/sections/model.typ).

## Layout

| Path | Contents |
| --- | --- |
| `asym/`, `internal/automata/` | library: invertible CA (private representation), ANF compilation (public map), key handling |
| `cmd/aquila/` | CLI: key generation, encrypt/decrypt, ANF export, analysis |
| `tools/hypotheses/` | registry checker |
| `article/` | Typst article, bibliography, abstracts in English and Ukrainian |
| `hypotheses.yaml` | claims |
| `references/` | Typst conversions of papers; the PDFs are git-ignored |

```sh
go test ./...                                 # library and CLI
go run ./tools/hypotheses check --run         # registry against its evidence
typst compile article/main.typ                # article
```

## Research program

- **Topic:** formalization, modification and experimental analysis of an asymmetric
  cryptographic transformation based on Guan's inhomogeneous cellular automata.
- **Object:** cryptographic transformations based on cellular automata.
- **Subject:** the modified invertible CA transformation of Guan and the influence of its structural
  parameters on the algebraic, diffusion and structural properties of the public map.
- **Aim:** improve and formalize Guan's model by defining construction algorithms for invertible
  Rules, introducing parameterized nonlinear components and per-Rule permutations, and
  experimentally studying the resulting transformation.
- **Research hypotheses:** `H01` (permutations vs. a static partition) and `H02` (effect of $d$).
- **Not a goal:** a proven-secure post-quantum encryption or signature scheme.

## Argument

The article's argument, and what each step currently rests on:

| # | Step | Hypotheses |
| --- | --- | --- |
| 1 | MPKC motivation; Guan's folded construction as the baseline | to re-admit (literature) |
| 2 | Inherited problems: $P(0)=0$ of the homogeneous map, persistent linear structure of static folds, public-key growth | to re-admit |
| 3 | Aquila's changes: affine shift $V$, per-Rule permutations, parameters $(q, f, d)$ | — |
| 4 | Precise definition of Aquila; invertibility | `H03`, `H04` (proven, derivation) |
| 5 | Implementation conformance to the definition | to re-admit (tests exist in `internal/`) |
| 6 | Structural validation of the corrections (static vs. permuted folds) | `H01` |
| 7 | Forward public-map diffusion | to define |
| 8 | Reduced public-map inversion $P(x)=y$ | to define |
| 9 | Public-ANF size and derivation-cost feasibility | `H02` |
| 10 | Novelty boundary relative to Guan and the repeated-triangular/TTM family | to re-admit |
| 11 | Bounded structural-attack screening (MinRank/kernel chain, decomposition, HOLE, cross-term) | to define |
| 12 | Bounded conclusions | follows from the above |

## Status

| Area | Established |
| --- | --- |
| model | $P$ is bijective (`H03`); $P(0)=V$ (`H04`) |
| everything else | nothing yet re-admitted; see the backlog below |

### Re-admission backlog

The former workspace (`master/doc`, commit `b66acdc`, a local repository without a remote) held
curated claims (`master-article/thesis-claims.yaml`, ids `TCnnn`), historical claims
(`course-project-claims.yaml`, `RMnnn`) and tickets. They are provenance only. Each enters
`hypotheses.yaml` as `open` and is proven again here. Remove a row once it is re-admitted or dropped.

| Legacy ticket | Topic |
| --- | --- |
| T002 | implementation conformance (tests already in the library) |
| T003 | affine shift sampled uniformly from nonzero states |
| T004 | unbiased permutation sampling |
| T005 | confusion map sampling |
| T006 | first fold must not be static across Rules |
| T007 | static-fold correction: proof and controlled ablation |
| T008 | forward-map diffusion over independent keys |
| T009 | classified public-preimage benchmark (SAT / timeout / OOM / error) |
| T010–T012 | degree-stratified ANF metrics, growth model, public-map feasibility |
| T013, T016 | novelty boundary, Guan baseline attribution |
| T014 | bounded structural-attack screening |
| T015 | terminology and API (see GLOSSARY: $q$ vs. `CompositionCount`) |

The bachelor thesis and its experiments are at tag `bachelor`. Its avalanche results measured the
inverse map and its SAT runs left failures unclassified, so neither counts as evidence for $P$.

## Open questions

Questions without a testable candidate answer yet. One becomes a hypothesis once a candidate
answer can be stated.

### Reference configurations

Which exact $(m, q, f, d)$ configurations the article reports for each block size, and the resource
limits under which the larger ones remain feasible.

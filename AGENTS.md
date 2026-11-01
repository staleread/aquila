# AGENTS.md

Rules for anyone, human or agent, changing this repository. It holds the
implementation of Aquila and the article that formalizes and evaluates it.

## Layout

- `README.md`: the only prose entry point: research program, argument, status, open questions.
- `GLOSSARY.md`: terms used by code, article and registry. Use them; add new ones there first.
- `hypotheses.yaml`: the ONLY store of claims. Nothing else (README, article, notes, commit
  messages) may state a claim as established unless a `proven` id backs it.
- `asym/`, `internal/`: the library. `cmd/aquila/`: the CLI. `tools/hypotheses/`: the registry checker.
- `article/`: the Typst article (English; the abstract also in Ukrainian, `abstract-uk.typ`).
  `article/references.bib` is the only bibliography.
- `references/`: Typst text conversions of papers (tracked); the PDFs are git-ignored, never commit them.
- `workdir/`: scratch notes and probes for current investigations. Not authoritative.

## Hypotheses

Each entry is keyed by id `Hnn` (ids only grow; new entries go at the end):

```yaml
H12:
  area: model            # research | model | impl | structure | diffusion | inversion | feasibility | novelty
  status: open           # open | proven | rejected
  evidence: derivation   # derivation | test | experiment | literature
  ref: article/sections/model.typ#h-p-zero   # see below; required once proven, except for tests
  parent: H01            # optional: the research hypothesis this one serves
  rejectedBy: [H13]      # optional
  statement: >-
    One falsifiable sentence with its scope (parameters, sample, protocol).
```

What `proven` means depends on the evidence kind, and the checker enforces it:

| Evidence | `ref` | Proven when |
| --- | --- | --- |
| `derivation` | `<file.typ>#<label>` | the label exists in that file and the proof under it is complete |
| `test` | none | a top-level Go test carries the id in its name (`TestH12_ZeroImage`, `TestSLE_H12`) and passes; skips do not count |
| `experiment` | path to the experiment's protocol and results | the recorded outcomes support the statement under the stated protocol |
| `literature` | bib key | the source states it; we did not verify it |

- `rejected` means a hypothesis in `rejectedBy` is proven. Until then the item stays `open` with
  `rejectedBy` listing the candidate rejector. Keep rejected entries: negative results are results.
- Evidence kinds are not interchangeable. A finite test or a bounded experiment never becomes a
  general, asymptotic or security claim; state the scope in the statement.
- Report timeouts, errors and counterexamples; never narrow a sample silently to make a claim pass.
- A question without a testable candidate answer goes to README "Open questions" and becomes a
  hypothesis once it has one.
- Legacy material (the bachelor thesis at tag `bachelor`, the former `master/doc` claims and
  tickets) is provenance, not evidence. A legacy claim enters the registry as `open` and is
  proven only by evidence in this repository.

Commands:

```sh
go run ./tools/hypotheses check          # validate registry, refs and test ids (after every change)
go run ./tools/hypotheses check --run    # also run every tagged test; proven ids must pass
go run ./tools/hypotheses test H12 H13   # run the tests of these ids and of their rejectors
go run ./tools/hypotheses next           # next free id
yq -P '.H12' hypotheses.yaml             # read one entry (mikefarah yq v4)
yq '.[] | key + " " + .status + " " + .area' hypotheses.yaml   # compact list
```

## Article

- The article cites hypotheses through the Typst labels named in `ref` (`<h-…>`). Do not rename
  such a label without updating the registry.
- Text in the article states only proven claims as facts; open ones are worded as hypotheses or
  left out.
- Both abstracts say the same thing; change them together.

## Code

- `go vet ./... && go test ./...` before committing; add a test next to the code it covers.
- The CLI builds per configuration with build tags: from `cmd/aquila/`, `go run . build 16c1f8d2`
  runs `go build -tags "block16 comp1 fold8 deg2"`; a plain build uses the untagged default.
- Smallest change that satisfies the request; say how it was verified.

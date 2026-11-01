// The article: English text, with the annotation in English and Ukrainian.
// Build: typst compile article/main.typ  (from the repository root)

#set document(title: "Aquila: a formalized Guan-style multivariate construction")
#set text(lang: "en")
#set page(paper: "a4", numbering: "1")
#set par(justify: true)
#set heading(numbering: "1.1")
#set math.equation(numbering: "(1)")

#align(center, text(16pt, strong[Aquila: a formalized Guan-style multivariate construction]))

#include "abstract-en.typ"
#include "abstract-uk.typ"

#include "sections/introduction.typ"
#include "sections/model.typ"
#include "sections/implementation.typ"

#bibliography("references.bib", style: "dstu-gost-7-1-2006.csl")

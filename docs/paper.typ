#import "template.typ": paper
#show: paper

//= РЕФЕРАТ
//
//#pagebreak()
//
//= ABSTRACT
//
//The graduation thesis of the bachelor consists of a software complex, which
//contains a cryptographic transformation library, and a thesis. The
//explanatory note to the final thesis contains [X] pages, [Y] illustrations
//and [Z] tables.
//
//The purpose of the work is the modification, software implementation, and
//experimental study of the computational complexity limits of a public-key
//cryptosystem based on inhomogeneous cellular automata, originally proposed by
//#box[P. Guan]. The key goal is to verify the architectural viability and
//correctness of this mathematical model when scaling to larger block sizes, as
//well as to solve the problem of the combinatorial explosion of monomials
//during rule composition. Important aspects of the system are the tracking of
//the avalanche effect, the experimental evaluation of polynomial behavior, and
//the establishing of practical parameters for modern cryptographic standards.
//
//The work defines the main stages of system development using the Go
//programming language and a dedicated configuration CLI/TUI tool, and analyzes
//existing asymmetric analogues such as RSA and ECC. Formed mathematical
//requirements are presented, including a derived predictive complexity formula
//for polynomial growth, a detailed description of the architecture is
//developed using UML 2.0 diagrams, and optimized matrix manipulation
//algorithms over the GF(2) field are givn, reflecting in detail the
//dependencies and computational pathways of the system components.
//
//The result of the work is a functioning cryptographic software complex
//providing production-ready library components that implement standard
//"crypto" interfaces, successfully demonstrating block processing and verified
//key generation stability based on #box[P. Guan's] foundational concepts.
//
//KEY WORDS: PUBLIC-KEY CRYPTOSYSTEM, CELLULAR AUTOMATA, GO LANGUAGE,
//POLYNOMIAL COMPOSITION
//
//#pagebreak()
//
//#outline(title: [ЗМІСТ]) #pagebreak()
//
//= ПЕРЕЛІК УМОВНИХ ПОЗНАЧЕНЬ
//
//#pagebreak()

#include "sections/intro.typ"

#pagebreak()

#include "sections/section1.typ"

#include "sections/section2.typ"

#include "sections/section3.typ"

#pagebreak()

#include "sections/summary.typ"

#pagebreak()

= Список використаних джерел

#bibliography(
  "references.bib",
  title: none,
  style: "dstu-gost-7-1-2006.csl",
)

//#pagebreak()
//#include "sections/code.typ"
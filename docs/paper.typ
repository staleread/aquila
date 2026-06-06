#import "template.typ": paper
#show: paper

#include "sections/intro.typ"

#pagebreak()

#include "sections/section1.typ"

#include "sections/section2.typ"

#include "sections/section3.typ"

#pagebreak()

#include "sections/summary.typ"

#pagebreak()

= СПИСОК ВИКОРИСТАНИХ ДЖЕРЕЛ

#bibliography(
  "references.bib",
  title: none,
  style: "dstu-gost-7-1-2006.csl",
)

//#pagebreak()
//#include "sections/code.typ"
#import "template.typ": template

#show: template

#include "sections/cover.typ"
#include "sections/intro.typ"
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
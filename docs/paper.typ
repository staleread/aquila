#import "template.typ": template

#show: template

#[
  #set page(numbering: none)
  #include "sections/cover.typ"
  #include "sections/task.typ"
  #include "sections/abstract.typ"
]

= ЗМІСТ
#v(-1em)
#outline()
#pagebreak()

#include "sections/terminology.typ"
#pagebreak()
#include "sections/intro.typ"
#pagebreak()
#include "sections/section1.typ"
#pagebreak()
#include "sections/section2.typ"
#pagebreak()
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

#pagebreak()
#include "sections/lib-src.typ"
#pagebreak()
#include "sections/cli-src.typ"
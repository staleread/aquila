#import "template.typ": template

#show: template


#[
  #set page(numbering: none)
  #include "sections/cover.typ"
  #include "sections/task.typ"
  #include "sections/abstract.typ"
]

= ЗМІСТ
#v(-1.5em)
#outline()

#pagebreak()

#include "sections/intro.typ"
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
#import "template.typ": template

#show: template

#[
  #set page(numbering: none)
  #include "sections/cover.typ"
]

#outline(title: "ЗМІСТ", indent: 0em)

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
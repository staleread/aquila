#import "../common.typ": common

#show: common

#set page(
  paper: "a4",
  numbering: "1",
  number-align: top + right,
  margin: (
    top: 20mm,
    bottom: 20mm,
    left: 25mm,
    right: 15mm,
  ),
)

#set par(
  first-line-indent: (
    amount: 1.25cm,
    all: true,
  ),
  justify: true,
  leading: 1em,
  spacing: 1em
)

#set outline(title: none, indent: 0em)

#let fontSize = 14pt;
#let offset = fontSize + 1.5em

#set text(size: fontSize)

#let chapter-counter = counter("chapter")

#set heading(numbering: none)

#show heading: it => {
  set text(size: fontSize)
  block(above: offset, below: offset, width: 100%, par(justify: true, it.body))
}

#show heading.where(level: 1): it => {
  if it.body.has("text") and it.body.text.starts-with("РОЗДІЛ") {
    chapter-counter.step()
    counter(math.equation).update(0)
    counter(figure).update(0)
  }
  align(center)[
    #it
  ]
}

#show raw.where(block: true): it => {
  set par(first-line-indent: 0pt)
  set text(font: "Courier New", size: 8pt)
  it
}

#set math.equation(
  numbering: eq => {
    let chapter = chapter-counter.get().at(0)
    let equation = counter(math.equation).get().at(0)
    
    if chapter == 0 {
      str(numbering("(1)", equation))
    } else {
      str(numbering("(1.1)", chapter, equation))
    }
  }
)

#show math.equation.where(block: true): it => block(above: offset, below: offset, it)

#set figure(
  supplement: [Рис.],
  numbering: _ => {
    let chapter = chapter-counter.get().at(0)
    let figure = counter(figure).get().at(0)
    
    str(numbering("1.1", chapter, figure))
  },
)
#set figure.caption(separator: [ -- ])
#show figure: it => block(above: offset, below: offset, it)

#show table: it => block(below: offset, it)

#[
  #set page(numbering: none)
  #include "cover.typ"
  #include "task.typ"
  #include "abstract.typ"
]

= ЗМІСТ
#v(-1em)
#outline()
#pagebreak()

#include "terminology.typ"
#pagebreak()
#include "intro.typ"
#pagebreak()
#include "section1.typ"
#pagebreak()
#include "section2.typ"
#pagebreak()
#include "section3.typ"
#pagebreak()
#include "summary.typ"
#pagebreak()

= СПИСОК ВИКОРИСТАНИХ ДЖЕРЕЛ

#bibliography(
  "../references.bib",
  title: none,
  style: "../dstu-gost-7-1-2006.csl",
)

#pagebreak()
#include "lib-src.typ"
#pagebreak()
#include "cli-src.typ"
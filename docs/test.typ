#import "common.typ": common

#show: common

#set page(
  paper: "a4",
  margin: (
    top: 10mm,
    bottom: 20mm,
    left: 10mm,
    right: 15mm,
  ),
)

#set par(
  first-line-indent: (amount: 0cm),
  leading: 0em,
  spacing: 0em
)

#let name = "Ратушняк М.А"

#set text(size: 0.7cm)
#box(stroke: 0.5pt, inset: 0.3cm)[#name]

#box(stroke: 0.5pt, inset: (0.3cm))[#name, бакалаврська робота]

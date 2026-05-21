#let paper(body) = {
  set page(
    paper: "a4",
    margin: (
      top: 20mm,
      bottom: 20mm,
      left: 25mm,
      right: 15mm,
    ),
  )
  
  set par(
    first-line-indent: (
      amount: 1.25cm,
      all: true,
    ),
    justify: true,
    leading: 1.5em,
  )
  
  let fontSize = 14pt;
  
  set text(
    font: "Times New Roman",
    size: fontSize,
    lang: "ua"
  )
  
  set figure(
    supplement: [Рисунок],
    numbering: _ => {
      let headingCnt = str(counter(heading).get().at(0));
      let figureCnt = str(counter(figure).get().at(0));
      headingCnt + "." + figureCnt
    }
  )
  
  set figure.caption(separator: [ -- ])
  
  set heading(numbering: (..nums) => {
     let numbers = nums.pos()
     if numbers.len() == 2 {
        numbering("1.1.", ..numbers)
     }
  })
  
  show heading.where(level: 1): it => {
    counter(figure).update(0)
    align(center)[
      #set text(size: fontSize)
      #block(
        above: 3.5em,
        below: 2.5em,
        upper(it),
      )
    ]
  }
  
  show heading.where(level: 2): it => {
    set text(size: fontSize);
    block(above: 3.5em, below: 2.5em, it)
  }
  
  show outline: it => {
    show heading: set align(center)
    it
  }
  
  show raw.where(block: true): it => {
    set par(first-line-indent: 0pt)
    set text(font: "Courier New", size: 10pt)
    it
  }
  
  body
}
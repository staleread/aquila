#let paper(body) = {
  set page(
    paper: "a4",
    numbering: "1",
    number-align: right,
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
    leading: 1em,
    spacing: 1em
  )
  
  let fontSize = 14pt;
  
  set text(
    font: "Times New Roman",
    size: fontSize,
    lang: "ua"
  )
  
  let chapter-counter = counter("chapter")

  set heading(numbering: none)

  show heading.where(level: 1): it => {
    if it.body.has("text") and it.body.text.starts-with("РОЗДІЛ") {
      chapter-counter.step()
      counter(math.equation).update(0)
    }
    align(center)[
      #set text(size: fontSize)
      #block(
        above: 3.5em,
        below: 2.5em,
        upper(it.body)
      )
    ]
  }
  
  show heading.where(level: 2): it => {
    set text(size: fontSize);
    block(above: 3.5em, below: 2.5em, it)
  }
  
  show heading.where(level: 3): it => {
    set text(size: fontSize);
    block(above: 2em, below: 2em, it)
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
  
  set math.equation(
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
    
  set figure(
    supplement: [Рис.],
    numbering: _ => {
      let heading = chapter-counter.get().at(0)
      let figure = counter(figure).get().at(0)
      
      str(numbering("1.1", chapter, figure))
    },
    caption: [ -- ]
  )
  
  body
}
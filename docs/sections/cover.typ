#align(center)[
  #set par(leading: 1em)
  #set text(size: 14pt)
  
  Міністерство освіти і науки України \
  Чернівецький національний університет імені Юрія Федьковича \
  Навчально-науковий інститут фізико-технічних та комп'ютерних наук \
  Кафедра програмного забезпечення комп'ютерних систем
]

#v(6em)

// Заголовок роботи
#align(center)[
  #set par(leading: 1em)
  #set text(size: 20pt)
  
  #strong("ПОЯСНЮВАЛЬНА ЗАПИСКА") \
  
  #set text(size: 18pt)
  до кваліфікаційної роботи бакалавра \
  
  #set par(justify: false)
  #set text(size: 20pt)
  на тему: «Розробка та дослідження асиметричного шифру на основі неоднорідного клітинного автомату»
]
 
// Блок інформації про студента та керівника (праворуч)
#place(top + right, dy: 27em)[
  #block(width: 70%, align(left)[
    #set par(leading: 0.5em, first-line-indent: 0em, justify: true)
    #set text(size: 14pt)

    #grid(
      columns: (1fr),
      row-gutter: 1.2em,
      [
        Виконав студент _\_#underline("4-го")\__ курсу, групи _\_#underline("443б")\__ 
        спеціальності _\_#underline("121" + " " + "Інженерія програмного забезпечення")_ \
        #h(11em) #text(size: 10pt)[(шифр і назва спеціальності)]
      ]
    )
    
    #v(0.5cm)

    #block(width: 90%, align(left)[
      #grid(
        columns: (1.5fr, 2.5fr, 2.5fr),
        row-gutter: 1em,
        column-gutter: 1em,
        align: (left, center, center),
        
        [], [#text("_"*11) \ #text(size: 10pt)[(підпис)]], [\_#underline[М.А. Ратушняк]\_ \ #text(size: 10pt)[(ініціали та прізвище)]],
        [Керівник], [#text("_"*11) \ #text(size: 10pt)[(підпис)]], [\_\_#underline[С.Е. Остапов]\_\_ \ #text(size: 10pt)[(ініціали та прізвище)]]
      )
    ])
  ])
]

// Блок допуску до захисту (ліворуч)
#place(top+left, dy: 42em)[
  #block(width: 50%, align(left)[
    #set text(size: 12pt)
    #set par(leading: 0.5em, first-line-indent: 0em, spacing: 0.5em)
    
    *До захисту допущено:* \
    *Протокол засідання кафедри №* \_\_#underline("14")\_\_ \
    від «02» червня 2025 р.
    
    #grid(
      columns: (70pt, 1fr, 2fr),
      align: (left, center, center),
      row-gutter: 0.5em,
      [*зав.кафедри*], [#text("_"*10)], [\_\_\_#underline("К.П. Газдюк")\_\_\_],
      [], [#text(size: 10pt)[(підпис)]], [#text(size: 10pt)[(ініціали та прізвище)]]
    )
  ])
]

#v(1fr)

#align(center)[
  #set text(size: 14pt)
  Чернівці – 2026
]

#pagebreak()
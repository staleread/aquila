#align(center)[
  #set par(leading: 1em)
  #text(size: 14pt)[
    Міністерство освіти і науки України \
    Чернівецький національний університет імені Юрія Федьковича \
    Навчально-науковий інститут фізико-технічних та комп'ютерних наук \
    Кафедра програмного забезпечення комп'ютерних систем
  ]
]

#v(6em)

#align(center)[
  #set par(leading: 1em)
  
  #text(size: 20pt)[*ПОЯСНЮВАЛЬНА ЗАПИСКА*] \
  #text(size: 18pt)[до кваліфікаційної роботи бакалавра] \
  
  #set par(justify: false)
  #text(size: 20pt)[
    на тему: «Розробка та дослідження асиметричного шифру на основі
    неоднорідного клітинного автомату»
  ]
]
 
#place(top + right, dy: 27em)[
  #block(width: 70%, align(left)[
    #set par(leading: 0.5em, first-line-indent: 0em, justify: true)
    #set text(size: 14pt)
    
    #grid(
      columns: (1fr),
      row-gutter: 0.5em,
      [
        Виконав студент групи _443б_ \
        спеціальності _#underline[121 Інженерія програмного забезпечення]_ \
      ],
      [#align(center)[#text(size: 10pt)[(шифр і назва спеціальності)]]]
    )
    
    #v(2em)

    #block(width: 90%, align(left)[
      #grid(
        columns: (1.5fr, 2.5fr, 2.5fr),
        row-gutter: 1em,
        column-gutter: 1em,
        align: (left, center, center),
        
        [Студент],
        [#underline(" "*22) \ #text(size: 10pt)[(підпис)]],
        [#underline[~М.А. Ратушняк~] \ #text(size: 10pt)[(ініціали та прізвище)]],
        
        [Керівник],
        [#underline(" "*22) \ #text(size: 10pt)[(підпис)]], 
        [#underline[~~~С.Е. Остапов~~~] \ #text(size: 10pt)[(ініціали та прізвище)]]
      )
    ])
  ])
]

#place(top+left, dy: 42em)[
  #block(width: 50%, align(left)[
    #set text(size: 12pt)
    #set par(leading: 0.5em, first-line-indent: 0em, spacing: 0.5em)
    
    *До захисту допущено:* \
    *Протокол засідання кафедри № 16* \
    від «21» травня 2026 р.
    
    #grid(
      columns: (70pt, 1fr, 2fr),
      align: (left, center, center),
      row-gutter: 0.5em,
      
      [*зав.кафедри*], [#underline(" "*18)], [#underline[~~~~~К.П. Газдюк~~~~~]],
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
**Jqmas** (pronunciado "Jcumas") es un lenguaje de programación general en desarrollo.

## Estado

**v0.1.0** — Fundación. Aritmética, strings, y la función `pres` para imprimir.

## Filosofía

- Sintaxis propia, inventada, con identidad visual única
- Basado en bloques `nombre[contenido]`
- Salto de línea controlado con `notN/` (de 1 a 34 saltos)
- Local-first, portable, open source

## Sintaxis básica

```

// Comentario con //

pres['hola']                  // imprime: hola
pres[3 + 3]                   // imprime: 6
pres[10 * (2 + 3)]            // imprime: 50
pres[19, 'hola', not/'hi']    // imprime:
//   19 hola
//   hi
pres[not2/'con 2 saltos']     // imprime:
//   (línea vacía)
//   (línea vacía)
//   con 2 saltos

```

## Instalación (Termux / Linux)

```bash
pip install -e .
```

Uso

```bash
jqmas run archivo.jqm        # ejecutar un archivo
jqmas repl                   # REPL interactivo
jqmas version                # versión
```

Reglas de notN/

· not/ = 1 salto
· not1/ = 1 salto (igual que not/)
· not2/ a not34/ = N saltos
· not0/ y not35+/ = error
· not/ sin argumento después = no hace nada
· not sin / = error

Roadmap

· v0.1 — Aritmética, strings, pres, notN/
· v0.2 — Variables (haz[...])
· v0.3 — Condicionales (si[...])
· v0.4 — Loops (mientras[...])
· v0.5 — Funciones (haz[fn(a,b) : ...]) → Turing-completo
· v0.6+ — Strings, listas, módulos, tipos

Licencia

MIT
EOF

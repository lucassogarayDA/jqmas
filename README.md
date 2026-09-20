# Jqmas

**Jqmas** (pronunciado "Jcumas") es un lenguaje de programación general en desarrollo, con sintaxis propia y filosofía local-first.

## Estado

**v0.3.0** — Core reescrito en Go. Aritmética, strings, variables, `pres` para imprimir, y `notN/` para saltos de línea.

## Filosofía

Jqmas nació de una necesidad simple: **crear algo propio, diferente.**

No pretende reemplazar a Python, JavaScript o Rust. No busca ser 
el lenguaje más rápido, ni el más popular, ni el más completo. 
Busca ser **Jqmas**.

### 1. Aprendizaje sin fricción

La mayoría de los lenguajes tienen una curva de aprendizaje larga 
y empinada. Jqmas apuesta por lo contrario: **aprender haciendo, 
sin trabas.**

- Sintaxis corta y visual
- Bloques `nombre[contenido]` en vez de llaves y paréntesis
- Errores amigables con sugerencias
- Salto de línea controlado con `notN/`

### 2. Diseño con identidad

Jqmas no quiere parecerse a otros lenguajes. Quiere **tener su 
propia cara.** Cada decisión de diseño busca que un programa en 
Jqmas se vea como **Jqmas**, no como "Python con otra sintaxis".

### 3. Local-first, portable, open source

Jqmas corre donde estés. En Termux, en Linux, donde sea. Sin 
dependencias pesadas, sin cuentas, sin nube. **Tuyo, en tu 
dispositivo.**

### 4. Firme, no débil

Jqmas no es un juguete. No es "un lenguaje de prueba". Es un 
proyecto en serio, con core en Go, tests, y un roadmap claro 
hacia ser **Turing-completo**.

### 5. Sin copiar

Jqmas no copia a otros lenguajes. No es "Python pero con bloques". 
No es "Rust pero más fácil". Es **Jqmas**, con sus propias reglas, 
su propia sintaxis, y su propia identidad.

---

**Jqmas no tiene significado como palabra. Lo tiene como proyecto.**

## Sintaxis básica

```jqmas
// Comentario con //

// Aritmética
pres[3 + 3]                   // imprime: 6
pres[10 * (2 + 3)]            // imprime: 50
pres[2.5 + 1.5]               // imprime: 4.0

// Strings
pres['hola']                  // imprime: hola
pres['línea 1\nlínea 2']      // imprime con salto

// Variables
haz[x = 5]                    // guarda 5 en x
pres[x + 3]                   // imprime: 8

// Múltiples valores
pres[19, 'hola', not/'hi']    // imprime:
//   19 hola
//   hi

// Saltos de línea
pres[not2/'con 2 saltos']     // imprime:
//   (línea vacía)
//   (línea vacía)
//   con 2 saltos
```

Reglas de notN/

Sintaxis Significado
not/ 1 salto
not1/ 1 salto (igual que not/)
not2/ a not34/ N saltos
not0/ ❌ Error
not35+/ ❌ Error
not/ sin argumento No hace nada
not sin / ❌ Error

Instalación (Termux / Linux)

Requisitos

· Python 3.10+
· Go 1.21+ (para el core)

Desde el código fuente

```bash
git clone https://github.com/lucassogarayDA/jqmas.git
cd jqmas
./build.sh          # compila el core en Go
pip install -e .    # instala el CLI en Python
```

Uso

```bash
jqmas run archivo.jqm        # ejecutar un archivo
jqmas repl                   # REPL interactivo
jqmas version                # versión
```

Arquitectura

```
┌─────────────────┐
│  CLI/REPL       │
│  (Python)       │
└────────┬────────┘
         │ JSON
┌────────▼────────┐
│  Core           │
│  (Go)           │
│  - Lexer        │
│  - Parser       │
│  - Evaluator    │
└─────────────────┘
```

· Go → lexer, parser, evaluator (rápido, compilado)
· Python → CLI, REPL (cómodo, flexible)
· JSON → comunicación entre ambos

Tests

31 tests en Go:

· 9 tests de lexer
· 9 tests de parser
· 13 tests de evaluator

```bash
cd go
go test ./...
```

Roadmap

· ✅ v0.1 — Aritmética, strings, pres, notN/
· ✅ v0.2 — Variables (haz[...])
· ✅ v0.3 — Core en Go, errores amigables
· ⏳ v0.4 — Condicionales (si[...])
· ⏳ v0.5 — Loops (mientras[...])
· ⏳ v0.6 — Funciones (haz[fn(a,b) : ...]) → Turing-completo
· ⏳ v0.7+ — Strings avanzados, listas, módulos, tipos

Licencia

MIT

Autor

Lucas Sogaray

· GitHub: @lucassogarayDA
· Reddit: u/PapuSOGA
· TikTok: @Lucassogaray1

---

Si te gusta Jqmas, dejá una ⭐ en el repo. ¡Gracias!

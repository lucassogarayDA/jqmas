# Sintaxis de Jqmas v0.1

## Comentarios

```

// esto es un comentario

```

## Imprimir

```

pres['hola']              // imprime: hola
pres[3 + 3]               // imprime: 6
pres[19, 'hola']          // imprime: 19 hola

```

## Strings

Con comillas simples. Soporta escapes `\n`, `\t`, `\\`, `\'`.

```

pres['hola']              // hola
pres['a\nb']              // a (salto) b

```

## Números

Enteros y decimales.

```

pres[42]                  // 42
pres[3.14]                // 3.14

```

## Operadores

`+`, `-`, `*`, `/` con precedencia estándar. Paréntesis para agrupar.

```

pres[2 + 3 * 4]           // 14
pres[(2 + 3) * 4]         // 20

```

Con `+` también se concatenan strings:

```

pres['a' + 'b']           // ab

```

## Saltos de línea: notN/

Inserta N saltos antes del argumento que sigue.

```

pres[not/'x']             // (salto) x
pres[not2/'x']            // (2 saltos) x
pres[not34/'x']           // 34 saltos (máximo)

```

Reglas:
- `not/` = 1 salto
- `not1/` = igual que `not/`
- `not0/` = error
- `not35/` o mayor = error
- `not/` sin argumento después = no hace nada
- `not` sin `/` = error

## pres[] vacío

Imprime solo un salto de línea.

```

pres[]                    // (solo un salto)

```

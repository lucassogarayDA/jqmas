package lexer

import (
"fmt"
"strconv"
"strings"
)

// Error de lexer con toda la info para formatear lindo desde Python.
type LexerError struct {
Mensaje      string
Linea        int
Columna      int
Sugerencias  []string
Tips         []string
}

func (e *LexerError) Error() string {
return e.Mensaje
}

// Lexer convierte código Jqmas en una lista de tokens.
type Lexer struct {
source string
pos    int
line   int
column int
tokens []Token
}

var keywords = map[string]TokenType{
"pres": PRES,
"haz":  HAZ,
"not":  NOT,
}

var singleChar = map[byte]TokenType{
'+': PLUS,
'-': MINUS,
'*': STAR,
'/': SLASH,
'=': EQUALS,
'(': LPAREN,
')': RPAREN,
'[': LBRACKET,
']': RBRACKET,
',': COMMA,
}

const caracteresValidos = "letras, números, espacios, + - * / = ( ) [ ] , ' y //"

// NewLexer crea un lexer para el código fuente dado.
func NewLexer(source string) *Lexer {
return &Lexer{
source: source,
pos:    0,
line:   1,
column: 1,
tokens: []Token{},
}
}

// Tokenize devuelve la lista completa de tokens.
func (l *Lexer) Tokenize() ([]Token, error) {
for !l.atEnd() {
if err := l.scanToken(); err != nil {
return nil, err
}
}
l.tokens = append(l.tokens, Token{Type: EOF, Value: nil, Line: l.line, Col: l.column})
return l.tokens, nil
}

func (l *Lexer) atEnd() bool {
return l.pos >= len(l.source)
}

func (l *Lexer) peek() byte {
if l.atEnd() {
return 0
}
return l.source[l.pos]
}

func (l *Lexer) peekNext() byte {
if l.pos+1 >= len(l.source) {
return 0
}
return l.source[l.pos+1]
}

func (l *Lexer) advance() byte {
ch := l.source[l.pos]
l.pos++
if ch == '\n' {
l.line++
l.column = 1
} else {
l.column++
}
return ch
}

func (l *Lexer) addToken(typ TokenType, value interface{}) {
l.tokens = append(l.tokens, Token{Type: typ, Value: value, Line: l.line, Col: l.column})
}

func (l *Lexer) error(mensaje string, tips []string) error {
return &LexerError{
Mensaje: mensaje,
Linea:   l.line,
Columna: l.column,
Tips:    tips,
}
}

func (l *Lexer) scanToken() error {
ch := l.peek()

// Espacios en blanco
if ch == ' ' || ch == '\t' || ch == '\r' {
l.advance()
return nil
}

// Salto de línea
if ch == '\n' {
l.advance()
return nil
}

// Comentarios //
if ch == '/' && l.peekNext() == '/' {
for !l.atEnd() && l.peek() != '\n' {
l.advance()
}
return nil
}

// Strings
if ch == '\'' {
return l.scanString()
}

// Números
if ch >= '0' && ch <= '9' {
return l.scanNumber()
}

// Identificadores / keywords
if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
return l.scanIdentifier()
}

// Símbolos de un solo carácter
if typ, ok := singleChar[ch]; ok {
l.advance()
l.addToken(typ, nil)
return nil
}

// Carácter inesperado
tips := []string{fmt.Sprintf("los caracteres válidos son: %s", caracteresValidos)}
switch ch {
case '"':
tips = append([]string{"en Jqmas los strings usan comillas simples: 'texto'"}, tips...)
case '.':
tips = append([]string{"los decimales se escriben sin espacio: 3.14"}, tips...)
case ':':
tips = append([]string{"': ' no es un operador válido en Jqmas"}, tips...)
case '{':
tips = append([]string{"Jqmas no usa llaves, los bloques van con [ ... ]"}, tips...)
}

return &LexerError{
Mensaje: fmt.Sprintf("Carácter inesperado: %q", string(ch)),
Linea:   l.line,
Columna: l.column,
Tips:    tips,
}
}

func (l *Lexer) scanString() error {
l.advance() // consume '
var result strings.Builder

for !l.atEnd() && l.peek() != '\'' {
ch := l.advance()
if ch == '\\' {
if l.atEnd() {
return &LexerError{
Mensaje: "String sin cerrar",
Linea:   l.line,
Columna: l.column,
Tips: []string{
"el string no termina con una comilla simple '",
"recordá cerrar con '",
},
}
}
escaped := l.advance()
switch escaped {
case 'n':
result.WriteByte('\n')
case 't':
result.WriteByte('\t')
case '\\':
result.WriteByte('\\')
case '\'':
result.WriteByte('\'')
default:
result.WriteByte(escaped)
}
} else {
result.WriteByte(ch)
}
}

if l.atEnd() {
return &LexerError{
Mensaje: "String sin cerrar",
Linea:   l.line,
Columna: l.column,
Tips: []string{
"agregá una comilla simple ' para cerrar el string",
"si querés un salto de línea usá \\n",
},
}
}

l.advance() // consume '
l.addToken(STRING, result.String())
return nil
}

func (l *Lexer) scanNumber() error {
start := l.pos
for !l.atEnd() && l.peek() >= '0' && l.peek() <= '9' {
l.advance()
}

// Parte decimal
esDecimal := false
if !l.atEnd() && l.peek() == '.' && l.peekNext() >= '0' && l.peekNext() <= '9' {
esDecimal = true
l.advance() // consume .
for !l.atEnd() && l.peek() >= '0' && l.peek() <= '9' {
l.advance()
}
}

text := l.source[start:l.pos]
if esDecimal {
f, _ := strconv.ParseFloat(text, 64)
l.addToken(NUMBER, f)
} else {
n, _ := strconv.Atoi(text)
l.addToken(NUMBER, n)
}
return nil
}

func (l *Lexer) scanIdentifier() error {
start := l.pos
startColumn := l.column

for !l.atEnd() {
ch := l.peek()
if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
(ch >= '0' && ch <= '9') || ch == '_' {
l.advance()
} else {
break
}
}

text := l.source[start:l.pos]

// Caso especial: notN → NOT + NUMBER
if strings.HasPrefix(text, "not") && len(text) > 3 {
numPart := text[3:]
if n, err := strconv.Atoi(numPart); err == nil {
l.tokens = append(l.tokens, Token{Type: NOT, Value: "not", Line: l.line, Col: startColumn})
l.tokens = append(l.tokens, Token{Type: NUMBER, Value: n, Line: l.line, Col: startColumn + 3})
return nil
}
}

typ, ok := keywords[text]
if !ok {
typ = IDENTIFIER
}
l.tokens = append(l.tokens, Token{Type: typ, Value: text, Line: l.line, Col: startColumn})
return nil
}

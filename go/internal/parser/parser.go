package parser

import (
"fmt"

"github.com/lucassogarayDA/jqmas/go/internal/lexer"
)

// ParserError es un error de sintaxis con info para formatear.
type ParserError struct {
Mensaje     string
Linea       int
Columna     int
Sugerencias []string
Tips        []string
}

func (e *ParserError) Error() string {
return e.Mensaje
}

var keywordsValidas = []string{"pres", "haz", "not"}

// Parser convierte tokens en un AST.
type Parser struct {
tokens []lexer.Token
pos    int
}

// NewParser crea un parser con la lista de tokens.
func NewParser(tokens []lexer.Token) *Parser {
return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) peek() lexer.Token {
return p.tokens[p.pos]
}

func (p *Parser) previous() lexer.Token {
return p.tokens[p.pos-1]
}

func (p *Parser) advance() lexer.Token {
if !p.atEnd() {
p.pos++
}
return p.previous()
}

func (p *Parser) atEnd() bool {
return p.peek().Type == lexer.EOF
}

func (p *Parser) check(t lexer.TokenType) bool {
return p.peek().Type == t
}

func (p *Parser) match(types ...lexer.TokenType) bool {
for _, t := range types {
if p.check(t) {
p.advance()
return true
}
}
return false
}

func (p *Parser) consume(t lexer.TokenType, msg string, tips []string) (lexer.Token, error) {
if p.check(t) {
return p.advance(), nil
}
tok := p.peek()
sugerenciasTips := []string{}
sugerenciasTips = append(sugerenciasTips, tips...)

switch t {
case lexer.RBRACKET:
sugerenciasTips = append(sugerenciasTips, "revisá que cada '[' tenga su ']'")
case lexer.RPAREN:
sugerenciasTips = append(sugerenciasTips, "revisá que cada '(' tenga su ')'")
}

return lexer.Token{}, &ParserError{
Mensaje: msg,
Linea:   tok.Line,
Columna: tok.Col,
Tips:    sugerenciasTips,
}
}

// Parse convierte los tokens en un Program.
func (p *Parser) Parse() (*Program, error) {
statements := []Node{}
for !p.atEnd() {
stmt, err := p.statement()
if err != nil {
return nil, err
}
statements = append(statements, stmt)
}
return &Program{Statements: statements}, nil
}

func (p *Parser) statement() (Node, error) {
if p.check(lexer.PRES) {
return p.presCall()
}
if p.check(lexer.HAZ) {
return p.hazAssignment()
}

tok := p.peek()
if tok.Type == lexer.IDENTIFIER {
name := tok.Value.(string)
sugerencias := buscarSimilar(name, keywordsValidas)
if len(sugerencias) > 0 {
return nil, &ParserError{
Mensaje:     fmt.Sprintf("'%s' no es una instrucción válida", name),
Linea:       tok.Line,
Columna:     tok.Col,
Sugerencias: sugerencias,
Tips:        []string{"las instrucciones válidas son: pres, haz"},
}
}
return nil, &ParserError{
Mensaje: fmt.Sprintf("'%s' no es una instrucción válida", name),
Linea:   tok.Line,
Columna: tok.Col,
Tips:    []string{"las instrucciones válidas son: pres, haz"},
}
}

return nil, &ParserError{
Mensaje: "Se esperaba 'pres' o 'haz'",
Linea:   tok.Line,
Columna: tok.Col,
Tips:    []string{"una instrucción empieza con 'pres[' o con 'haz['"},
}
}

func (p *Parser) hazAssignment() (Node, error) {
if _, err := p.consume(lexer.HAZ, "Se esperaba 'haz'", nil); err != nil {
return nil, err
}
if _, err := p.consume(lexer.LBRACKET, "Se esperaba '[' después de haz",
[]string{"la sintaxis es haz[nombre = valor]"}); err != nil {
return nil, err
}

// Nombre de variable
if !p.check(lexer.IDENTIFIER) {
tok := p.peek()
if tok.Type == lexer.NUMBER {
return nil, &ParserError{
Mensaje: "El nombre de variable no puede ser un número",
Linea:   tok.Line,
Columna: tok.Col,
Tips: []string{
"los nombres empiezan con letra o _",
"por ejemplo: haz[x = 5], haz[contador = 0]",
},
}
}
if tok.Type == lexer.STRING {
return nil, &ParserError{
Mensaje: "El nombre de variable no puede ser un string",
Linea:   tok.Line,
Columna: tok.Col,
Tips: []string{
"los nombres van sin comillas",
"por ejemplo: haz[nombre = 'Lucas']",
},
}
}
return nil, &ParserError{
Mensaje: "Se esperaba un nombre de variable",
Linea:   tok.Line,
Columna: tok.Col,
Tips:    []string{"la sintaxis es haz[nombre = valor]"},
}
}

nameTok := p.advance()
name := nameTok.Value.(string)

// Igual
if !p.check(lexer.EQUALS) {
tok := p.peek()
if tok.Type == lexer.RBRACKET {
return nil, &ParserError{
Mensaje: "Falta el '= valor' en la asignación",
Linea:   tok.Line,
Columna: tok.Col,
Tips: []string{
fmt.Sprintf("escribiste: haz[%s]", name),
fmt.Sprintf("la sintaxis correcta es: haz[%s = valor]", name),
},
}
}
return nil, &ParserError{
Mensaje: "Se esperaba '=' después del nombre de variable",
Linea:   tok.Line,
Columna: tok.Col,
Tips:    []string{"la sintaxis es haz[nombre = valor]"},
}
}

p.advance() // consume =

value, err := p.expression()
if err != nil {
return nil, err
}

if _, err := p.consume(lexer.RBRACKET, "Se esperaba ']' para cerrar haz",
[]string{"la sintaxis es haz[nombre = valor]"}); err != nil {
return nil, err
}

return &Assignment{Name: name, Value: value}, nil
}

func (p *Parser) presCall() (Node, error) {
if _, err := p.consume(lexer.PRES, "Se esperaba 'pres'", nil); err != nil {
return nil, err
}
if _, err := p.consume(lexer.LBRACKET, "Se esperaba '[' después de pres",
[]string{"la sintaxis es pres[expresión]"}); err != nil {
return nil, err
}

args := []Argument{}
if !p.check(lexer.RBRACKET) {
arg, err := p.argument()
if err != nil {
return nil, err
}
args = append(args, arg)
for p.match(lexer.COMMA) {
arg, err := p.argument()
if err != nil {
return nil, err
}
args = append(args, arg)
}
}

if _, err := p.consume(lexer.RBRACKET, "Se esperaba ']' para cerrar pres",
[]string{"la sintaxis es pres[expresión]"}); err != nil {
return nil, err
}

return &PresCall{Arguments: args}, nil
}

func (p *Parser) argument() (Argument, error) {
notCount := 0

if p.check(lexer.NOT) {
p.advance()

if p.check(lexer.NUMBER) {
numTok := p.advance()
n := numTok.Value.(int)
if n < 1 {
return Argument{}, &ParserError{
Mensaje: "not0/ no es válido: el mínimo es 1",
Linea:   numTok.Line,
Columna: numTok.Col,
Tips: []string{
"para 1 salto usá not/ o not1/",
"para 2 saltos usá not2/",
},
}
}
if n > 34 {
return Argument{}, &ParserError{
Mensaje: fmt.Sprintf("not%d/ excede el límite de 34 saltos", n),
Linea:   numTok.Line,
Columna: numTok.Col,
Tips:    []string{"el máximo es not34/"},
}
}
notCount = n
} else {
notCount = 1
}

if _, err := p.consume(lexer.SLASH, "Se esperaba '/' después de not",
[]string{"la sintaxis es notN/ antes de un argumento"}); err != nil {
return Argument{}, err
}

if p.check(lexer.RBRACKET) || p.check(lexer.COMMA) {
return Argument{Value: &NumberLiteral{Value: 0}, NotCount: 0}, nil
}
}

value, err := p.expression()
if err != nil {
return Argument{}, err
}
return Argument{Value: value, NotCount: notCount}, nil
}

func (p *Parser) expression() (Node, error) {
return p.addition()
}

func (p *Parser) addition() (Node, error) {
left, err := p.multiplication()
if err != nil {
return nil, err
}
for p.match(lexer.PLUS, lexer.MINUS) {
opTok := p.previous()
op := "+"
if opTok.Type == lexer.MINUS {
op = "-"
}
right, err := p.multiplication()
if err != nil {
return nil, err
}
left = &BinaryOp{Left: left, Operator: op, Right: right}
}
return left, nil
}

func (p *Parser) multiplication() (Node, error) {
left, err := p.unary()
if err != nil {
return nil, err
}
for p.match(lexer.STAR, lexer.SLASH) {
opTok := p.previous()
op := "*"
if opTok.Type == lexer.SLASH {
op = "/"
}
right, err := p.unary()
if err != nil {
return nil, err
}
left = &BinaryOp{Left: left, Operator: op, Right: right}
}
return left, nil
}

func (p *Parser) unary() (Node, error) {
if p.match(lexer.MINUS) {
operand, err := p.unary()
if err != nil {
return nil, err
}
return &BinaryOp{Left: &NumberLiteral{Value: 0}, Operator: "-", Right: operand}, nil
}
return p.primary()
}

func (p *Parser) primary() (Node, error) {
if p.match(lexer.NUMBER) {
return &NumberLiteral{Value: toFloat(p.previous().Value)}, nil
}
if p.match(lexer.STRING) {
return &StringLiteral{Value: p.previous().Value.(string)}, nil
}
if p.match(lexer.IDENTIFIER) {
return &Variable{Name: p.previous().Value.(string)}, nil
}
if p.match(lexer.LPAREN) {
expr, err := p.expression()
if err != nil {
return nil, err
}
if _, err := p.consume(lexer.RPAREN, "Se esperaba ')'", nil); err != nil {
return nil, err
}
return expr, nil
}

tok := p.peek()
tips := []string{"una expresión puede ser: número, string, variable, o a + b"}
if tok.Type == lexer.PRES {
tips = append([]string{"pres[...] no va dentro de una expresión"}, tips...)
} else if tok.Type == lexer.HAZ {
tips = append([]string{"haz[...] es una instrucción, no una expresión"}, tips...)
}

return nil, &ParserError{
Mensaje: "Se esperaba una expresión",
Linea:   tok.Line,
Columna: tok.Col,
Tips:    tips,
}
}

// buscarSimilar busca strings similares a `nombre` en `opciones`.
func buscarSimilar(nombre string, opciones []string) []string {
nombreLower := toLower(nombre)

// Casos específicos
especificas := map[string]string{
"print": "pres", "echo": "pres", "puts": "pres", "printf": "pres",
"pres_": "pres", "pre": "pres", "prs": "pres",
"has": "haz", "haz_": "haz", "hac": "haz",
"nt": "not", "n0t": "not",
}
if sug, ok := especificas[nombreLower]; ok {
for _, op := range opciones {
if op == sug {
return []string{sug}
}
}
}

// Búsqueda genérica simple (Levenshtein simplificado)
resultado := []string{}
for _, op := range opciones {
if distanciaLevenshtein(nombreLower, toLower(op)) <= 2 {
resultado = append(resultado, op)
}
}
return resultado
}

func toLower(s string) string {
result := []byte(s)
for i := range result {
if result[i] >= 'A' && result[i] <= 'Z' {
result[i] += 32
}
}
return string(result)
}

func distanciaLevenshtein(a, b string) int {
la, lb := len(a), len(b)
if la == 0 {
return lb
}
if lb == 0 {
return la
}
prev := make([]int, lb+1)
for j := 0; j <= lb; j++ {
prev[j] = j
}
for i := 1; i <= la; i++ {
curr := make([]int, lb+1)
curr[0] = i
for j := 1; j <= lb; j++ {
costo := 1
if a[i-1] == b[j-1] {
costo = 0
}
curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+costo)
}
prev = curr
}
return prev[lb]
}

func min3(a, b, c int) int {
if a < b {
if a < c {
return a
}
return c
}
if b < c {
return b
}
return c
}
func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case float64:
		return n
	}
	return 0
}

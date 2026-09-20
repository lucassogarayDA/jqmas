package lexer

// TokenType es el tipo de un token de Jqmas.
type TokenType string

const (
// Literales
NUMBER TokenType = "NUMBER"
STRING TokenType = "STRING"

// Identificadores y palabras clave
PRES TokenType = "PRES"
HAZ  TokenType = "HAZ"
NOT  TokenType = "NOT"
IDENTIFIER TokenType = "IDENTIFIER"

// Operadores
PLUS   TokenType = "PLUS"
MINUS  TokenType = "MINUS"
STAR   TokenType = "STAR"
SLASH  TokenType = "SLASH"
EQUALS TokenType = "EQUALS"

// Delimitadores
LPAREN   TokenType = "LPAREN"
RPAREN   TokenType = "RPAREN"
LBRACKET TokenType = "LBRACKET"
RBRACKET TokenType = "RBRACKET"
COMMA    TokenType = "COMMA"

// Especiales
EOF TokenType = "EOF"
)

// Token representa un token del lenguaje Jqmas.
type Token struct {
Type  TokenType   `json:"type"`
Value interface{} `json:"value"`
Line  int         `json:"line"`
Col   int         `json:"col"`
}

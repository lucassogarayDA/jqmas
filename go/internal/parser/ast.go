package parser

// Node es la interfaz base del AST.
type Node interface {
node()
}

// Program es un programa completo.
type Program struct {
Statements []Node
}

func (p *Program) node() {}

// NumberLiteral es un número.
type NumberLiteral struct {
Value float64
}

func (n *NumberLiteral) node() {}

// StringLiteral es un string.
type StringLiteral struct {
Value string
}

func (s *StringLiteral) node() {}

// Variable es una referencia a variable.
type Variable struct {
Name string
}

func (v *Variable) node() {}

// BinaryOp es una operación binaria.
type BinaryOp struct {
Left     Node
Operator string
Right    Node
}

func (b *BinaryOp) node() {}

// Assignment es haz[nombre = expr].
type Assignment struct {
Name  string
Value Node
}

func (a *Assignment) node() {}

// Argument es un argumento de pres, con not_count opcional.
type Argument struct {
Value    Node
NotCount int
}

// PresCall es pres[...].
type PresCall struct {
Arguments []Argument
}

func (p *PresCall) node() {}

package parser

import (
	"testing"

	"github.com/lucassogarayDA/jqmas/go/internal/lexer"
)

func parse(t *testing.T, source string) *Program {
	t.Helper()
	tokens, err := lexer.NewLexer(source).Tokenize()
	if err != nil {
		t.Fatalf("error lexer: %v", err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("error parser: %v", err)
	}
	return prog
}

func TestPresSimple(t *testing.T) {
	prog := parse(t, "pres['hola']")
	if len(prog.Statements) != 1 {
		t.Fatalf("esperaba 1 statement, obtuve %d", len(prog.Statements))
	}
	pres, ok := prog.Statements[0].(*PresCall)
	if !ok {
		t.Fatalf("esperaba PresCall, obtuve %T", prog.Statements[0])
	}
	if len(pres.Arguments) != 1 {
		t.Fatalf("esperaba 1 argumento, obtuve %d", len(pres.Arguments))
	}
	lit, ok := pres.Arguments[0].Value.(*StringLiteral)
	if !ok {
		t.Fatalf("esperaba StringLiteral, obtuve %T", pres.Arguments[0].Value)
	}
	if lit.Value != "hola" {
		t.Errorf("esperaba hola, obtuve %s", lit.Value)
	}
}

func TestPresVacio(t *testing.T) {
	prog := parse(t, "pres[]")
	pres := prog.Statements[0].(*PresCall)
	if len(pres.Arguments) != 0 {
		t.Errorf("esperaba 0 argumentos, obtuve %d", len(pres.Arguments))
	}
}

func TestPresMultiplesArgumentos(t *testing.T) {
	prog := parse(t, "pres[1, 'a', 2]")
	pres := prog.Statements[0].(*PresCall)
	if len(pres.Arguments) != 3 {
		t.Fatalf("esperaba 3 argumentos, obtuve %d", len(pres.Arguments))
	}
}

func TestPrecedenciaMultiplicacion(t *testing.T) {
	prog := parse(t, "pres[2 + 3 * 4]")
	pres := prog.Statements[0].(*PresCall)
	expr, ok := pres.Arguments[0].Value.(*BinaryOp)
	if !ok {
		t.Fatalf("esperaba BinaryOp, obtuve %T", pres.Arguments[0].Value)
	}
	if expr.Operator != "+" {
		t.Errorf("esperaba +, obtuve %s", expr.Operator)
	}
	right, ok := expr.Right.(*BinaryOp)
	if !ok {
		t.Fatalf("esperaba BinaryOp a la derecha, obtuve %T", expr.Right)
	}
	if right.Operator != "*" {
		t.Errorf("esperaba *, obtuve %s", right.Operator)
	}
}

func TestParentesisCambianPrecedencia(t *testing.T) {
	prog := parse(t, "pres[(2 + 3) * 4]")
	pres := prog.Statements[0].(*PresCall)
	expr := pres.Arguments[0].Value.(*BinaryOp)
	if expr.Operator != "*" {
		t.Errorf("esperaba *, obtuve %s", expr.Operator)
	}
	left := expr.Left.(*BinaryOp)
	if left.Operator != "+" {
		t.Errorf("esperaba +, obtuve %s", left.Operator)
	}
}

func TestNotSinNumero(t *testing.T) {
	prog := parse(t, "pres[not/'x']")
	pres := prog.Statements[0].(*PresCall)
	if pres.Arguments[0].NotCount != 1 {
		t.Errorf("esperaba 1, obtuve %d", pres.Arguments[0].NotCount)
	}
}

func TestNotConNumero(t *testing.T) {
	prog := parse(t, "pres[not3/'x']")
	pres := prog.Statements[0].(*PresCall)
	if pres.Arguments[0].NotCount != 3 {
		t.Errorf("esperaba 3, obtuve %d", pres.Arguments[0].NotCount)
	}
}

func TestHazAssignment(t *testing.T) {
	prog := parse(t, "haz[x = 5]")
	asig, ok := prog.Statements[0].(*Assignment)
	if !ok {
		t.Fatalf("esperaba Assignment, obtuve %T", prog.Statements[0])
	}
	if asig.Name != "x" {
		t.Errorf("esperaba x, obtuve %s", asig.Name)
	}
	lit := asig.Value.(*NumberLiteral)
	if lit.Value != 5 {
		t.Errorf("esperaba 5, obtuve %v", lit.Value)
	}
}

func TestVariable(t *testing.T) {
	prog := parse(t, "pres[x]")
	pres := prog.Statements[0].(*PresCall)
	v, ok := pres.Arguments[0].Value.(*Variable)
	if !ok {
		t.Fatalf("esperaba Variable, obtuve %T", pres.Arguments[0].Value)
	}
	if v.Name != "x" {
		t.Errorf("esperaba x, obtuve %s", v.Name)
	}
}

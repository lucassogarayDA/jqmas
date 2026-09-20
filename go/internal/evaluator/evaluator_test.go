package evaluator

import (
	"testing"

	"github.com/lucassogarayDA/jqmas/go/internal/lexer"
	"github.com/lucassogarayDA/jqmas/go/internal/parser"
)

func ejecutar(t *testing.T, source string) string {
	t.Helper()
	tokens, err := lexer.NewLexer(source).Tokenize()
	if err != nil {
		t.Fatalf("error lexer: %v", err)
	}
	prog, err := parser.NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("error parser: %v", err)
	}
	ev := NewEvaluator()
	if err := ev.Evaluate(prog); err != nil {
		t.Fatalf("error evaluator: %v", err)
	}
	return ev.Output.String()
}

func TestAritmeticaBasica(t *testing.T) {
	casos := map[string]string{
		"pres[2 + 3]":  "5\n",
		"pres[10 - 4]": "6\n",
		"pres[3 * 4]":  "12\n",
		"pres[10 / 2]": "5\n",
	}
	for src, esperado := range casos {
		got := ejecutar(t, src)
		if got != esperado {
			t.Errorf("%q: esperaba %q, obtuve %q", src, esperado, got)
		}
	}
}

func TestPrecedencia(t *testing.T) {
	if got := ejecutar(t, "pres[2 + 3 * 4]"); got != "14\n" {
		t.Errorf("esperaba 14, obtuve %q", got)
	}
	if got := ejecutar(t, "pres[(2 + 3) * 4]"); got != "20\n" {
		t.Errorf("esperaba 20, obtuve %q", got)
	}
}

func TestFlotantes(t *testing.T) {
	if got := ejecutar(t, "pres[2.5 + 1.5]"); got != "4\n" {
		t.Errorf("esperaba 4, obtuve %q", got)
	}
}

func TestStrings(t *testing.T) {
	if got := ejecutar(t, "pres['hola']"); got != "hola\n" {
		t.Errorf("esperaba hola, obtuve %q", got)
	}
}

func TestConcatenacionStrings(t *testing.T) {
	if got := ejecutar(t, "pres['a' + 'b']"); got != "ab\n" {
		t.Errorf("esperaba ab, obtuve %q", got)
	}
}

func TestMultiplesArgumentos(t *testing.T) {
	if got := ejecutar(t, "pres[1, 2, 3]"); got != "1 2 3\n" {
		t.Errorf("esperaba 1 2 3, obtuve %q", got)
	}
}

func TestNotSimple(t *testing.T) {
	if got := ejecutar(t, "pres[not/'x']"); got != "\nx\n" {
		t.Errorf("esperaba \\nx\\n, obtuve %q", got)
	}
}

func TestNotDoble(t *testing.T) {
	if got := ejecutar(t, "pres[not2/'x']"); got != "\n\nx\n" {
		t.Errorf("esperaba dos saltos, obtuve %q", got)
	}
}

func TestPresVacio(t *testing.T) {
	if got := ejecutar(t, "pres[]"); got != "\n" {
		t.Errorf("esperaba salto, obtuve %q", got)
	}
}

func TestVariables(t *testing.T) {
	if got := ejecutar(t, "haz[x = 5]\npres[x]"); got != "5\n" {
		t.Errorf("esperaba 5, obtuve %q", got)
	}
}

func TestVariableReasignacion(t *testing.T) {
	if got := ejecutar(t, "haz[x = 5]\nhaz[x = 10]\npres[x]"); got != "10\n" {
		t.Errorf("esperaba 10, obtuve %q", got)
	}
}

func TestVariableEnExpresion(t *testing.T) {
	if got := ejecutar(t, "haz[a = 2]\nhaz[b = 3]\npres[a + b]"); got != "5\n" {
		t.Errorf("esperaba 5, obtuve %q", got)
	}
}

func TestDivisionPorCero(t *testing.T) {
	tokens, _ := lexer.NewLexer("pres[10 / 0]").Tokenize()
	prog, _ := parser.NewParser(tokens).Parse()
	ev := NewEvaluator()
	err := ev.Evaluate(prog)
	if err == nil {
		t.Error("esperaba error de división por cero")
	}
}

func TestVariableNoDefinida(t *testing.T) {
	tokens, _ := lexer.NewLexer("pres[x]").Tokenize()
	prog, _ := parser.NewParser(tokens).Parse()
	ev := NewEvaluator()
	err := ev.Evaluate(prog)
	if err == nil {
		t.Error("esperaba error de variable no definida")
	}
}

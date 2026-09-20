package lexer

import "testing"

func tipos(t *testing.T, source string) []TokenType {
	t.Helper()
	tokens, err := NewLexer(source).Tokenize()
	if err != nil {
		t.Fatalf("error al tokenizar %q: %v", source, err)
	}
	result := []TokenType{}
	for _, tok := range tokens {
		if tok.Type != EOF {
			result = append(result, tok.Type)
		}
	}
	return result
}

func TestNumerosEnteros(t *testing.T) {
	tokens, err := NewLexer("42").Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Type != NUMBER {
		t.Errorf("esperaba NUMBER, obtuve %s", tokens[0].Type)
	}
	if tokens[0].Value.(int) != 42 {
		t.Errorf("esperaba 42, obtuve %v", tokens[0].Value)
	}
}

func TestNumerosDecimales(t *testing.T) {
	tokens, err := NewLexer("3.14").Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Value.(float64) != 3.14 {
		t.Errorf("esperaba 3.14, obtuve %v", tokens[0].Value)
	}
}

func TestStrings(t *testing.T) {
	tokens, err := NewLexer("'hola'").Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Value.(string) != "hola" {
		t.Errorf("esperaba hola, obtuve %v", tokens[0].Value)
	}
}

func TestStringsConEscapes(t *testing.T) {
	tokens, err := NewLexer(`'a\nb\tc\\d\'e'`).Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	esperado := "a\nb\tc\\d'e"
	if tokens[0].Value.(string) != esperado {
		t.Errorf("esperaba %q, obtuve %q", esperado, tokens[0].Value)
	}
}

func TestPresYParentesis(t *testing.T) {
	tipos := tipos(t, "pres['x']")
	esperado := []TokenType{PRES, LBRACKET, STRING, RBRACKET}
	for i, tt := range esperado {
		if tipos[i] != tt {
			t.Errorf("posición %d: esperaba %s, obtuve %s", i, tt, tipos[i])
		}
	}
}

func TestOperadoresAritmeticos(t *testing.T) {
	tipos := tipos(t, "1 + 2 - 3 * 4 / 5")
	if len(tipos) != 9 {
		t.Fatalf("esperaba 9 tokens, obtuve %d", len(tipos))
	}
	esperado := []TokenType{NUMBER, PLUS, NUMBER, MINUS, NUMBER, STAR, NUMBER, SLASH, NUMBER}
	for i, tt := range esperado {
		if tipos[i] != tt {
			t.Errorf("posición %d: esperaba %s, obtuve %s", i, tt, tipos[i])
		}
	}
}

func TestNotSimple(t *testing.T) {
	tipos := tipos(t, "not/")
	if len(tipos) != 2 || tipos[0] != NOT || tipos[1] != SLASH {
		t.Errorf("esperaba [NOT SLASH], obtuve %v", tipos)
	}
}

func TestNotConNumero(t *testing.T) {
	tokens, err := NewLexer("not2/").Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Type != NOT {
		t.Errorf("esperaba NOT, obtuve %s", tokens[0].Type)
	}
	if tokens[1].Type != NUMBER {
		t.Errorf("esperaba NUMBER, obtuve %s", tokens[1].Type)
	}
	if tokens[1].Value.(int) != 2 {
		t.Errorf("esperaba 2, obtuve %v", tokens[1].Value)
	}
	if tokens[2].Type != SLASH {
		t.Errorf("esperaba SLASH, obtuve %s", tokens[2].Type)
	}
}

func TestComentariosIgnorados(t *testing.T) {
	tokens, err := NewLexer("42 // esto es un comentario").Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Type != NUMBER {
		t.Errorf("esperaba NUMBER, obtuve %s", tokens[0].Type)
	}
	if tokens[1].Type != EOF {
		t.Errorf("esperaba EOF, obtuve %s", tokens[1].Type)
	}
}

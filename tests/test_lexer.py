"""Tests del lexer de Jqmas."""

from jqmas.lexer.lexer import Lexer
from jqmas.lexer.token import TokenType


def tipos(source: str) -> list[TokenType]:
    """Helper: devuelve solo los tipos de token (sin EOF)."""
    return [t.type for t in Lexer(source).tokenize() if t.type != TokenType.EOF]


def test_numeros_enteros():
    tokens = Lexer("42").tokenize()
    assert tokens[0].type == TokenType.NUMBER
    assert tokens[0].value == 42


def test_numeros_decimales():
    tokens = Lexer("3.14").tokenize()
    assert tokens[0].type == TokenType.NUMBER
    assert tokens[0].value == 3.14


def test_strings():
    tokens = Lexer("'hola'").tokenize()
    assert tokens[0].type == TokenType.STRING
    assert tokens[0].value == "hola"


def test_strings_con_escapes():
    tokens = Lexer(r"'a\nb\tc\\d\'e'").tokenize()
    assert tokens[0].value == "a\nb\tc\\d'e"


def test_pres_y_parentesis():
    assert tipos("pres['x']") == [
        TokenType.PRES,
        TokenType.LBRACKET,
        TokenType.STRING,
        TokenType.RBRACKET,
    ]


def test_operadores_aritmeticos():
    assert tipos("1 + 2 - 3 * 4 / 5") == [
        TokenType.NUMBER,
        TokenType.PLUS,
        TokenType.NUMBER,
        TokenType.MINUS,
        TokenType.NUMBER,
        TokenType.STAR,
        TokenType.NUMBER,
        TokenType.SLASH,
        TokenType.NUMBER,
    ]


def test_not_simple():
    assert tipos("not/") == [TokenType.NOT, TokenType.SLASH]


def test_not_con_numero():
    tokens = Lexer("not2/").tokenize()
    assert tokens[0].type == TokenType.NOT
    assert tokens[1].type == TokenType.NUMBER
    assert tokens[1].value == 2
    assert tokens[2].type == TokenType.SLASH


def test_comentarios_ignorados():
    tokens = Lexer("42 // esto es un comentario").tokenize()
    assert tokens[0].type == TokenType.NUMBER
    assert tokens[0].value == 42
    assert tokens[1].type == TokenType.EOF

"""Tests del evaluator de Jqmas."""

import io

import pytest

from jqmas.evaluator.evaluator import Evaluator, EvaluatorError
from jqmas.lexer.lexer import Lexer
from jqmas.parser.parser import Parser


def ejecutar(source: str) -> str:
    """Ejecuta código Jqmas y devuelve lo que imprimió."""
    buf = io.StringIO()
    ast = Parser(Lexer(source).tokenize()).parse()
    Evaluator(output=buf).evaluate(ast)
    return buf.getvalue()


def test_aritmetica_basica():
    assert ejecutar("pres[2 + 3]") == "5\n"
    assert ejecutar("pres[10 - 4]") == "6\n"
    assert ejecutar("pres[3 * 4]") == "12\n"
    assert ejecutar("pres[10 / 2]") == "5\n"


def test_precedencia():
    assert ejecutar("pres[2 + 3 * 4]") == "14\n"
    assert ejecutar("pres[(2 + 3) * 4]") == "20\n"


def test_flotantes():
    assert ejecutar("pres[2.5 + 1.5]") == "4\n"
    assert ejecutar("pres[1.5 * 2]") == "3\n"


def test_division_por_cero():
    with pytest.raises(EvaluatorError):
        ejecutar("pres[1 / 0]")


def test_strings():
    assert ejecutar("pres['hola']") == "hola\n"


def test_concatenacion_strings():
    assert ejecutar("pres['a' + 'b']") == "ab\n"


def test_multiples_argumentos():
    assert ejecutar("pres[1, 2, 3]") == "1 2 3\n"
    assert ejecutar("pres['a', 'b']") == "a b\n"


def test_not_simple():
    assert ejecutar("pres[not/'x']") == "\nx\n"


def test_not_doble():
    assert ejecutar("pres[not2/'x']") == "\n\nx\n"


def test_not_entre_argumentos():
    assert ejecutar("pres[19, 'hola', not/'hi']") == "19 hola\nhi\n"


def test_pres_vacio():
    assert ejecutar("pres[]") == "\n"


def test_multiple_pres():
    assert ejecutar("pres['a']\npres['b']") == "a\nb\n"

"""Tests del parser de Jqmas."""

import pytest

from jqmas.lexer.lexer import Lexer
from jqmas.parser.ast_nodes import (
    BinaryOp,
    NumberLiteral,
    PresCall,
    StringLiteral,
)
from jqmas.parser.parser import Parser, ParserError


def parse(source: str):
    return Parser(Lexer(source).tokenize()).parse()


def test_pres_simple():
    ast = parse("pres['hola']")
    assert len(ast.statements) == 1
    pres = ast.statements[0]
    assert isinstance(pres, PresCall)
    assert len(pres.arguments) == 1
    assert isinstance(pres.arguments[0].value, StringLiteral)
    assert pres.arguments[0].value.value == "hola"
    assert pres.arguments[0].not_count == 0


def test_pres_vacio():
    ast = parse("pres[]")
    assert len(ast.statements[0].arguments) == 0


def test_pres_multiples_argumentos():
    ast = parse("pres[1, 'a', 2]")
    args = ast.statements[0].arguments
    assert len(args) == 3
    assert isinstance(args[0].value, NumberLiteral)
    assert isinstance(args[1].value, StringLiteral)
    assert isinstance(args[2].value, NumberLiteral)


def test_precedencia_multiplicacion():
    ast = parse("pres[2 + 3 * 4]")
    expr = ast.statements[0].arguments[0].value
    assert isinstance(expr, BinaryOp)
    assert expr.operator == "+"
    assert isinstance(expr.right, BinaryOp)
    assert expr.right.operator == "*"


def test_parentesis_cambian_precedencia():
    ast = parse("pres[(2 + 3) * 4]")
    expr = ast.statements[0].arguments[0].value
    assert isinstance(expr, BinaryOp)
    assert expr.operator == "*"
    assert isinstance(expr.left, BinaryOp)
    assert expr.left.operator == "+"


def test_not_sin_numero():
    ast = parse("pres[not/'x']")
    assert ast.statements[0].arguments[0].not_count == 1


def test_not_con_numero():
    ast = parse("pres[not3/'x']")
    assert ast.statements[0].arguments[0].not_count == 3


def test_not0_error():
    with pytest.raises(ParserError):
        parse("pres[not0/'x']")


def test_not35_error():
    with pytest.raises(ParserError):
        parse("pres[not35/'x']")


def test_not_sin_slash_error():
    with pytest.raises(Exception):
        parse("pres[not 'x']")

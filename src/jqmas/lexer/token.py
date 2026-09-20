"""Tipos de token del lenguaje Jqmas."""

from dataclasses import dataclass
from enum import Enum, auto


class TokenType(Enum):
    """Tipos de token posibles en Jqmas."""

    # Literales
    NUMBER = auto()
    STRING = auto()

    # Identificadores y palabras clave
    PRES = auto()
    NOT = auto()
    IDENTIFIER = auto()

    # Operadores
    PLUS = auto()
    MINUS = auto()
    STAR = auto()
    SLASH = auto()

    # Delimitadores
    LPAREN = auto()
    RPAREN = auto()
    LBRACKET = auto()
    RBRACKET = auto()
    COMMA = auto()

    # Especiales
    EOF = auto()


@dataclass
class Token:
    """Un token del lenguaje Jqmas."""

    type: TokenType
    value: object
    line: int
    column: int

    def __repr__(self) -> str:
        return f"Token({self.type.name}, {self.value!r}, {self.line}:{self.column})"

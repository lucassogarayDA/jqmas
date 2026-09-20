"""Nodos del AST de Jqmas."""

from dataclasses import dataclass, field
from typing import Any


@dataclass
class Node:
    """Nodo base del AST."""
    pass


@dataclass
class NumberLiteral(Node):
    """Un número literal (entero o decimal)."""
    value: float


@dataclass
class StringLiteral(Node):
    """Un string literal."""
    value: str


@dataclass
class BinaryOp(Node):
    """Operación binaria: left OP right."""
    left: Node
    operator: str
    right: Node


@dataclass
class Argument(Node):
    """Un argumento de pres, opcionalmente con prefijo notN/."""
    value: Node
    not_count: int = 0   # 0 = sin not, >=1 = N saltos


@dataclass
class PresCall(Node):
    """Llamada a pres[...]."""
    arguments: list[Argument] = field(default_factory=list)


@dataclass
class Program(Node):
    """Un programa completo: lista de sentencias."""
    statements: list[Node] = field(default_factory=list)

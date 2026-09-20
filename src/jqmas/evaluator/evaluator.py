"""Evaluator de Jqmas — ejecuta el AST e imprime resultados."""

import sys

from jqmas.parser.ast_nodes import (
    Argument,
    BinaryOp,
    NumberLiteral,
    PresCall,
    Program,
    StringLiteral,
)


class EvaluatorError(Exception):
    """Error durante la evaluación."""


class Evaluator:
    """Recorre el AST y ejecuta cada sentencia."""

    def __init__(self, output=None):
        self.output = output if output is not None else sys.stdout

    def evaluate(self, node: object) -> object:
        """Evalúa un nodo y devuelve su valor."""
        if isinstance(node, Program):
            for stmt in node.statements:
                self.evaluate(stmt)
            return None

        if isinstance(node, PresCall):
            return self._eval_pres(node)

        if isinstance(node, NumberLiteral):
            return node.value

        if isinstance(node, StringLiteral):
            return node.value

        if isinstance(node, BinaryOp):
            return self._eval_binary(node)

        raise EvaluatorError(f"Nodo desconocido: {type(node).__name__}")

    def _eval_binary(self, node: BinaryOp) -> object:
        left = self.evaluate(node.left)
        right = self.evaluate(node.right)

        if node.operator == "+":
            if isinstance(left, str) and isinstance(right, str):
                return left + right
            if isinstance(left, str) or isinstance(right, str):
                return str(left) + str(right)
            return left + right
        if node.operator == "-":
            return left - right
        if node.operator == "*":
            return left * right
        if node.operator == "/":
            if right == 0:
                raise EvaluatorError("División por cero")
            return left / right

        raise EvaluatorError(f"Operador desconocido: {node.operator}")

    def _format_value(self, value: object) -> str:
        """Formatea un valor para imprimir."""
        if isinstance(value, float):
            if value.is_integer():
                return str(int(value))
        return str(value)

    def _eval_pres(self, node: PresCall) -> None:
        """Ejecuta pres[...], imprimiendo los argumentos."""
        # Construimos la línea actual, y volcamos cuando aparece un notN/
        linea: list[str] = []

        def volcar_linea() -> None:
            if linea:
                self.output.write(" ".join(linea))
                linea.clear()

        for arg in node.arguments:
            # Si hay saltos, primero volcamos la línea actual y después tiramos N saltos
            if arg.not_count > 0:
                volcar_linea()
                for _ in range(arg.not_count):
                    self.output.write("\n")

            value = self.evaluate(arg.value)
            linea.append(self._format_value(value))

        # Volcar lo que quede pendiente
        volcar_linea()

        # Salto final automático
        self.output.write("\n")

"""REPL interactivo de Jqmas."""

import sys

from jqmas import __version__
from jqmas.evaluator.evaluator import Evaluator
from jqmas.lexer.lexer import Lexer
from jqmas.parser.parser import Parser


BANNER = f"""Jqmas v{__version__} — REPL interactivo
Escribí 'salir' o 'exit' para terminar.
Probá: pres['hola']
"""

SALIR = {"salir", "exit", "quit", "q"}


def run_repl() -> None:
    """Inicia el REPL interactivo."""
    print(BANNER)

    evaluator = Evaluator()

    while True:
        try:
            linea = input("jqmas> ")
        except (EOFError, KeyboardInterrupt):
            print()
            break

        linea = linea.strip()

        if not linea:
            continue

        if linea.lower() in SALIR:
            break

        try:
            tokens = Lexer(linea).tokenize()
            ast = Parser(tokens).parse()
            evaluator.evaluate(ast)
        except SyntaxError as e:
            print(f"Error de sintaxis: {e}", file=sys.stderr)
        except Exception as e:
            print(f"Error: {e}", file=sys.stderr)

    print("¡Hasta luego!")

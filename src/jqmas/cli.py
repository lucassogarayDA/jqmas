"""Interfaz de línea de comandos de Jqmas."""

import argparse
import sys
from pathlib import Path

from jqmas import __version__
from jqmas.evaluator.evaluator import Evaluator
from jqmas.lexer.lexer import Lexer
from jqmas.parser.parser import Parser
from jqmas.repl.repl import run_repl


def run_file(path: str) -> int:
    """Ejecuta un archivo .jqm."""
    file_path = Path(path)
    if not file_path.exists():
        print(f"Error: no existe el archivo '{path}'", file=sys.stderr)
        return 1
    if file_path.suffix != ".jqm":
        print(f"Advertencia: '{path}' no termina en .jqm", file=sys.stderr)

    source = file_path.read_text(encoding="utf-8")
    return run_source(source)


def run_source(source: str) -> int:
    """Ejecuta código fuente Jqmas."""
    try:
        tokens = Lexer(source).tokenize()
        ast = Parser(tokens).parse()
        evaluator = Evaluator()
        evaluator.evaluate(ast)
        return 0
    except SyntaxError as e:
        print(f"Error de sintaxis: {e}", file=sys.stderr)
        return 1
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


def main() -> None:
    """Entry point del CLI."""
    parser = argparse.ArgumentParser(
        prog="jqmas",
        description="Jqmas — un lenguaje de programación general.",
    )
    parser.add_argument(
        "--version",
        action="version",
        version=f"Jqmas v{__version__}",
    )

    subparsers = parser.add_subparsers(dest="command")

    run_parser = subparsers.add_parser("run", help="Ejecutar un archivo .jqm")
    run_parser.add_argument("file", help="Ruta del archivo .jqm")

    subparsers.add_parser("repl", help="Abrir el REPL interactivo")
    subparsers.add_parser("version", help="Mostrar la versión")

    args = parser.parse_args()

    if args.command == "run":
        sys.exit(run_file(args.file))
    elif args.command == "repl":
        run_repl()
    elif args.command == "version":
        print(f"Jqmas v{__version__}")
    else:
        parser.print_help()


if __name__ == "__main__":
    main()

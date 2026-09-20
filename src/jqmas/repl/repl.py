"""REPL interactivo de Jqmas."""

import sys

from jqmas import __version__
from jqmas.engine import EngineError, run as engine_run
from jqmas.errors import formatear_error_detail


BANNER = f"""Jqmas v{__version__} — REPL interactivo
Escribí 'salir' o 'exit' para terminar.
Probá: pres['hola'] y haz[x = 5]
"""

SALIR = {"salir", "exit", "quit", "q"}


def run_repl() -> None:
    """Inicia el REPL interactivo."""
    print(BANNER)

    # Acumulamos el estado de las variables entre líneas
    variables_acumuladas: dict = {}

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

        # Componer el código con las variables acumuladas al principio
        # Así cada línea "ve" las variables declaradas antes
        codigo = _componer_codigo(variables_acumuladas, linea)

        try:
            result = engine_run(codigo)
        except EngineError as e:
            print(f"Error: {e}", file=sys.stderr)
            continue

        if result.errors:
            for err in result.errors:
                print(formatear_error_detail(err), file=sys.stderr)
            continue

        # Actualizar variables acumuladas
        variables_acumuladas.update(result.variables)

        # Imprimir solo el output de la nueva línea
        nuevo_output = _extraer_output_nuevo(result.output, variables_acumuladas)
        sys.stdout.write(nuevo_output)

        # Si la línea fue una asignación, mostrar "x = valor"
        asignacion = _detectar_asignacion(linea)
        if asignacion:
            nombre = asignacion
            if nombre in result.variables:
                valor = result.variables[nombre]
                valor_str = _format_valor(valor)
                print(f"{nombre} = {valor_str}")


def _componer_codigo(variables: dict, nueva_linea: str) -> str:
    """Compone el código fuente con las variables acumuladas + la nueva línea."""
    partes = []
    for nombre, valor in variables.items():
        if isinstance(valor, str):
            partes.append(f"haz[{nombre} = '{valor}']")
        else:
            partes.append(f"haz[{nombre} = {valor}]")
    partes.append(nueva_linea)
    return "\n".join(partes)


def _extraer_output_nuevo(output_completo: str, variables: dict) -> str:
    """Como no sabemos qué parte del output es nueva, devolvemos todo.

    Mejora futura: que el core devuelva solo el output de la nueva línea."""
    return output_completo


def _detectar_asignacion(linea: str) -> str | None:
    """Detecta si la línea es haz[nombre = ...] y devuelve el nombre."""
    linea = linea.strip()
    if not linea.startswith("haz["):
        return None
    resto = linea[4:]
    if "=" not in resto:
        return None
    nombre = resto.split("=")[0].strip()
    if nombre.isidentifier():
        return nombre
    return None


def _format_valor(v) -> str:
    """Formatea un valor para mostrar."""
    if isinstance(v, float):
        if v.is_integer():
            return str(int(v))
    return str(v)

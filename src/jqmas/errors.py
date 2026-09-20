"""Sistema de errores amigables de Jqmas."""

from jqmas.engine import ErrorDetail


def formatear_error_detail(err: ErrorDetail) -> str:
    """Formatea un ErrorDetail del engine para mostrarlo lindo."""
    partes = [f"Error: {err.mensaje}"]

    if err.linea is not None and err.columna is not None:
        partes.append(f"       (línea {err.linea}, columna {err.columna})")

    if err.sugerencias:
        if len(err.sugerencias) == 1:
            partes.append(f"       ¿Quizás quisiste decir '{err.sugerencias[0]}'?")
        elif len(err.sugerencias) == 2:
            partes.append(
                f"       ¿Quizás quisiste decir '{err.sugerencias[0]}' "
                f"o '{err.sugerencias[1]}'?"
            )
        else:
            lista = ", ".join(f"'{s}'" for s in err.sugerencias[:-1])
            partes.append(
                f"       ¿Quizás quisiste decir {lista} "
                f"o '{err.sugerencias[-1]}'?"
            )

    for tip in err.tips:
        partes.append(f"       Tip: {tip}")

    return "\n".join(partes)

"""Engine de Jqmas — habla con el binario Go jqmas-core."""

import json
import shutil
import subprocess
from dataclasses import dataclass, field
from pathlib import Path


class EngineError(Exception):
    """Error al comunicarse con jqmas-core."""
    pass


@dataclass
class ErrorDetail:
    """Un error devuelto por el core."""
    mensaje: str
    linea: int | None = None
    columna: int | None = None
    sugerencias: list[str] = field(default_factory=list)
    tips: list[str] = field(default_factory=list)


@dataclass
class Result:
    """Resultado de ejecutar código Jqmas."""
    output: str
    errors: list[ErrorDetail]
    variables: dict


def _find_binary() -> str:
    """Busca el binario jqmas-core."""
    # 1. En el PATH
    path = shutil.which("jqmas-core")
    if path:
        return path

    # 2. En la raíz del repo (desarrollo)
    repo_root = Path(__file__).parent.parent.parent
    local = repo_root / "jqmas-core"
    if local.exists():
        return str(local)

    # 3. En ~/.local/bin
    local_bin = Path.home() / ".local" / "bin" / "jqmas-core"
    if local_bin.exists():
        return str(local_bin)

    raise EngineError(
        "No se encontró el binario 'jqmas-core'.\n"
        "       Compilalo con: ./build.sh\n"
        "       O instalalo con: make install"
    )


def run(source: str) -> Result:
    """Ejecuta código Jqmas y devuelve el resultado."""
    binary = _find_binary()

    try:
        proc = subprocess.run(
            [binary],
            input=json.dumps({"source": source}),
            capture_output=True,
            text=True,
            timeout=10,
        )
    except subprocess.TimeoutExpired:
        raise EngineError("El core tardó demasiado (timeout de 10s)")
    except FileNotFoundError:
        raise EngineError(f"No se pudo ejecutar '{binary}'")

    if not proc.stdout:
        raise EngineError(
            f"El core no devolvió nada.\n"
            f"       stderr: {proc.stderr.strip()}"
        )

    try:
        data = json.loads(proc.stdout)
    except json.JSONDecodeError as e:
        raise EngineError(f"Respuesta inválida del core: {e}\n{proc.stdout[:200]}")

    errors = [
        ErrorDetail(
            mensaje=e.get("mensaje", ""),
            linea=e.get("linea"),
            columna=e.get("columna"),
            sugerencias=e.get("sugerencias", []),
            tips=e.get("tips", []),
        )
        for e in data.get("errors", [])
    ]

    return Result(
        output=data.get("output", ""),
        errors=errors,
        variables=data.get("variables", {}),
    )

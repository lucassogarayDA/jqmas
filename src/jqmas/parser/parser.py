"""Parser de Jqmas — convierte tokens en AST."""

from jqmas.lexer.token import Token, TokenType
from jqmas.parser.ast_nodes import (
    Argument,
    BinaryOp,
    NumberLiteral,
    PresCall,
    Program,
    StringLiteral,
)


class ParserError(SyntaxError):
    """Error durante el análisis sintáctico."""


class Parser:
    """Parser de descenso recursivo para Jqmas."""

    def __init__(self, tokens: list[Token]):
        self.tokens = tokens
        self.pos = 0

    # ----- helpers -----

    def _peek(self) -> Token:
        return self.tokens[self.pos]

    def _previous(self) -> Token:
        return self.tokens[self.pos - 1]

    def _advance(self) -> Token:
        if not self._at_end():
            self.pos += 1
        return self._previous()

    def _at_end(self) -> bool:
        return self._peek().type == TokenType.EOF

    def _check(self, type_: TokenType) -> bool:
        return self._peek().type == type_

    def _match(self, *types: TokenType) -> bool:
        for t in types:
            if self._check(t):
                self._advance()
                return True
        return False

    def _consume(self, type_: TokenType, msg: str) -> Token:
        if self._check(type_):
            return self._advance()
        tok = self._peek()
        raise ParserError(
            f"{msg} (encontrado {tok.type.name} en línea {tok.line}, columna {tok.column})"
        )

    def _error(self, msg: str) -> None:
        tok = self._peek()
        raise ParserError(f"{msg} (línea {tok.line}, columna {tok.column})")

    # ----- parse -----

    def parse(self) -> Program:
        """Parsea el programa completo."""
        statements = []
        while not self._at_end():
            statements.append(self._statement())
        return Program(statements=statements)

    def _statement(self) -> object:
        """Una sentencia: por ahora solo pres[...]."""
        if self._check(TokenType.PRES):
            return self._pres_call()
        self._error("Se esperaba 'pres'")

    def _pres_call(self) -> PresCall:
        """Parsea pres[args...]."""
        self._consume(TokenType.PRES, "Se esperaba 'pres'")
        self._consume(TokenType.LBRACKET, "Se esperaba '[' después de pres")

        arguments: list[Argument] = []

        # Caso especial: pres[] vacío
        if not self._check(TokenType.RBRACKET):
            arguments.append(self._argument())
            while self._match(TokenType.COMMA):
                arguments.append(self._argument())

        self._consume(TokenType.RBRACKET, "Se esperaba ']' para cerrar pres")
        return PresCall(arguments=arguments)

    def _argument(self) -> Argument:
        """Parsea un argumento, con posible prefijo notN/."""
        not_count = 0

        # Detectar notN/ (not, opcionalmente número, slash)
        if self._check(TokenType.NOT):
            self._advance()  # consume 'not'

            # Número opcional entre not y /
            if self._check(TokenType.NUMBER):
                num_tok = self._advance()
                n = num_tok.value
                if not isinstance(n, int):
                    raise ParserError(
                        f"notN/ solo acepta número entero, no decimal (línea {num_tok.line})"
                    )
                if n < 1:
                    raise ParserError(
                        f"not0/ no es válido: el mínimo es 1 (línea {num_tok.line})"
                    )
                if n > 34:
                    raise ParserError(
                        f"not{n}/ excede el límite de 34 saltos (línea {num_tok.line})"
                    )
                not_count = n
            else:
                not_count = 1

            # El slash obligatorio
            self._consume(TokenType.SLASH, "Se esperaba '/' después de not")

            # Si lo que sigue es ] o , entonces notN/ va sin argumento: no hace nada
            if self._check(TokenType.RBRACKET) or self._check(TokenType.COMMA):
                return Argument(value=NumberLiteral(0), not_count=0)

        # Parsear la expresión del argumento
        value = self._expression()
        return Argument(value=value, not_count=not_count)

    # ----- expresiones -----

    def _expression(self) -> object:
        return self._addition()

    def _addition(self) -> object:
        left = self._multiplication()
        while self._match(TokenType.PLUS, TokenType.MINUS):
            op_tok = self._previous()
            op = "+" if op_tok.type == TokenType.PLUS else "-"
            right = self._multiplication()
            left = BinaryOp(left=left, operator=op, right=right)
        return left

    def _multiplication(self) -> object:
        left = self._unary()
        while self._match(TokenType.STAR, TokenType.SLASH):
            op_tok = self._previous()
            op = "*" if op_tok.type == TokenType.STAR else "/"
            right = self._unary()
            left = BinaryOp(left=left, operator=op, right=right)
        return left

    def _unary(self) -> object:
        # Acepta -expr como número negativo (útil)
        if self._match(TokenType.MINUS):
            operand = self._unary()
            return BinaryOp(
                left=NumberLiteral(0),
                operator="-",
                right=operand,
            )
        return self._primary()

    def _primary(self) -> object:
        if self._match(TokenType.NUMBER):
            return NumberLiteral(self._previous().value)
        if self._match(TokenType.STRING):
            return StringLiteral(self._previous().value)
        if self._match(TokenType.LPAREN):
            expr = self._expression()
            self._consume(TokenType.RPAREN, "Se esperaba ')'")
            return expr
        self._error("Se esperaba una expresión")

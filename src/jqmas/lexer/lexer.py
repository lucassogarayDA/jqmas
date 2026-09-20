"""Lexer de Jqmas — convierte código fuente en tokens."""

from jqmas.lexer.token import Token, TokenType


class LexerError(SyntaxError):
    """Error durante el análisis léxico."""


class Lexer:
    """Convierte texto de Jqmas en una lista de tokens."""

    KEYWORDS = {
        "pres": TokenType.PRES,
        "not": TokenType.NOT,
    }

    SINGLE_CHAR = {
        "+": TokenType.PLUS,
        "-": TokenType.MINUS,
        "*": TokenType.STAR,
        "/": TokenType.SLASH,
        "(": TokenType.LPAREN,
        ")": TokenType.RPAREN,
        "[": TokenType.LBRACKET,
        "]": TokenType.RBRACKET,
        ",": TokenType.COMMA,
    }

    def __init__(self, source: str):
        self.source = source
        self.pos = 0
        self.line = 1
        self.column = 1
        self.tokens: list[Token] = []

    def tokenize(self) -> list[Token]:
        """Devuelve la lista completa de tokens."""
        while not self._at_end():
            self._scan_token()
        self.tokens.append(Token(TokenType.EOF, None, self.line, self.column))
        return self.tokens

    # ----- helpers -----

    def _at_end(self) -> bool:
        return self.pos >= len(self.source)

    def _peek(self) -> str:
        if self._at_end():
            return "\0"
        return self.source[self.pos]

    def _peek_next(self) -> str:
        if self.pos + 1 >= len(self.source):
            return "\0"
        return self.source[self.pos + 1]

    def _advance(self) -> str:
        char = self.source[self.pos]
        self.pos += 1
        if char == "\n":
            self.line += 1
            self.column = 1
        else:
            self.column += 1
        return char

    def _add_token(self, type_: TokenType, value: object = None) -> None:
        self.tokens.append(Token(type_, value, self.line, self.column))

    def _error(self, msg: str) -> None:
        raise LexerError(f"{msg} (línea {self.line}, columna {self.column})")

    # ----- scan -----

    def _scan_token(self) -> None:
        char = self._peek()

        if char in " \t\r":
            self._advance()
            return

        if char == "\n":
            self._advance()
            return

        if char == "/" and self._peek_next() == "/":
            while not self._at_end() and self._peek() != "\n":
                self._advance()
            return

        if char == "'":
            self._scan_string()
            return

        if char.isdigit():
            self._scan_number()
            return

        if char.isalpha() or char == "_":
            self._scan_identifier()
            return

        if char in self.SINGLE_CHAR:
            self._advance()
            self._add_token(self.SINGLE_CHAR[char])
            return

        self._error(f"Carácter inesperado: {char!r}")

    def _scan_string(self) -> None:
        self._advance()
        result = []
        while not self._at_end() and self._peek() != "'":
            char = self._advance()
            if char == "\\":
                if self._at_end():
                    self._error("String sin cerrar")
                escaped = self._advance()
                if escaped == "n":
                    result.append("\n")
                elif escaped == "t":
                    result.append("\t")
                elif escaped == "\\":
                    result.append("\\")
                elif escaped == "'":
                    result.append("'")
                else:
                    result.append(escaped)
            else:
                result.append(char)

        if self._at_end():
            self._error("String sin cerrar")

        self._advance()
        self._add_token(TokenType.STRING, "".join(result))

    def _scan_number(self) -> None:
        start = self.pos
        while not self._at_end() and self._peek().isdigit():
            self._advance()

        if not self._at_end() and self._peek() == "." and self._peek_next().isdigit():
            self._advance()
            while not self._at_end() and self._peek().isdigit():
                self._advance()

        text = self.source[start:self.pos]
        if "." in text:
            value = float(text)
        else:
            value = int(text)
        self._add_token(TokenType.NUMBER, value)

    def _scan_identifier(self) -> None:
        start = self.pos
        start_column = self.column
        while not self._at_end() and (self._peek().isalnum() or self._peek() == "_"):
            self._advance()

        text = self.source[start:self.pos]

        # Caso especial: notN (not seguido de dígitos) → NOT + NUMBER
        if text.startswith("not") and len(text) > 3 and text[3:].isdigit():
            self.tokens.append(Token(TokenType.NOT, "not", self.line, start_column))
            numero = int(text[3:])
            self.tokens.append(
                Token(TokenType.NUMBER, numero, self.line, start_column + 3)
            )
            return

        type_ = self.KEYWORDS.get(text, TokenType.IDENTIFIER)
        self.tokens.append(Token(type_, text, self.line, start_column))

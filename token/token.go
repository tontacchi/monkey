package token

const (
	// Special Types
	ILLEGAL = "ILLEGAL"  // token / character not covered by lexer
	EOF     = "EOF"      // end of file (parser can stop)

	// Identifiers + Literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Balanced
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// type alias (change to enums later?)
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}



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

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

type Token struct {
	Type    TokenType
	Literal string
}

// Checks if identifier is in keywords map
func LookupIdentifier(identifier string) TokenType {
	if keyword, ok := keywords[identifier]; ok {
		return keyword
	}

	return IDENT
}


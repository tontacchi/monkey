package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int     // index of current char in input
	readPosition int     // index of char after current char in input
	char         byte    // current byte examined (pointed to by position)
}

//---[ Public Package Methods ]-------------------------------------------------
func New(input string) (newLexer *Lexer) {
	newLexer = &Lexer{
		input: input,
	}
	newLexer.readChar()

	return newLexer
}

//---[ Public Package Methods ]-------------------------------------------------


//---[ Lexer API Methods ]------------------------------------------------------
func (lex *Lexer) NextToken() (nextToken token.Token) {
	// Ignore whitespace
	lex.skipWhitespace()

	// Decide next token
	switch lex.char {
	case '=':
		nextToken = newToken(token.ASSIGN, lex.char)
	case ';':
		nextToken = newToken(token.SEMICOLON, lex.char)
	case '(':
		nextToken = newToken(token.LPAREN, lex.char)
	case ')':
		nextToken = newToken(token.RPAREN, lex.char)
	case ',':
		nextToken = newToken(token.COMMA, lex.char)
	case '+':
		nextToken = newToken(token.PLUS, lex.char)
	case '{':
		nextToken = newToken(token.LBRACE, lex.char)
	case '}':
		nextToken = newToken(token.RBRACE, lex.char)
	case '[':
		nextToken = newToken(token.LBRACKET, lex.char)
	case ']':
		nextToken = newToken(token.RBRACKET, lex.char)
	case 0:
		nextToken.Literal = ""
		nextToken.Type    = token.EOF
	default:
		if isLetter(lex.char) {
			nextToken.Literal = lex.readIdentifier()
			nextToken.Type    = token.LookupIdentifier(nextToken.Literal)
			return nextToken
		}

		nextToken = newToken(token.ILLEGAL, lex.char)
	}

	// Move the cursor to the next character
	lex.readChar()
	return nextToken
}

//---[ Lexer API Methods ]------------------------------------------------------


//---[ Lexer Helper Methods ]---------------------------------------------------
func (lex *Lexer) readChar() {
	// EOF reached
	if lex.readPosition >= len(lex.input) {
		lex.char = 0
		return
	}

	// read next char + ready readPos for next char
	lex.char = lex.input[lex.readPosition]

	lex.position = lex.readPosition
	lex.readPosition++
}

func (lex *Lexer) readIdentifier() string {
	start := lex.position

	for isLetter(lex.char) {
		lex.readChar()
	}

	until := lex.position 

	return lex.input[start:until]
}

func (lex *Lexer) skipWhitespace() {
	for lex.char == ' ' || lex.char == '\t' || lex.char == '\n' || lex.char == '\r' {
		lex.readChar()
	}
}

//---[ Lexer Helper Methods ]---------------------------------------------------


//---[ Package Helper Methods ]-------------------------------------------------
func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
	}
} 

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

//---[ Package Helper Methods ]-------------------------------------------------


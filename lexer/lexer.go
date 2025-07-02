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

func New(input string) (newLexer *Lexer) {
	newLexer = &Lexer{
		input: input,
	}
	newLexer.readChar()

	return newLexer
}

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

func (lex *Lexer) NextToken() (nextToken token.Token) {
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
		nextToken.Literal = "ILLEGAL"
		nextToken.Type    = token.ILLEGAL
	}

	lex.readChar()

	return nextToken
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
	}
} 


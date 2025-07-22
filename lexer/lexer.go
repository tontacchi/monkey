package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string  // Full source code
	position     int     // cursor's current index -> char's index
	readPosition int     // index after cursor's current index
	char         byte    // current byte examined (pointed to by position)
}

//---[ Public Package Methods ]-------------------------------------------------

func New(input string) (newLexer *Lexer) {
	newLexer = &Lexer{
		input: input,
	}

	// everything set to 0
	// important: readPosition & position both: 0
	newLexer.readChar()

	// invariants: (load char, sets up readPosition)
	// - char: input[0], position: 0
	// - readPosition: 1

	return newLexer
}

//---[ Public Package Methods ]-------------------------------------------------


//---[ Lexer API Methods ]------------------------------------------------------

func (lex *Lexer) NextToken() (nextToken token.Token) {
	// Ignore whitespace
	lex.skipWhitespace()

	// invariant: char: input[position] is an alphanum character

	// Decide next token
	switch lex.char {
	case '=':
		if lex.peekChar() == '=' {
			char := lex.char
			lex.readChar()

			nextToken.Type    = token.EQ
			nextToken.Literal = string(char) + string(lex.char)
		} else {
			nextToken = newToken(token.ASSIGN, lex.char)
		}
	case '+':
		nextToken = newToken(token.PLUS, lex.char)
	case '-':
		nextToken = newToken(token.MINUS, lex.char)
	case '!':
		if lex.peekChar() == '=' {
			currChar := lex.char
			lex.readChar()

			nextToken.Type    = token.NOT_EQ
			nextToken.Literal = string(currChar) + string(lex.char)
		} else {
			nextToken = newToken(token.BANG, lex.char)
		}
	case '/':
		nextToken = newToken(token.SLASH, lex.char)
	case '*':
		nextToken = newToken(token.ASTERISK, lex.char)
	case '<':
		nextToken = newToken(token.LT, lex.char)
	case '>':
		nextToken = newToken(token.GT, lex.char)
	case ';':
		nextToken = newToken(token.SEMICOLON, lex.char)
	case '(':
		nextToken = newToken(token.LPAREN, lex.char)
	case ')':
		nextToken = newToken(token.RPAREN, lex.char)
	case ',':
		nextToken = newToken(token.COMMA, lex.char)
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
			// char is letter / _ -> identifier (variable) or keyword
			nextToken.Literal = lex.readIdentifier()
			nextToken.Type    = token.LookupIdentifier(nextToken.Literal)

			// invariant: Token either a:
			// - variable (Type: IDENT)
			// - keyword  (Type: LET, FUNC, etc.)
			//   - Literal: text from input (varName, let, function)

			return nextToken
		} else if isDigit(lex.char) {
			// char is number -> automatically an int value
			nextToken.Literal = lex.readInt()
			nextToken.Type    = token.INT

			// invariant: Token is a:
			// - integer token (Type: INT)
			// - literal -> number as a string typed out in code

			return nextToken
		} else {
			// char not alphanum or other symbols -> Illegal
			nextToken = newToken(token.ILLEGAL, lex.char)

			// invariant: Token is a:
			// - unknown combo of chars (Type: ILLEGAL)
			// - literal -> garbled combo of characters
		}
	}

	// Move the cursor beyond end of current token
	lex.readChar()
	return nextToken
}

//---[ Lexer API Methods ]------------------------------------------------------


//---[ Lexer Helper Methods ]---------------------------------------------------

func (lex *Lexer) readChar() {
	// EOF / char harvesting control flow
	if lex.readPosition >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.readPosition]
	}

	// setup for next char to read
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

func (lex *Lexer) readInt() string {
	start := lex.position

	for isDigit(lex.char) {
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

func (lex *Lexer) peekChar() byte {
	if lex.readPosition >= len(lex.input) { return 0 }

	return lex.input[lex.readPosition]
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

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

//---[ Package Helper Methods ]-------------------------------------------------


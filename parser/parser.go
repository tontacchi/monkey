package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lex       *lexer.Lexer
	currToken token.Token
	peekToken token.Token
}


//---[ Module API Functions ]---------------------------------------------------

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lex: lex,
	}

	// sets currToken & peekToken
	parser.nextToken()
	parser.nextToken()

	return parser
}

//---[ Module API Functions ]---------------------------------------------------



//---[ Parser API Methods ]-----------------------------------------------------

func (parser *Parser) ParseProgram() *ast.Program {
	return nil
}

//---[ Parser API Methods ]-----------------------------------------------------



//---[ Parser Helper Methods ]--------------------------------------------------

func (parser *Parser) nextToken() {
	parser.currToken = parser.peekToken
	parser.peekToken = parser.lex.NextToken()
}

//---[ Parser Helper Methods ]--------------------------------------------------


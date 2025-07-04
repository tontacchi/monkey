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
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.currToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.nextToken()
	}

	return program
}

//---[ Parser API Methods ]-----------------------------------------------------



//---[ Parser Helper Methods ]--------------------------------------------------

func (parser *Parser) nextToken() {
	parser.currToken = parser.peekToken
	parser.peekToken = parser.lex.NextToken()
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: parser.currToken,
	}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: parser.currToken,
		Value: parser.currToken.Literal,
	}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skipping expressions until semicolon encountered
	for !parser.currTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) currTokenIs(tokenType token.TokenType) bool {
	return parser.currToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekTokenIs(tokenType) {
		parser.nextToken()
		return true
	}

	return false
}

//---[ Parser Helper Methods ]--------------------------------------------------


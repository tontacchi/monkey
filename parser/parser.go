package parser

import (
	"fmt"

	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

const (
	_ int = iota
	LOWEST      
	EQUALS       // ==
	LESSGREATER  // <, >
	SUM          // +
	PRODUCT      // *
	PREFIX       // -x, !x
	CALL         // f()
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lex       *lexer.Lexer
	errors    []string

	currToken token.Token
	peekToken token.Token

	prefixParseMap map[token.TokenType]prefixParseFn
	infixParseMap  map[token.TokenType]infixParseFn
}


//---[ Module API Functions ]---------------------------------------------------

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lex:    lex,
		errors: []string{},
	}

	parser.prefixParseMap = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)

	// sets currToken & peekToken
	parser.nextToken()
	parser.nextToken()

	return parser
}

//---[ Module API Functions ]---------------------------------------------------



//---[ Parser API Methods ]-----------------------------------------------------

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for parser.currToken.Type != token.EOF {
		statement := parser.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.nextToken()
	}

	return program
}

func (parser *Parser) Errors() []string {
	return parser.errors
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
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
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

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: parser.currToken,
	}

	parser.nextToken()

	// TODO: skipping expressions until semicolon encountered
	for !parser.currTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token: parser.currToken,
	}

	statement.Expression = parser.parseExpression(LOWEST)

	// optional semicolons -> easier REPL input
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}


func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefixFunc, ok := parser.prefixParseMap[parser.currToken.Type]
	if !ok {
		return nil
	}

	leftExp := prefixFunc()
	return leftExp
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: parser.currToken,
		Value: parser.currToken.Literal,
	}
}

// helpers for parseLetStatement()
func (parser *Parser) currTokenIs(tokenType token.TokenType) bool {
	return parser.currToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if !parser.peekTokenIs(tokenType) {
		parser.peekError(tokenType)
		return false
	}

	parser.nextToken()
	return true
}

// helper for Errors()
func (parser *Parser) peekError(tokenType token.TokenType) {
	message := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		tokenType,
		parser.peekToken.Type,
	)
	
	parser.errors = append(parser.errors, message)
}


// helper for token -> parse function maps setup
func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseMap[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseMap[tokenType] = fn
}

//---[ Parser Helper Methods ]--------------------------------------------------


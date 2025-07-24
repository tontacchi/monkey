package parser

import (
	"fmt"
	"strconv"

	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
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

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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

	// register tokens + associated parse functions
	parser.prefixParseMap = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT,   parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG,  parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

	parser.infixParseMap = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS,     parser.parseInfixExpression)
	parser.registerInfix(token.MINUS,    parser.parseInfixExpression)
	parser.registerInfix(token.SLASH,    parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ,       parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ,   parser.parseInfixExpression)
	parser.registerInfix(token.LT,       parser.parseInfixExpression)
	parser.registerInfix(token.GT,       parser.parseInfixExpression)

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
		parser.noPrefixParseFuncError(parser.currToken.Type)
		return nil
	}

	leftExp := prefixFunc()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infixFunc, ok := parser.infixParseMap[parser.peekToken.Type]
		if !ok {
			return leftExp
		}

		parser.nextToken()
		leftExp = infixFunc(leftExp)
	}

	return leftExp
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: parser.currToken,
		Value: parser.currToken.Literal,
	}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{
		Token: parser.currToken,
	}

	value, err := strconv.ParseInt(parser.currToken.Literal, 0, 64)
	if err != nil {
		errMsg := fmt.Sprintf("could not parse %q as int64", parser.currToken.Literal)
		parser.errors = append(parser.errors, errMsg)
		
		return nil
	}

	literal.Value = value

	return literal
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.currToken,
		Operator: parser.currToken.Literal,
	}

	parser.nextToken()
	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.currToken,
		Operator: parser.currToken.Literal,
		Left:     left,
	}

	precedence := parser.currPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
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


// helper for precedence
func (parser *Parser) currPrecedence() int {
	if precedence, ok := precedences[parser.currToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedences[parser.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
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


// parseExpression() helper for better error messages
func (parser *Parser) noPrefixParseFuncError(t token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %s found", t)
	parser.errors = append(parser.errors, message)
}

//---[ Parser Helper Methods ]--------------------------------------------------


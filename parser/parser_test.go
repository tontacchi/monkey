package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

//―――[ Main Tests ]―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedValue any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		parser := New(lexer.New(test.input))

		program := parser.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		numStatements := len(program.Statements)
		if numStatements != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, numStatements)
		}

		checkParserErrors(t, parser)

		statement := program.Statements[0]
		testLetStatement(t, statement, test.expectedIdent)

		// val := statement.(*ast.LetStatement).Value
		// testLiteralExpression(t, val, test.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	
	lex     := lexer.New(input)
	parser  := New(lex)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf(
			"program.Statements does not contain 3 statements. got=%d",
			len(program.Statements),
		)
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got=%T", statement)
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf(
				"returnStatement.TokenLiteral() not `return`, got %q",
				returnStatement.TokenLiteral(),
			)
		}
	}
}


func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lex     := lexer.New(input)
	parser  := New(lex)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program contains incorrect number of statements. got=%d",
			len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement) 
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf(
			"expression is not *ast.Identifier. got=%T",
			statement.Expression,
		)
	}

	if identifier.Value != "foobar" {
		t.Errorf(
			"identifier.Value not %s. got=%s",
			"foobar",
			identifier.Value,
		)
	}
	
	if identifier.TokenLiteral() != "foobar" {
		t.Errorf(
			"identifer.TokenLiteral is not %s. got=%s",
			"foobar",
			identifier.TokenLiteral(),
		)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lex     := lexer.New(input)
	parser  := New(lex)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program does not have 1 statement. got=%d",
			len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf(
			"expression not *ast.IntegerLiteral. got=%T",
			statement.Expression,
		)
	}

	if literal.Value != 5 {
		t.Errorf(
			"literal.Value not %d. got=%d",
			5,
			literal.Value,
		)
	}

	if literal.TokenLiteral() != "5" {
		 t.Errorf(
			 "literal.TokenLiteral not %s. got=%s",
			 "5",
			 literal.TokenLiteral(),
		 )
	}
}

//―――[ Main Tests ]―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――


//―――[ Helper Functions ]―――――――――――――――――――――――――――――――――――――――――――――――――――――――

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
	}
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	numErrors := len(errors)
	if numErrors == 0 { return }

	t.Errorf("parser has %d errors", numErrors)
	for _, errorMsg := range errors {
		t.Errorf("parser error: %q", errorMsg)
	}

	t.FailNow()
}

//―――[ Helper Functions ]―――――――――――――――――――――――――――――――――――――――――――――――――――――――


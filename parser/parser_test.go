package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedValue interface{}
	}{
		{"let x  5;", "x", 5},
		{"let  = true;", "y", true},
		{"let  25;", "foobar", "y"},
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
	for _, errMsg := range errors {
		t.Errorf("parser error: %q", errMsg)
	}

	t.FailNow()
}

package parser

import (
	"testing"
	"fmt"

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

func TestBooleanExpressions(t *testing.T) {
	tests := []struct{
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		lex     := lexer.New(test.input)
		parser  := New(lex)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := statement.Expression.(*ast.Boolean) 
		if !ok {
			t.Fatalf("expression not *ast.Boolean. got=%T", statement.Expression)
		}

		if boolean.Value != test.expected {
			t.Errorf("boolean.Value not %t. got=%t", test.expected,
				boolean.Value)
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct{
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range tests {
		lex     := lexer.New(test.input)
		parser  := New(lex)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statement. got=%d",
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

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(
				"statement.Expression is not *ast.PrefixExpression. got=%T",
				statement.Expression,
			)
		}

		if expression.Operator != test.operator {
			t.Fatalf(
				"expression.Operator is not '%s'. got='%s'",
				test.operator,
				expression.Operator,
			)
		}

		if !testLiteralExpression(t, expression.Right, test.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct{
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range tests {
		lex     := lexer.New(test.input)
		parser  := New(lex)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
		}

		if !testInfixExpression(t, statement.Expression, test.leftValue, test.operator, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct{
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, test := range tests {
		lex     := lexer.New(test.input)
		parser  := New(lex)
		program := parser.ParseProgram()

		for _, statement := range program.Statements {
			fmt.Printf("%v\n", statement)
		}

		checkParserErrors(t, parser)

		actual := program.String()
		fmt.Printf("%s -> %s\n", test.input, actual)

		if actual != test.expected {
			t.Errorf("expected=%q, got%q", test.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	lex     := lexer.New(input)
	parser  := New(lex)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain %d statements. got=%d\n",
			1,
			len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf(
			"statement.Expression is not ast.IfExpression. got=%T",
			statement.Expression,
		)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf(
			"consequence is not 1 statement. got=%d\n",
			expression.Consequence.Statements[0],
		)
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Consequence.Statements[0],
		)
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf(
			"expression.Alternative.Statements not nil. got=%+v",
			expression.Alternative,
		)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	lex     := lexer.New(input)
	parser  := New(lex)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain %d statements. got=%d\n",
			1,
			len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf(
			"statement.Expression is not ast.IfExpression. got=%T",
			statement.Expression,
		)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf(
			"consequence is not 1 statement. got=%d\n",
			expression.Consequence.Statements[0],
		)
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Consequence.Statements[0],
		)
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf(
			"alternative is not 1 statement. got=%d\n",
			expression.Alternative.Statements[0],
		)
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Alternative.Statements[0],
		)
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	lex     := lexer.New(input)
	parser  := New(lex)
	program := parser.ParseProgram()

	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain 1 statements. got=%d\n",
			len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf(
			"statement.Expression is not ast.FunctionLiteral. got=%T",
			statement.Expression,
		)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf(
			"function literal parameters wrong. got=%d\n",
			len(function.Parameters),
		)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf(
			"function.Body.Statements does not have 1 statement. got=%d\n",
			len(function.Body.Statements),
		)
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"function body statement is not ast.ExpressionStatement, got=%T",
			function.Body.Statements[0],
		)
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct{
		input    string
		expected []string
	}{
		{input: "fn() {}",        expected: []string{}},
		{input: "fn(x) {}",       expected: []string{"x"}},
		{input: "fn(x, y, z) {}", expected: []string{"x", "y", "z"}},
	}

	for _, test := range tests {
		lex     := lexer.New(test.input)
		parser  := New(lex)
		program := parser.ParseProgram()

		checkParserErrors(t, parser)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function  := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(test.expected) {
			t.Errorf(
				"length parameters wrong. want %d, got=%d\n",
				len(test.expected),
				len(function.Parameters),
			)
		}

		for i, identifier := range test.expected {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}
	}
}

//―――[ Main Tests ]―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――


//―――[ Helper Functions ]―――――――――――――――――――――――――――――――――――――――――――――――――――――――

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


func testLiteralExpression(
	t          *testing.T,
	expression ast.Expression,
	expected   any,
) bool {
	switch castedValue := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(castedValue))
	case int64:
		return testIntegerLiteral(t, expression, castedValue)
	case bool:
		return testBooleanLiteral(t, expression, castedValue)
	case string:
		return testIdentifier(t, expression, castedValue)
	}

	t.Errorf("type of expression not handled. got=%T", expression)
	return false
}

func testIntegerLiteral(t *testing.T, intExp ast.Expression, value int64) bool {
	intLit, ok := intExp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf(
			"intExp is not *ast.IntegerLiteral. got=%T",
			intExp,
		)
		return false
	}

	if intLit.Value != value {
		t.Errorf(
			"intLit value is not %d. got =%d",
			value,
			intLit.Value,
		)
		return false
	}

	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf(
			"intLit.TokenLiteral is not %d. got=%s",
			value,
			intLit.TokenLiteral(),
		)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expression not *ast.Boolean. got=%T", expression)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t, got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
	}

	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression not *ast.Identifier. got=%T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value not %s. got=%s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() not %s. got=%s",
			value,
			identifier.TokenLiteral(),
		)
		return false
	}

	return true
}

func testInfixExpression(
	t          *testing.T,
	expression ast.Expression,
	left       any,
	operator   string,
	right      any,
) bool {
	infix, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf(
			"expression is not ast.InfixExpression. got=%T(%s)",
			expression,
			expression,
		)
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}

	if infix.Operator != operator {
		t.Errorf("infix.Operator is not '%s'. got=%q", operator, infix.Operator)
		return false
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
} 

//―――[ Helper Functions ]―――――――――――――――――――――――――――――――――――――――――――――――――――――――


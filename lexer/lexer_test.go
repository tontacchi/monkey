package lexer

import (
	"testing"

	"monkey/token"
)

func TestNextTokenBasicSymbols(t *testing.T) {
	input := "=+(){}[],;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS,   "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.COMMA,  ","},
		{token.SEMICOLON, ";"},
		{token.EOF,    ""},
	}

	lex := New(input)

	for index, test := range tests {
		testToken := lex.NextToken()

		if test.expectedType != testToken.Type {
			t.Fatalf("tests[%d] - incorrect token type. expected=%v, got=%v",
				index, test.expectedType, testToken.Type,
			)
		}

		if testToken.Literal != test.expectedLiteral {
			t.Fatalf("test[%d] - incorrect literal. expected=%v, got=%v",
				index, test.expectedLiteral, testToken.Literal,
			)
		}
	}
}

func TestNextTokenSubsetInput(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// fill in from provided repo
		{token.EOF, ""},
	}

	lex := New(input)

	for index, test := range tests {
		testToken := lex.NextToken()

		if test.expectedType != testToken.Type {
			t.Fatalf("tests[%d] - incorrect token type. expected=%v, got=%v",
				index, test.expectedType, testToken.Type,
			)
		}

		if testToken.Literal != test.expectedLiteral {
			t.Fatalf("test[%d] - incorrect literal. expected=%v, got=%v",
				index, test.expectedLiteral, testToken.Literal,
			)
		}
	}
}

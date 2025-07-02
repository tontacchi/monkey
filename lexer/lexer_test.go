package lexer

import (
	"testing"

	"monkey/token"
)

func TestNextToken(t *testing.T) {
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

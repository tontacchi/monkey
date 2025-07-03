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
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
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

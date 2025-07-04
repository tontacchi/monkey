package ast

import (
	"monkey/token"
)

//---[ Node Interfaces ]--------------------------------------------------------
type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

//---[ Node Interfaces ]--------------------------------------------------------


type Program struct {
	Statements []Statement
}

// Program now implements the Node interface
func (program *Program) TokenLiteral() string {
	if len(program.Statements) <= 0 {
		return ""
	}

	return program.Statements[0].TokenLiteral()
}


type Identifier struct {
	Token token.Token  // token.IDENT
	Value string
}

func (identifier *Identifier) expressionNode() {}  // simplifies using identifiers as RHS values

func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}


type LetStatement struct {
	Token  token.Token  // should always be the token.LET token
	Name  *Identifier   // variable used in binding
	Value  Expression   // RHS produces value -> bind to variable
}

func (let *LetStatement) statementNode() {}

func (let *LetStatement) TokenLiteral() string {
	return let.Token.Literal
}


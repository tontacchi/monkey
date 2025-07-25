package ast

import (
	"monkey/token"
	"bytes"
)

//---[ Node Interfaces ]--------------------------------------------------------
type Node interface {
	TokenLiteral() string
	String()       string
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

func (program *Program) TokenLiteral() string {
	if len(program.Statements) <= 0 {
		return ""
	}

	return program.Statements[0].TokenLiteral()
}

func (program *Program) String() string {
	var buffer bytes.Buffer

	for _, statement := range program.Statements {
		buffer.WriteString(statement.String())
	}

	return buffer.String()
}


type Identifier struct {
	Token token.Token  // token.IDENT
	Value string
}

func (identifier *Identifier) expressionNode() {}  // simplifies using identifiers as RHS values

func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}

func (identifier *Identifier) String() string {
	return identifier.Value
}


type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}


type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
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

func (let *LetStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(let.TokenLiteral() + " ")
	buffer.WriteString(let.Name.String())
	buffer.WriteString(" = ")

	if let.Value != nil {
		buffer.WriteString(let.Value.String())
	}

	buffer.WriteString(";")
	return buffer.String()
}


type ReturnStatement struct {
	Token       token.Token  // should always be token.RETURN
	ReturnValue Expression   // produces value to give to caller
}

func (ret *ReturnStatement) statementNode() {}

func (ret *ReturnStatement) TokenLiteral() string {
	return ret.Token.Literal
}

func (ret *ReturnStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(ret.String() + " ")

	if ret.ReturnValue != nil {
		buffer.WriteString(ret.ReturnValue.String())
	}

	buffer.WriteString(";")
	return buffer.String()
}


type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (exp *ExpressionStatement) statementNode() {}

func (exp *ExpressionStatement) TokenLiteral() string {
	return exp.Token.Literal
}

func (exp *ExpressionStatement) String() string {
	if exp.Expression != nil {
		return exp.Expression.String()
	}

	return ""
}



type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (prefix *PrefixExpression) expressionNode() {}

func (prefix *PrefixExpression) TokenLiteral() string {
	return prefix.Token.Literal
}

func (prefix *PrefixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(" + prefix.Operator + prefix.Right.String() + ")")

	return buffer.String()
}


type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (infix *InfixExpression) expressionNode() {}

func (infix *InfixExpression) TokenLiteral() string {
	return infix.Token.Literal
}

func (infix *InfixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	buffer.WriteString(infix.Left.String())
	buffer.WriteString(" " + infix.Operator + " ")
	buffer.WriteString(infix.Right.String())
	buffer.WriteString(")")

	return buffer.String()
}


type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (block *BlockStatement) statementNode() {}

func (block *BlockStatement) TokenLiteral() string {
	return block.Token.Literal
}

func (block *BlockStatement) String() string {
	var buffer bytes.Buffer

	for _, statement := range block.Statements {
		buffer.WriteString(statement.String())
	}

	return buffer.String()
}


type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExp *IfExpression) expressionNode() {}

func (ifExp *IfExpression) TokenLiteral() string {
	return ifExp.Token.Literal
}

func (ifExp *IfExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("if")
	buffer.WriteString(ifExp.Condition.String())
	buffer.WriteString(" ")
	buffer.WriteString(ifExp.Consequence.String())

	if ifExp.Alternative != nil {
		buffer.WriteString("else")
		buffer.WriteString(ifExp.Alternative.String())
	}

	return buffer.String()
}


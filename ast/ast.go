package ast

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

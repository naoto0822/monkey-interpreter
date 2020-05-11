package ast

import (
	"bytes"

	"github.com/naoto0822/monkey-interpreter/pkg/token"
)

// Node is interface
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is interface
type Statement interface {
	Node
	statementNode()
}

// Expression is interface
type Expression interface {
	Node
	expressionNode()
}

var _ (Node) = (*Program)(nil)

// Program is root node
type Program struct {
	Statements []Statement
}

// TokenLiteral implements Node
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// String implements Node
func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

var _ Statement = (*LetStatement)(nil)

// LetStatement ex. let x = 5;
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statementNode() {}

// TokenLiteral implements Statement
func (s *LetStatement) TokenLiteral() string {
	return s.Token.Literal
}

// String implements Statement
func (s *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	out.WriteString(" = ")

	if s.Value != nil {
		out.WriteString(s.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

var _ Expression = (*Identifier)(nil)

// Identifier is name of let
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral implements Expression
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String implements Expression
func (i *Identifier) String() string {
	return i.Value
}

var _ Statement = (*ReturnStatement)(nil)

// ReturnStatement ex return expression;
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (s *ReturnStatement) statementNode() {}

// TokenLiteral implemets Statement
func (s *ReturnStatement) TokenLiteral() string {
	return s.Token.Literal
}

// String implements Statement
func (s *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral() + " ")

	if s.ReturnValue != nil {
		out.WriteString(s.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

var _ Statement = (*ExpressionStatement)(nil)

// ExpressionStatement is ex x + 10;
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (s *ExpressionStatement) statementNode() {}

// TokenLiteral implements Statement
func (s *ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}

// String implements Statement
func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}

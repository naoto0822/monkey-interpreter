package ast

import (
	"github.com/naoto0822/monkey-interpreter/pkg/token"
)

// Node is interface
type Node interface {
	TokenLiteral() string
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

var _ Expression = (*Identifier)(nil)

// Identifier is name of let
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral implements Node
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

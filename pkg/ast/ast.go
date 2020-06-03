package ast

import (
	"bytes"
	"strings"

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

var _ Expression = (*IntegerLiteral)(nil)

// IntegerLiteral is ex 5
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

// TokenLiteral implmenets Expression
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// String implements Expression
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

var _ Expression = (*PrefixExpression)(nil)

// PrefixExpression is <prefix operator><exp>;
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}

// TokenLiteral implements Expression
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// String implements Expression
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

var _ (Expression) = (*InfixExpression)(nil)

// InfixExpression is <exp><operator><exp>;
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

// TokenLiteral implements Expression
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String implements Expression
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

var _ Expression = (*Boolean)(nil)

// Boolean is true or false
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral implements Expression
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String implements Expression
func (b *Boolean) String() string {
	return b.Token.Literal
}

var _ Expression = (*IfExpression)(nil)

// IfExpression is if else ~
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}

// TokenLiteral implements Expression
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String implements Expression
func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

var _ Statement = (*BlockStatement)(nil)

// BlockStatement is block { }
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}

// TokenLiteral implements Statement
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

// String implements Statement
func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

var _ Expression = (*FunctionLiteral)(nil)

// FunctionLiteral is fn(x, y){ return x; }
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}

// TokenLiteral implements Expression
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

// String implements Expression
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")

	if len(f.Parameters) > 0 {
		params := []string{}

		for _, p := range f.Parameters {
			params = append(params, p.String())
		}

		out.WriteString(strings.Join(params, ", "))
	}

	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

// CallExpression is add(2, 3)
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}

// TokenLiteral implements Expression
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

// String implements Expression
func (c *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(c.Function.String())
	out.WriteString("(")

	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

var _ Expression = (*StringLiteral)(nil)

// StringLiteral is "hoge"
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {}

// TokenLiteral implements Expression
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

// String implements Expression
func (s *StringLiteral) String() string {
	return s.Token.Literal
}

var _ Expression = (*ArrayLiteral)(nil)

// ArrayLiteral is [a, b,,,]
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}

// TokenLiteral implments Expression
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

// String implements Expression
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

var _ Expression = (*IndexExpression)(nil)

// IndexExpression is exp[exp]
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {}

// TokenLiteral implements Expression
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String implements Expression
func (i *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")

	return out.String()
}

// HashLiteral is {k:v}
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode() {}

// TokenLiteral implements Expression
func (h *HashLiteral) TokenLiteral() string {
	return h.Token.Literal
}

// String implements Expression
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for k, v := range h.Pairs {
		pair := k.String() + ":" + v.String()
		pairs = append(pairs, pair)
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

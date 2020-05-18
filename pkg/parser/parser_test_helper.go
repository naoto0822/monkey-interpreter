package parser

import (
	"fmt"
	"testing"

	"github.com/naoto0822/monkey-interpreter/pkg/ast"
)

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not let. got=%s", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s, got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not %s, got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, a ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, a, int64(v))
	case int64:
		return testIntegerLiteral(t, a, v)
	case string:
		return testIdentifier(t, a, v)
	case bool:
		return testBooleanLiteral(t, a, v)
	}

	t.Errorf("type of a not handled. got=%T", a)
	return false
}

func testIntegerLiteral(t *testing.T, a ast.Expression, value int64) bool {
	il, ok := a.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not ast.InetgerLiteral. got=%T", a)
		return false
	}

	if il.Value != value {
		t.Errorf("il.Value not %d, got=%d", value, il.Value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("il.TokenLiteral not %d, got=%s", value, il.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, a ast.Expression, value string) bool {
	id, ok := a.(*ast.Identifier)
	if !ok {
		t.Errorf("id is not ast.Identifier. got=%T", id)
		return false
	}

	if id.Value != value {
		t.Errorf("id.Value is not %s. got=%s", value, id.Value)
		return false
	}

	if id.TokenLiteral() != value {
		t.Errorf("id.TokenLiteral() is not %s. got=%s", value, id.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, a ast.Expression, value bool) bool {
	b, ok := a.(*ast.Boolean)
	if !ok {
		t.Errorf("a is not ast.Boolean. got=%T", a)
		return false
	}

	if b.Value != value {
		t.Errorf("b.Value is not %t. got=%t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral is not %t. got=%s", value, b.TokenLiteral())
		return false
	}

	return true
}

func testPrefixExpression(
	t *testing.T, a ast.Expression, operator string, right interface{}) bool {

	exp, ok := a.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("a is not ast.PrefixExpression. got=%T", a)
		return false
	}

	if exp.Operator != operator {
		t.Errorf("exp.Operator is not %s. got=%s", operator, exp.Operator)
		return false
	}

	if !testLiteralExpression(t, exp.Right, right) {
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	a ast.Expression,
	left interface{},
	operator string,
	right interface{}) bool {

	exp, ok := a.(*ast.InfixExpression)
	if !ok {
		t.Errorf("a is not ast.InfixExpression. got=%T", a)
		return false
	}

	if !testLiteralExpression(t, exp.Left, left) {
		return false
	}

	if exp.Operator != operator {
		t.Errorf("exp.Opearator is not %s. got=%s", operator, exp.Operator)
		return false
	}

	if !testLiteralExpression(t, exp.Right, right) {
		return false
	}

	return true
}

package evaluator

import (
	"testing"

	"github.com/naoto0822/monkey-interpreter/pkg/lexer"
	"github.com/naoto0822/monkey-interpreter/pkg/object"
	"github.com/naoto0822/monkey-interpreter/pkg/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	intObj, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not object.Integer. got=%T", obj)
		return false
	}

	if intObj.Value != expected {
		t.Errorf("intObj.Value is not %d. got=%d", expected, intObj.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	boolObj, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj is not object.Boolean. got=%T", obj)
		return false
	}

	if boolObj.Value != expected {
		t.Errorf("boolObj.Value is not %t. got=%t", expected, boolObj.Value)
		return false
	}

	return true

}

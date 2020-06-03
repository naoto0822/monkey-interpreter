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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 > 1", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"1 == 2", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func TestEvalIfExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, obj, int64(integer))
		} else {
			testNullObject(t, obj)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if(10 > 1){if(10 > 1){return 10;}return 1;}", 10},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {

	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1){ true + false }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)

		e, ok := obj.(*object.Error)
		if !ok {
			t.Errorf("e is not object.Error. got=%T", obj)
			continue
		}

		if e.Message != tt.expectedMessage {
			t.Errorf("e.Message is not %s. got=%s", tt.expectedMessage, e.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)

		if !testIntegerObject(t, obj, tt.expected) {
			continue
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x + 2; };`
	obj := testEval(input)

	fn, ok := obj.(*object.Function)
	if !ok {
		t.Errorf("obj is not object.Function. got=%T", obj)
		return
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("fn.Parameters has wrong parameters. got=%d",
			len(fn.Parameters))
		return
	}

	if fn.Parameters[0].String() != "x" {
		t.Errorf("p1 is not x. got=%s", fn.Parameters[0].String())
		return
	}

	if fn.Body.String() != "(x + 2)" {
		t.Errorf("fn.Body is not (x + 2). got=%s", fn.Body.String())
		return
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let identifier = fn(x) { x; }; identifier(5);",
			5,
		},
		{
			"let identifier = fn(x) { return x; }; identifier(5);",
			5,
		},
		{
			"let double = fn(x) { return x * 2; }; double(5);",
			10,
		},
		{
			"let add = fn(x, y) { return x + y; }; add(5, 5);",
			10,
		},
		{
			"let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));",
			20,
		},
		{
			"fn(x){x;}(5)", 5,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello world";`
	obj := testEval(input)

	stringObj, ok := obj.(*object.String)
	if !ok {
		t.Errorf("obj is not object.String. got=%T", obj)
		return
	}

	if stringObj.Value != "Hello world" {
		t.Errorf("stringObj.Value is not %s. got=%s", "Hello world", stringObj.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	obj := testEval(input)

	stringObj, ok := obj.(*object.String)
	if !ok {
		t.Errorf("obj is not object.String. got=%T", obj)
		return
	}

	if stringObj.Value != "Hello World!" {
		t.Errorf("stringObj.Value is not %s. got=%s", "Hello World!", stringObj.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`len("")`, 0,
		},
		{
			`len("four")`, 4,
		},
		{
			`len("hello world")`, 11,
		},
		{
			`len(1)`, "argument to `len` not supported, got INTEGER",
		},
		{
			`len("one", "two")`, "wrong number of arguments. got=2, want=1",
		},
		{
			`len([1, 2])`, 2,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, obj, int64(expected))
		case string:
			errObj, ok := obj.(*object.Error)
			if !ok {
				t.Errorf("obj is not object.Error. got=%T", obj)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("errObj.Message is not %s. got=%s",
					expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`
	obj := testEval(input)

	array, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("obj is not object.Array. got=%T", obj)
		return
	}

	if len(array.Elements) != 3 {
		t.Errorf("array.Elements does not contain 3. got=%d", len(array.Elements))
		return
	}

	if !testIntegerObject(t, array.Elements[0], 1) {
		return
	}

	if !testIntegerObject(t, array.Elements[1], 4) {
		return
	}

	if !testIntegerObject(t, array.Elements[2], 6) {
		return
	}
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, obj, int64(integer))
		} else {
			testNullObject(t, obj)
		}
	}
}

func TestParseHashLiteral(t *testing.T) {
	input := `let two = "two";

{
 "one": 10 - 9,
 two: 1 + 1,
 "thr" + "ee": 6 / 2,
 4: 4,
 true: 5,
 false: 6
}`
	obj := testEval(input)
	hash, ok := obj.(*object.Hash)
	if !ok {
		t.Errorf("obj is not object.Hash. got=%T", obj)
		return
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(hash.Pairs) != 6 {
		t.Errorf("hash.Pairs does not contain 6. got=%d", len(hash.Pairs))
		return
	}

	for k, v := range expected {
		pair, ok := hash.Pairs[k]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
			continue
		}

		testIntegerObject(t, pair.Value, v)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`, 5,
		},
		{
			`{"foo": 5}["bar"]`, nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		expectedInt, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, obj, int64(expectedInt))
		} else {
			testNullObject(t, obj)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("obj is not NULL. got=%T", obj)
		return false
	}

	return true
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

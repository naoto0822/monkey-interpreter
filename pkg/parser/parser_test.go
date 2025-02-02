package parser

import (
	"testing"

	"github.com/naoto0822/monkey-interpreter/pkg/ast"
	"github.com/naoto0822/monkey-interpreter/pkg/lexer"
)

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Errorf("program.Statements is not ast.LetStatement")
		}

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		if !testLiteralExpression(t, stmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("program.Statements is not ast.LetStatement")
		}

		if stmt.TokenLiteral() != "return" {
			t.Errorf("stmt.TokenLiteral not return. got=%T", stmt.TokenLiteral())
		}

		if !testLiteralExpression(t, stmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program has not enough statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if !testIdentifier(t, stmt.Expression, "foobar") {
		return
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}

func TestParsePrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testPrefixExpression(t, stmt.Expression, tt.operator, tt.value) {
			return
		}
	}
}

func TestParseInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statement does not contain 1. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a + b",
			"((-a) + b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if program.String() != tt.expected {
			t.Errorf("program.String is not %s. got=%s",
				tt.expected, program.String())
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statement does not contain 1. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tt.value) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements is not ast.ExpressionStatement")
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("stmt.Expression is not ast.IfExpression")
	}

	if !testInfixExpression(t, ifExp.Condition, "x", "<", "y") {
		return
	}

	if len(ifExp.Consequence.Statements) != 1 {
		t.Errorf("ifExp.Consequence.Statement does not contain. got=%d",
			len(ifExp.Consequence.Statements))
	}

	blockStmt, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("ifExp.Consequence.Statements is not ast.ExpressionStatement")
	}

	if !testIdentifier(t, blockStmt.Expression, "x") {
		return
	}

	if ifExp.Alternative != nil {
		t.Errorf("ifExp.Alternative was not nil. got=%+v", ifExp.Alternative)
	}
}

func TestFuntionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	fnExp, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(fnExp.Parameters) != 2 {
		t.Errorf("fnExp.Parameters does not contain 2. got=%d", len(fnExp.Parameters))
	}

	if !testIdentifier(t, fnExp.Parameters[0], "x") {
		return
	}

	if !testIdentifier(t, fnExp.Parameters[1], "y") {
		return
	}

	if len(fnExp.Body.Statements) != 1 {
		t.Errorf("fnExp.Body.Statements does not contain 1. got=%d",
			len(fnExp.Body.Statements))
	}

	bStmt, ok := fnExp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("fnExp.Body.Statements is not ast.ExpressionStatement. got=%T",
			fnExp.Body.Statements[0])
	}

	if !testInfixExpression(t, bStmt.Expression, "x", "+", "y") {
		return
	}
}

func TestParseFunctionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn(){};", []string{}},
		{"fn(x){};", []string{"x"}},
		{"fn(x, y, z){};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements is not ast.ExpressionStatement")
		}

		fnExp, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Errorf("stmt.Expression is not ast.FunctionLiteral")
		}

		if len(fnExp.Parameters) != len(tt.expectedParams) {
			t.Errorf("fnExp.Parameters does not contain %d. got=%d",
				len(tt.expectedParams), len(fnExp.Parameters))
		}

		for i, ev := range tt.expectedParams {
			testIdentifier(t, fnExp.Parameters[i], ev)
		}
	}
}

func TestParseCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statement is not ast.ExpressionStatement")
	}

	call, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Errorf("stmt.Expression is not ast.CallExpression")
	}

	if !testIdentifier(t, call.Function, "add") {
		return
	}

	if len(call.Arguments) != 3 {
		t.Errorf("call.Arguments does not contain 3. git=%d",
			len(call.Arguments))
	}

	if !testLiteralExpression(t, call.Arguments[0], 1) {
		return
	}

	if !testInfixExpression(t, call.Arguments[1], 2, "*", 3) {
		return
	}

	if !testInfixExpression(t, call.Arguments[2], 4, "+", 5) {
		return
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1. got=%d",
			len(program.Statements))
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
		return
	}

	exp, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not StringLiteral. got=%T", stmt.Expression)
		return
	}

	if exp.Value != "hello world" {
		t.Errorf("exp.Value is not %s. got=%s", "hello world", exp.Value)
	}
}

func TestParseArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1. got=%d",
			len(program.Statements))
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
		return
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not ast.ArrayLiteral. got=%T", stmt.Expression)
		return
	}

	if len(array.Elements) != 3 {
		t.Errorf("array.Elements does not contain 3. got=%d", len(array.Elements))
		return
	}

	if !testIntegerLiteral(t, array.Elements[0], 1) {
		return
	}

	if !testInfixExpression(t, array.Elements[1], 2, "*", 2) {
		return
	}

	if !testInfixExpression(t, array.Elements[2], 3, "+", 3) {
		return
	}
}

func TestParseIndexExpression(t *testing.T) {
	input := `myArray[1 + 1]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1. got=%d",
			len(program.Statements))
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
		return
	}

	index, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Errorf("stmt.Expression is not ast.IndexExpression. got=%T", stmt.Expression)
		return
	}

	if !testIdentifier(t, index.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, index.Index, 1, "+", 1) {
		return
	}
}

func TestParseHashLiteralStringKey(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1. got=%d",
			len(program.Statements))
		return
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
		return
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not ast.HashLiteral. got=%T", stmt.Expression)
		return
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", k)
			continue
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, v, expectedValue)
	}
}

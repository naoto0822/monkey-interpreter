package evaluator

import (
	"github.com/naoto0822/monkey-interpreter/pkg/ast"
	"github.com/naoto0822/monkey-interpreter/pkg/object"
)

// Eval start parsing ast.Node
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// statement
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

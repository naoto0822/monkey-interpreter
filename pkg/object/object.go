package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/naoto0822/monkey-interpreter/pkg/ast"
)

// Type is object type
type Type string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE_OBJ"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
)

// Object monkey value
type Object interface {
	Type() Type
	Inspect() string
}

var _ Object = (*Integer)(nil)

// Integer is exp int
type Integer struct {
	Value int64
}

// Type implements Type
func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

// Inspect implements Object
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

var _ Object = (*Boolean)(nil)

// Boolean is exp bool
type Boolean struct {
	Value bool
}

// Type implements Object
func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

// Inspect implements Object
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

var _ Object = (*Null)(nil)

// Null is null
type Null struct{}

// Type implements Object
func (n *Null) Type() Type {
	return NULL_OBJ
}

// Inspect implements Object
func (n *Null) Inspect() string {
	return "null"
}

var _ Object = (*ReturnValue)(nil)

// ReturnValue is return
type ReturnValue struct {
	Value Object
}

// Type implements Object
func (r *ReturnValue) Type() Type {
	return RETURN_VALUE_OBJ
}

// Inspect implements Object
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

var _ Object = (*Error)(nil)

// Error is error
type Error struct {
	Message string
}

// Type implements Object
func (e *Error) Type() Type {
	return ERROR_OBJ
}

// Inspect implements Object
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// Environment has let identifier
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment gen Environment
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

// NewEnclosedEnvironment gen Env w/ outer
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get is get object
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set is set object w/ name
func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

var _ Object = (*Function)(nil)

// Function is fn()
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type implements Object
func (f *Function) Type() Type {
	return FUNCTION_OBJ
}

// Inspect implements Object
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

var _ Object = (*String)(nil)

// String is string!
type String struct {
	Value string
}

// Type implements Object
func (s *String) Type() Type {
	return STRING_OBJ
}

// Inspect implements Object
func (s *String) Inspect() string {
	return s.Value
}

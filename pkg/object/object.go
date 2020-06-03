package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

// Object monkey value
type Object interface {
	Type() Type
	Inspect() string
}

// Hashable hash!
type Hashable interface {
	HashKey() HashKey
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

// HashKey gen hash key
func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
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

// HashKey gen hash key
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{
		Type:  b.Type(),
		Value: value,
	}
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

// HashKey gen hash key
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}

// BuiltinFunction return Object
type BuiltinFunction func(args ...Object) Object

var _ Object = (*Builtin)(nil)

// Builtin wrap BuiltinFunction
type Builtin struct {
	Fn BuiltinFunction
}

// Type implements Object
func (b *Builtin) Type() Type {
	return BUILTIN_OBJ
}

// Inspect implements Object
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Array is [a, b, c]
type Array struct {
	Elements []Object
}

// Type implements Object
func (a *Array) Type() Type {
	return ARRAY_OBJ
}

// Inspect implements Object
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// HashKey is hash key
type HashKey struct {
	Type  Type
	Value uint64
}

// HashPair is key and value
type HashPair struct {
	Key   Object
	Value Object
}

// Hash is {k:v}
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type implements Object
func (h *Hash) Type() Type {
	return HASH_OBJ
}

// Inspect implements Object
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, p := range h.Pairs {
		pair := fmt.Sprintf("%s: %s", p.Key.Inspect(), p.Value.Inspect())
		pairs = append(pairs, pair)
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

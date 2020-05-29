package object

import (
	"fmt"
)

// Type is object type
type Type string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE_OBJ"
	ERROR_OBJ        = "ERROR"
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
}

// NewEnvironment gen Environment
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

// Get is get object
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set is set object w/ name
func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

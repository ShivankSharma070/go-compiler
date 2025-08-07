package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/ShivankSharma070/go-compiler/ast"
	"github.com/ShivankSharma070/go-compiler/code"
)

const (
	INTEGER_OBJ           = "INTEGER"
	BOOLEAN_OBJ           = "BOOLEAN"
	RETURN_VALUE_OBJ      = "RETURN_VALUE"
	NULL_OBJ              = "NULL"
	ERROR_OBJ             = "ERROR"
	FUNCTION_OBJ          = "FUNCTION"
	STRING_OBJ            = "STRING"
	BUILTIN_OBJ           = "BUILTIN"
	ARRAY_OBJ             = "ARRAY"
	HASH_OBJ              = "HASH"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Interface to check if a object is hashable or not.
type Hashable interface {
	// TODO: Improve performance of HashKey() by caching the return values
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Inspect() string  { return "NULL" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "Error: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type FunctionLiteral struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fe *FunctionLiteral) Inspect() string {
	var result bytes.Buffer
	params := []string{}
	for _, f := range fe.Parameters {
		params = append(params, f.String())
	}

	result.WriteString("fn (")
	result.WriteString(strings.Join(params, ", "))
	result.WriteString(") {\n")
	result.WriteString(fe.Body.String())
	result.WriteString("\n}")
	return result.String()
}

func (fe *FunctionLiteral) Type() ObjectType {
	return FUNCTION_OBJ
}

// ============ ENVIRONMENT ==============
type Environment struct {
	Store map[string]Object
	Outer *Environment
}

func NewEnclosingEnvironment(enclosingEnv *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = enclosingEnv
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{Store: s, Outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.Store[name]
	if !ok && e.Outer != nil {
		value, ok = e.Outer.Get(name)
	}
	return value, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.Store[name] = value
	return value
}

// ================== BUILT-IN FUNCTION ===================

type BuiltInFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltInFunction
}

func (bu *Builtin) Inspect() string  { return "builtin function" }
func (bu *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// =================== ARRAY =========================
type Array struct {
	Elements []Object
}

func (ar *Array) Type() ObjectType { return ARRAY_OBJ }
func (ar *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ar.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// ==================== HASH ============================
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pair map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pair {
		pairs = append(pairs, fmt.Sprintf("%s : %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// =================== COMPILED FUNCTION ====================
type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals int
}

func (cf *CompiledFunction) Type() ObjectType {
	return COMPILED_FUNCTION_OBJ
} 
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)

}

package object

import (
	"bytes"
	"fmt"
	"github.com/qiuhoude/go-interpreter/ast"
	"hash/fnv"
	"strings"
)

type ObjectType string

type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	STRING_OBJ       ObjectType = "STRING"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	HASH_OBJ         ObjectType = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}
type HashKey struct {
	Type  ObjectType
	Value uint64
}

//type hashKeyObj interface {
//	HashKey() HashKey
//}

type cacheHashKey struct {
	key *HashKey
}

func (c *cacheHashKey) hashKey(keyCreateFn func() *HashKey) HashKey {
	if c.key == nil {
		c.key = keyCreateFn()
	}
	return *c.key
}

// integer
type Integer struct {
	cacheHashKey
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return i.hashKey(func() *HashKey {
		return &HashKey{Type: i.Type(), Value: uint64(i.Value)}
	})
}

// boolean
type Boolean struct {
	cacheHashKey
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) HashKey() HashKey {
	return b.hashKey(func() *HashKey {
		var value uint64
		if b.Value {
			value = 1
		} else {
			value = 0
		}
		return &HashKey{Type: b.Type(), Value: value}
	})

}

// null
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// return Val
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// error
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string
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

// string
type String struct {
	cacheHashKey
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	return s.hashKey(func() *HashKey {
		h := fnv.New64a()
		_, _ = h.Write([]byte(s.Value))
		return &HashKey{Type: s.Type(), Value: h.Sum64()}
	})
}

// 内建函数
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// 数组
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// hash
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	var pairs []string
	for _, p := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			p.Key.Inspect(), p.Key.Inspect()))
	}
	out.WriteString("hash")
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

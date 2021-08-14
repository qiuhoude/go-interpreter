package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
	NULL_OBJ    ObjectType = "NULL_OBJ"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// integer
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// boolean
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// null
type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

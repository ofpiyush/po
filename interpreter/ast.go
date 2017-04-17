package interpreter

import (
	"fmt"
	"strconv"
)

type Type int64

type InfixOperation int64

const (
	_ Type = iota
	INT_TYPE
	STRING_TYPE
	BOOL_TYPE
	HASH_TYPE
)

const (
	_ InfixOperation = iota
	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE
	MOD

	GREATER
	LESSER
	EQUAL
	NOTEQUAL

	AND
	OR
)

var TypeName = map[Type]string{
	INT_TYPE:    "INT",
	STRING_TYPE: "STRING",
	BOOL_TYPE:   "BOOL",
	HASH_TYPE:   "HASH",
}

// Expression is the basic struct of our language.
// Anything on the AST is a node.
type Expression interface {
	Info() string
}

// Object is the most basic an object can go in our language.
// Everything is an object
type Object struct {
	Name  string
	Type  Type
	Value interface{}
}

type HashKey struct {
	Type    Type
	KeyHash uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Lookup struct {
	Name string
}

// Concat is to concat the values of various Expressions
type Concat struct {
	Name   string
	Source []Expression
}

// Assign is used to assign value of Expression on the RHS to the LHS
type Assign struct {
	LHS Expression
	RHS Expression
}

// BlockExpr is a list of Expressions that represent an indent in python / brace in C-like languages
// Eg: if and else blocks, func block
type BlockExpr struct {
	Expressions []Expression
}

// If allows for basic conditional branching
// I am still not sure if we should expose this outside
type If struct {
	Condition Expression
	Success   *BlockExpr
	Fail      *BlockExpr
}

type Infix struct {
	Operator InfixOperation
	Left     Expression
	Right    Expression
}

// Return is used at various places in lambdas
type Return struct {
	Expr Expression
}

// Lambda allows us to write basic functions
type Lambda struct {
	Params map[string]Expression
	Body   *BlockExpr
}

func (*Object) Info() string    { return "I am an Object" }
func (*Lookup) Info() string    { return "I am a Lookup" }
func (*Return) Info() string    { return "I am a Return" }
func (*Concat) Info() string    { return "I am a Concat" }
func (*Assign) Info() string    { return "I am an Assign" }
func (*BlockExpr) Info() string { return "I am a BlockExpr" }
func (*If) Info() string        { return "I am an If" }
func (*Lambda) Info() string    { return "I am a Lambda" }
func (*Infix) Info() string     { return "I am an infix operation" }

func (o *Object) String() string {
	switch o.Type {
	case STRING_TYPE:
		return o.Value.(string)
	case INT_TYPE:
		return strconv.FormatInt(int64(o.Value.(int)), 10)
	case BOOL_TYPE:
		return fmt.Sprintf("%t", o.Value.(bool))
	case HASH_TYPE:
		return fmt.Sprintf("%v", o.Value)
	}
	return ""
}

func (o *Object) HashKey() (key HashKey) {
	key.Type = o.Type
	switch o.Type {
	case STRING_TYPE:
	case INT_TYPE:
		key.KeyHash = uint64(o.Value.(int))
	case BOOL_TYPE:
		key.KeyHash = 0
		if isTruthy(o) {
			key.KeyHash = 1
		}
	case HASH_TYPE:
		// Phat ke phlawar ho jaao piliz
	}
	return
}

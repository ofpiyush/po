package interpreter

import (
	"errors"
	"log"
	"strings"
)

var (
	NilObj   = &Object{}
	TrueObj  = &Object{Type: BOOL, Value: true}
	FalseObj = &Object{Type: BOOL, Value: false}
)

func Eval(s *Scope, expr Expression) (*Object, error) {
	switch expr.(type) {
	case *Object:
		return expr.(*Object), nil
	case *Concat:
		return evalConcat(s, expr.(*Concat))
	case *Assign:
		return evalAssign(s, expr.(*Assign))
	case *If:
		return evalIfStmt(s, expr.(*If))
	case *BlockExpr:
		return evalBlockExpr(s, expr.(*BlockExpr))
	case *Lookup:
		return evalLookup(s, expr.(*Lookup))
	case *Return:
		return Eval(s, expr.(*Return).Expr)
	case *Lambda:
		return evalLambda(s, expr.(*Lambda))
	case *Infix:
		return evalInfixExpr(s, expr.(*Infix))
	default:
		log.Fatal("Unknown expression ", expr.Info())
	}
	return nil, nil
}

func evalConcat(s *Scope, c *Concat) (*Object, error) {
	str := []string{}

	for _, v := range c.Source {
		tmp, err := NilErrCheck(Eval(s, v))
		if err != nil {
			return nil, err
		}
		str = append(str, tmp.String())
	}

	return &Object{Type: STRING, Value: strings.Join(str, "")}, nil
}

func evalAssign(s *Scope, a *Assign) (*Object, error) {

	// Eval RHS
	rhs, err := Eval(s, a.RHS)
	if err != nil {
		return nil, err
	}

	// Find the object and check if it exists in scope
	// Burst if it is not an object or child field of a struct
	switch a.LHS.(type) {
	case *Object:
		obj, err := s.Update(a.LHS.(*Object).Name, rhs)
		if err != nil {
			return nil, err
		}
		return obj, nil
	default:
		return nil, errors.New("Can not assign to LHS")
	}
}

func evalIfStmt(s *Scope, i *If) (*Object, error) {
	obj, err := Eval(s, i.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(obj) {
		return Eval(s, i.Success)
	}
	return Eval(s, i.Fail)
}

func evalBlockExpr(s *Scope, b *BlockExpr) (*Object, error) {
	for _, v := range b.Expressions {
		res, err := Eval(s, v)
		if err != nil {
			return nil, err
		}
		switch v.(type) {
		case *Return:
			return res, nil
		}
	}
	return NilObj, nil
}

func evalLookup(s *Scope, l *Lookup) (*Object, error) {
	obj := s.Lookup(l.Name)
	if obj.Type == NilObj.Type {
		return nil, errors.New(l.Name + " not found")
	}
	return obj, nil
}

func evalLambda(s *Scope, l *Lambda) (*Object, error) {
	child := NewScope("child", s)
	for k, v := range l.Params {
		data, err := Eval(s, v)
		if err != nil {
			return nil, err
		}
		child.Update(k, data)
	}
	return evalBlockExpr(child, l.Body)
}

func evalInfixExpr(s *Scope, i *Infix) (*Object, error) {
	var IntInput = map[InfixOperation]bool{
		ADD:      true,
		SUBTRACT: true,
		MULTIPLY: true,
		DIVIDE:   true,
		MOD:      true,
		GREATER:  true,
		LESSER:   true,
	}

	l, err := NilErrCheck(Eval(s, i.Left))
	if err != nil {
		return nil, err
	}
	r, err := NilErrCheck(Eval(s, i.Right))
	if err != nil {
		return nil, err
	}
	var returnObj = &Object{}

	if IntInput[i.Operator] {
		if l.Type != INT {
			return nil, errors.New("LHS is not an integer")
		}
		if r.Type != INT {
			return nil, errors.New("RHS is not an integer")
		}
	}

	switch i.Operator {
	case ADD:
		returnObj.Type = INT
		returnObj.Value = l.Value.(int) + r.Value.(int)
	case SUBTRACT:
		returnObj.Type = INT
		returnObj.Value = l.Value.(int) - r.Value.(int)
	case MULTIPLY:
		returnObj.Type = INT
		returnObj.Value = l.Value.(int) * r.Value.(int)
	case DIVIDE:
		returnObj.Type = INT
		returnObj.Value = l.Value.(int) / r.Value.(int)
	case MOD:
		returnObj.Type = INT
		returnObj.Value = l.Value.(int) % r.Value.(int)

	case GREATER:
		returnObj.Type = BOOL
		returnObj.Value = l.Value.(int) > r.Value.(int)
	case LESSER:
		returnObj.Type = BOOL
		returnObj.Value = l.Value.(int) < r.Value.(int)
	case EQUAL:
		returnObj.Type = BOOL
		returnObj.Value = false
		if l.Type == r.Type && l.Value == r.Value {
			returnObj.Value = true
		}
	case NOTEQUAL:
		returnObj.Type = BOOL
		returnObj.Value = true
		if l.Type == r.Type && l.Value == r.Value {
			returnObj.Value = false
		}
	case AND:
		returnObj.Type = BOOL
		returnObj.Value = isTruthy(l) && isTruthy(r)
	case OR:
		returnObj.Type = BOOL
		returnObj.Value = isTruthy(l) || isTruthy(r)
	}
	return returnObj, nil
}

func NilErrCheck(obj *Object, err error) (*Object, error) {
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, errors.New("Got nil from eval")
	}
	return obj, err
}

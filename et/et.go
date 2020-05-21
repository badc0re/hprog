package et

import (
	"fmt"
	"github.com/badc0re/hprog/token"
	"github.com/pkg/errors"
)

type InternalType int

const (
	ILLEGALInternalType InternalType = iota

	//
	IntInternalType
	f64InternalType
	StringInternalType

	// FunctionInternalType
)

type Expr interface {
	Accept(Expr) (Object, error)
}

type Object struct {
	Value        interface{}
	InternalType string
}

type Binary struct {
	Operator token.Token
	Left     Expr
	Right    Expr
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

type Literal struct {
	TypedObject Object
}

type Grouping struct {
	Expression Expr
}

func (bexpr Binary) Accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Binary)
	return visitor.visitBinaryExpr(bexpr)
}

func (uexpr Unary) Accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Unary)
	return visitor.visitUnaryExpr(uexpr)
}

func (lexpr Literal) Accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Literal)
	return visitor.visitLiteralExpr(lexpr)
}

func (gexpr Grouping) Accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Grouping)
	return visitor.visitGroupingExpr(gexpr)
}

func evaluate(expr Expr) (Object, error) {
	return expr.Accept(expr)
}

func (thisExpr Binary) visitBinaryExpr(inputExpr Binary) (Object, error) {
	// I know this is ugly
	left, _ := evaluate(inputExpr.Left)
	right, _ := evaluate(inputExpr.Right)

	if left.InternalType != right.InternalType {
		return Object{}, errors.New(fmt.Sprintf("Missmatch type %s and %s", left.InternalType, right.InternalType))
	}

	if inputExpr.Operator.Type == token.MINUS {
		if left.InternalType == "int" && right.InternalType == "int" {
			left_value := left.Value.(int)
			right_value := right.Value.(int)
			return Object{left_value - right_value, "int"}, nil
		}
		if left.InternalType == "f64" && right.InternalType == "f64" {
			left_value := left.Value.(float64)
			right_value := right.Value.(float64)
			return Object{left_value - right_value, "f64"}, nil
		}
	} else if inputExpr.Operator.Type == token.SLASH {
		if left.InternalType == "int" && right.InternalType == "int" {
			left_value := left.Value.(int)
			right_value := right.Value.(int)
			return Object{left_value / right_value, "int"}, nil
		}
		if left.InternalType == "f64" && right.InternalType == "f64" {
			left_value := left.Value.(float64)
			right_value := right.Value.(float64)
			return Object{left_value / right_value, "f64"}, nil
		}
	} else if inputExpr.Operator.Type == token.STAR {
		if left.InternalType == "int" && right.InternalType == "int" {
			left_value := left.Value.(int)
			right_value := right.Value.(int)
			return Object{left_value * right_value, "int"}, nil
		}
		if left.InternalType == "f64" && right.InternalType == "f64" {
			left_value := left.Value.(float64)
			right_value := right.Value.(float64)
			return Object{left_value * right_value, "f64"}, nil
		}
	} else if inputExpr.Operator.Type == token.PLUS {
		if left.InternalType == "int" && right.InternalType == "int" {
			left_value := left.Value.(int)
			right_value := right.Value.(int)
			return Object{left_value + right_value, "int"}, nil
		}
		if left.InternalType == "f64" && right.InternalType == "f64" {
			left_value := left.Value.(float64)
			right_value := right.Value.(float64)
			return Object{left_value + right_value, "f64"}, nil
		}
	}
	return Object{}, nil
}

func (thisExpr Unary) visitUnaryExpr(inputExpr Unary) (Object, error) {
	right, _ := evaluate(inputExpr.Right)
	if inputExpr.Operator.Type == token.MINUS {
		if right.InternalType == "int" {
			return Object{-right.Value.(int), "int"}, nil
		}
		if right.InternalType == "f64" {
			return Object{-right.Value.(float64), "f64"}, nil
		}
	}
	return Object{}, nil
}

func (thisExpr Literal) visitLiteralExpr(inputExpr Literal) (Object, error) {
	return inputExpr.TypedObject, nil
}

func (thisExpr Grouping) visitGroupingExpr(inputExpr Expr) (Object, error) {
	return evaluate(thisExpr.Expression)
}

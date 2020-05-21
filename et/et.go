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
	accept(Expr) (Object, error)
}

type Object struct {
	value        interface{}
	internalType string
}

type Binary struct {
	operator token.Token
	left     Expr
	right    Expr
}

type Unary struct {
	operator token.Token
	right    Expr
}

type Literal struct {
	Object Object
}

type Grouping struct {
	expression Expr
}

func (bexpr Binary) accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Binary)
	return visitor.visitBinaryExpr(bexpr)
}

func (uexpr Unary) accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Unary)
	return visitor.visitUnaryExpr(uexpr)
}

func (lexpr Literal) accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Literal)
	return visitor.visitLiteralExpr(lexpr)
}

func (gexpr Grouping) accept(expr Expr) (Object, error) {
	// tests type
	visitor, _ := expr.(Grouping)
	return visitor.visitGroupingExpr(gexpr)
}

func evaluate(expr Expr) (Object, error) {
	return expr.accept(expr)
}

func (thisExpr Binary) visitBinaryExpr(inputExpr Binary) (Object, error) {
	// I know this is ugly
	left, _ := evaluate(inputExpr.left)
	right, _ := evaluate(inputExpr.right)

	if left.internalType != right.internalType {
		return Object{}, errors.New(fmt.Sprintf("Missmatch type %s and %s", left.internalType, right.internalType))
	}

	if inputExpr.operator.Type == token.MINUS {
		if left.internalType == "int" && right.internalType == "int" {
			left_value := left.value.(int)
			right_value := right.value.(int)
			return Object{left_value - right_value, "int"}, nil
		}
		if left.internalType == "f64" && right.internalType == "f64" {
			left_value := left.value.(float64)
			right_value := right.value.(float64)
			return Object{left_value - right_value, "f64"}, nil
		}
	} else if inputExpr.operator.Type == token.SLASH {
		if left.internalType == "int" && right.internalType == "int" {
			left_value := left.value.(int)
			right_value := right.value.(int)
			return Object{left_value / right_value, "int"}, nil
		}
		if left.internalType == "f64" && right.internalType == "f64" {
			left_value := left.value.(float64)
			right_value := right.value.(float64)
			return Object{left_value / right_value, "f64"}, nil
		}
	} else if inputExpr.operator.Type == token.STAR {
		if left.internalType == "int" && right.internalType == "int" {
			left_value := left.value.(int)
			right_value := right.value.(int)
			return Object{left_value * right_value, "int"}, nil
		}
		if left.internalType == "f64" && right.internalType == "f64" {
			left_value := left.value.(float64)
			right_value := right.value.(float64)
			return Object{left_value * right_value, "f64"}, nil
		}
	} else if inputExpr.operator.Type == token.PLUS {
		if left.internalType == "int" && right.internalType == "int" {
			left_value := left.value.(int)
			right_value := right.value.(int)
			return Object{left_value + right_value, "int"}, nil
		}
		if left.internalType == "f64" && right.internalType == "f64" {
			left_value := left.value.(float64)
			right_value := right.value.(float64)
			return Object{left_value + right_value, "f64"}, nil
		}
	}
	return Object{}, nil
}

func (thisExpr Unary) visitUnaryExpr(inputExpr Unary) (Object, error) {
	right, _ := evaluate(inputExpr.right)
	if inputExpr.operator.Type == token.MINUS {
		if right.internalType == "int" {
			return Object{-right.value.(int), "int"}, nil
		}
		if right.internalType == "f64" {
			return Object{-right.value.(float64), "f64"}, nil
		}
	}
	return Object{}, nil
}

func (thisExpr Literal) visitLiteralExpr(inputExpr Literal) (Object, error) {
	return inputExpr.object, nil
}

func (thisExpr Grouping) visitGroupingExpr(inputExpr Expr) (Object, error) {
	return evaluate(thisExpr.expression)
}

package main

import (
	"fmt"
	"strconv"
)

type Expr interface {
	accept(Expr) Object
}

type Object interface{}

type Binary struct {
	operator Token
	left     Expr
	right    Expr
}

type Unary struct {
	operator Token
	right    Expr
}

type Literal struct {
	object Object
}

type Grouping struct {
	expression Expr
}

func (bexpr Binary) accept(expr Expr) Object {
	// tests type
	visitor, _ := expr.(Binary)
	return visitor.visitBinaryExpr(bexpr)
}

func (uexpr Unary) accept(expr Expr) Object {
	// tests type
	visitor, _ := expr.(Unary)
	return visitor.visitUnaryExpr(uexpr)
}

func (lexpr Literal) accept(expr Expr) Object {
	// tests type
	visitor, _ := expr.(Literal)
	return visitor.visitLiteralExpr(lexpr)
}

func (gexpr Grouping) accept(expr Expr) Object {
	// tests type
	visitor, _ := expr.(Grouping)
	return visitor.visitGroupingExpr(gexpr)
}

func evaluate(expr Expr) Object {
	return expr.accept(expr)
}

func (thisExpr Binary) visitBinaryExpr(inputExpr Binary) Object {
	fmt.Println("This binary:", inputExpr)
	left := evaluate(inputExpr.left)
	right := evaluate(inputExpr.right)
	fmt.Println("left:", inputExpr.left, "right:", inputExpr.right)

	if inputExpr.operator.tokenType == MINUS {
		right_value := right.(float64)
		left_value := left.(float64)
		return left_value - right_value
	} else if inputExpr.operator.tokenType == SLASH {
		return nil
	} else if inputExpr.operator.tokenType == STAR {
		return nil
	}
	return nil
}

func (thisExpr Unary) visitUnaryExpr(inputExpr Unary) Object {
	right := evaluate(inputExpr.right)

	if inputExpr.operator.tokenType == MINUS {
		objectToString, _ := right.(string)
		stringToValue, _ := strconv.ParseFloat(objectToString, 32)
		return -stringToValue
	}

	return nil
}

func (thisExpr Literal) visitLiteralExpr(inputExpr Literal) Object {
	fmt.Println("inputExpr:", inputExpr)
	return inputExpr.object
}

func (thisExpr Grouping) visitGroupingExpr(inputExpr Expr) Object {
	// fmt.Println("This grouping:", inputExpr)
	// fmt.Println("This grouping expr:", thisExpr.expression)
	return evaluate(thisExpr.expression)
}

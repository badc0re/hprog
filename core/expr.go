package main

import "fmt"

type Expr interface {
	accept(Expr) Expr
}

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
	value interface{}
}

type Grouping struct {
	expression Expr
}

func (bexpr Binary) accept(expr Expr) Expr {
	// tests type
	visitor, _ := expr.(Binary)
	return visitor.visitBinaryExpr(bexpr)
}

func (uexpr Unary) accept(expr Expr) Expr {
	// tests type
	visitor, _ := expr.(Unary)
	return visitor.visitUnaryExpr(uexpr)
}

func (lexpr Literal) accept(expr Expr) Expr {
	// tests type
	visitor, _ := expr.(Literal)
	return visitor.visitLiteralExpr(lexpr)
}

func (gexpr Grouping) accept(expr Expr) Expr {
	// tests type
	visitor, _ := expr.(Grouping)
	return visitor.visitGroupingExpr(gexpr)
}

func (thisExpr Binary) visitBinaryExpr(inputExpr Expr) Expr {
	fmt.Println("This binary:", inputExpr)
	return thisExpr
}

func (thisExpr Unary) visitUnaryExpr(inputExpr Expr) Expr {
	fmt.Println("This unary:", inputExpr)
	return thisExpr
}

func (thisExpr Literal) visitLiteralExpr(inputExpr Expr) Expr {
	fmt.Println("This literal:", inputExpr)
	return thisExpr
}

func (thisExpr Grouping) visitGroupingExpr(inputExpr Expr) Expr {
	// fmt.Println("This grouping:", inputExpr)
	// fmt.Println("This grouping expr:", thisExpr.expression)
	return thisExpr
	//return inputExpr.accept(thisExpr.expression)
}

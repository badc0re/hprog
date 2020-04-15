package main

import "fmt"

type Expr interface {
	accept(Expr) string
}

type Token struct {
	value string
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

func (bexpr Binary) accept(expr Expr) string {
	// tests type
	binary, _ := expr.(Binary)
	return bexpr.visitBinaryExpr(binary)
}

func (uexpr Unary) accept(expr Expr) string {
	// tests type
	unary, _ := expr.(Unary)
	return uexpr.visitUnaryExpr(unary)
}

func (lexpr Literal) accept(expr Expr) string {
	// tests type
	literal, _ := expr.(Literal)
	return lexpr.visitLiteralExpr(literal)
}

func (gexpr Grouping) accept(expr Expr) string {
	// tests type
	grouping, _ := expr.(Grouping)
	return gexpr.visitGroupingExpr(grouping)
}
func (thisExpr Binary) visitBinaryExpr(inputExpr Expr) string {
	binary, _ := inputExpr.(Binary)
	return parents(binary.operator.value, binary.left, binary.right)
}

func (thisExpr Unary) visitUnaryExpr(inputExpr Expr) string {
	unary, _ := inputExpr.(Unary)
	return parents(unary.operator.value, unary.right)
}

func (thisExpr Literal) visitLiteralExpr(inputExpr Expr) string {
	literal, _ := inputExpr.(Literal)
	return literal.value.(string)
}

func (thisExpr Grouping) visitGroupingExpr(inputExpr Expr) string {
	grouping, _ := inputExpr.(Grouping)
	return parents("group", grouping.expression)
}

func parents(name string, expr ...Expr) string {
	var ast string
	ast += "(" + name
	for _, e := range expr {
		ast += " "
		ast += e.accept(e)
	}
	ast += ")"
	return ast
}

func prints(expr Expr) {
	fmt.Println(expr.accept(expr))
}

func main() {
	b := Binary{
		operator: Token{
			value: "*",
		},
		left: Unary{
			operator: Token{
				value: "-",
			},
			right: Literal{
				value: "123",
			},
		},
		right: Grouping{
			expression: Literal{
				value: "45.67",
			},
		},
	}
	prints(b)
}

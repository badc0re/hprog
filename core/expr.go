package main

//import "fmt"

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
	binary, _ := expr.(Binary)
	return bexpr.visitBinaryExpr(binary)
}

func (uexpr Unary) accept(expr Expr) Expr {
	// tests type
	unary, _ := expr.(Unary)
	return uexpr.visitUnaryExpr(unary)
}

func (lexpr Literal) accept(expr Expr) Expr {
	// tests type
	literal, _ := expr.(Literal)
	return lexpr.visitLiteralExpr(literal)
}

func (gexpr Grouping) accept(expr Expr) Expr {
	// tests type
	grouping, _ := expr.(Grouping)
	return gexpr.visitGroupingExpr(grouping)
}

func (thisExpr Binary) visitBinaryExpr(inputExpr Expr) Expr {
	return inputExpr
	//return parents()

}

func (thisExpr Unary) visitUnaryExpr(inputExpr Expr) Expr {
	return inputExpr
}

func (thisExpr Literal) visitLiteralExpr(inputExpr Expr) Expr {
	return inputExpr
}

func (thisExpr Grouping) visitGroupingExpr(inputExpr Expr) Expr {
	return inputExpr
}

func parents(name string, expr ...Expr) {

}

func bla() {
	b := Binary{
		operator: Token{
			tokenType: STAR,
			pos:       0,
			end:       0,
			line:      0,
			value:     "*",
		},
		left: Unary{
			operator: Token{
				tokenType: MINUS,
				pos:       0,
				end:       0,
				line:      0,
				value:     "-",
			},
			right: Literal{
				value: 123,
			},
		},
		right: Grouping{
			expression: Literal{
				value: 45.67,
			},
		},
	}
}

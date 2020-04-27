package main

import (
	"fmt"
)

func parents(name string, expr ...Expr) string {
	var ast string
	ast += "(" + name
	for _, e := range expr {
		ast += " "
		current_exp := e.accept(e)
		lit, _ := current_exp.(Literal)
		ast += lit.value.(string)
	}
	ast += ")"
	return ast
}

func prints(expr Expr) {
	fmt.Println(expr.accept(expr))
}

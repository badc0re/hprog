package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func astPrinter(expr Expr) {
	// fmt.Printf("%+v\n", expr)
	fmt.Println(spew.Sdump(expr))
	fmt.Println("Result:")
	fmt.Println(spew.Sdump(expr.accept(expr)))
}

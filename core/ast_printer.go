package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func prints(expr Expr) {
	// fmt.Printf("%+v\n", expr)
	fmt.Println(spew.Sdump(expr))
	expr.accept(expr)
}

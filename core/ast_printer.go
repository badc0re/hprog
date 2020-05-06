package main

import (
	"fmt"
)

func prints(expr Expr) {
	fmt.Printf("%#v\n", expr)
	expr.accept(expr)
}

package main

import "fmt"

type Expr interface {
	readVariable() Pos
}

type Binary struct {
	operator Token
	left     Expr
	right    Expr
}

func (b Binary) readVariable() Pos {
	return b.operator.pos
}

func read(b Expr) {
	fmt.Println("%d", b.readVariable())
}

func accept() {
	token := Token{
		tokenType: DOT,
		pos:       30,
		end:       1,
		line:      1,
		value:     ".",
	}
	b := Binary{operator: token}
	b2 := Binary{operator: token, left: b}
	read(b2.left)
}

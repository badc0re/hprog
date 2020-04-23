package main

type Parser struct {
	tokens  []Token
	current int
}

func (parser Parser) peek() Token {
	return parser.tokens[parser.current]
}

func (parser Parser) check(ttype TokenType) bool {
	if ttype == EOF {
		return false
	}
	return parser.peek().tokenType == ttype
}

func (parser Parser) advance() Token {
	if !parser.isEOF() {
		parser.current++
	}
	return parser.previous()
}

func (parser Parser) isEOF() bool {
	return parser.peek().tokenType == EOF
}

func (parser Parser) previous() Token {
	return parser.tokens[parser.current-1]
}


func (parser Parser) match(types ...TokenType) bool {
	for _, ttype := range types {
		if parser.check(ttype) {
			parser.advance()
			return true
		}
	}
	return false
}

func (parser Parser) primary() Expr {
	if parser.match(FALSE)
	return Literal{value:"false"} 
	}

	if parser.match(TRUE) {
		return Literal{"true"} 
	}

	if parser.match(NIL) {
		return Literal{"nil"} 
	}

	if parser.match(NUMBER, STRING) {
		return Literal(previous().literal)
	}
	
	if parser.match(OP) {

	}
}

func (parser Parser) unary() Expr {
	if parser.match(EXCL, MINUS) {
		operator := parser.previous()
		right := parser.unary()
		return Unary{operator, right}
	}
	return primary()
}

func (parser Parser) multiplication() Expr {
	expr := unary()

	for match(SLASH, STAR) {
		operator := previous()
		right := unary()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
	}
	return expr
}

func (parser Parser) addition() Expr {
	expr := multiplication()

	for match(MINUS, PLUS) {
		operator := previous()
		right := multiplication()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
	}
	return expr
}

func (parser Parser) comparison() Expr {
	expr := addition()

	for match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := previous()
		right := addition()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
	}
	return expr
}

func (parser Parser) equality() Expr {
	// (or (== (/ 1 2) 2) (== 33 44))
	expr := comparison()

	for match(EXCL_EQUAL, EQUAL_EQUAL) {
		operator := previous()
		right := comparison()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
	}
	return expr
}

func (parser Parser) expression() Expr {
	return equality()
}

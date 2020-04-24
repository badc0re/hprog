package main

import (
	"errors"
	"fmt"
)

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
	fmt.Println(parser.current)
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
			fmt.Println(types)
			parser.advance()
			return true
		}
	}
	return false
}

func (parser Parser) consume(tokenType TokenType, msg string) (Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}
	return parser.peek(), errors.New(msg)
}

func (parser Parser) primary() (Expr, error) {
	if parser.match(BOOL_FALSE) {
		return Literal{value: "false"}, nil
	}
	if parser.match(BOOL_TRUE) {
		return Literal{value: "true"}, nil
	}

	if parser.match(NIL) {
		return Literal{value: "nil"}, nil
	}

	if parser.match(NUMBER, STRING) {
		return Literal{parser.previous().value}, nil
	}

	if parser.match(OP) {
		expr, err := parser.expression()
		parser.consume(CP, "Expected ')'")
		return expr, err
	}
	return nil, errors.New("Expected expression.")
}

func (parser Parser) unary() (Expr, error) {
	if parser.match(EXCL, MINUS) {
		operator := parser.previous()
		right, err := parser.unary()
		return Unary{operator, right}, err
	}
	return parser.primary()
}

func (parser Parser) multiplication() (Expr, error) {
	expr, err := parser.unary()

	for parser.match(SLASH, STAR) {
		operator := parser.previous()
		right, _err := parser.unary()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
		err = _err
	}
	return expr, err
}

func (parser Parser) addition() (Expr, error) {
	expr, err := parser.multiplication()

	for parser.match(MINUS, PLUS) {
		operator := parser.previous()
		right, _err := parser.multiplication()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
		err = _err
	}
	return expr, err
}

func (parser Parser) comparison() (Expr, error) {
	expr, err := parser.addition()

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right, _err := parser.addition()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
		err = _err
	}
	return expr, err
}

func (parser Parser) equality() (Expr, error) {
	// (or (== (/ 1 2) 2) (== 33 44))
	expr, err := parser.comparison()

	for parser.match(EXCL_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right, _err := parser.comparison()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
		// NOTE: can be propagated back without changing the current
		err = _err
	}
	return expr, err
}

func (parser Parser) expression() (Expr, error) {
	fmt.Println("ass")
	bla, err := parser.equality()
	fmt.Println("error", err)
	return bla, err
}

func bla() {
	var tokens []Token

	lex := startGrinding("(<= 10 10)")
	for {
		token := lex.nextToken()
		tokens = append(tokens, token)
		//fmt.Println(reverseKeys[token.tokenType])
		if token.tokenType == EOF || token.tokenType == ERR {
			break
		}
		//fmt.Println(token.value)
	}
	parser := Parser{tokens: tokens, current: 0}
	parser.expression()
}

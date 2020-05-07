package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type Parser struct {
	tokens  []Token
	current int
}

func (parser *Parser) peek() Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) check(ttype TokenType) bool {
	if ttype == EOF {
		return false
	}

	return parser.peek().tokenType == ttype
}

func (parser *Parser) advance() Token {
	if !parser.isEOF() {
		parser.current++
	}
	return parser.previous()
}

func (parser *Parser) isEOF() bool {
	return parser.peek().tokenType == EOF
}

func (parser *Parser) previous() Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) match(types ...TokenType) bool {
	for _, ttype := range types {
		if parser.check(ttype) {
			parser.advance()
			return true
		}
	}
	return false
}

func (parser *Parser) consume(tokenType TokenType, msg string) (Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}
	return parser.peek(), errors.New(msg)
}

func (parser *Parser) primary() (Expr, error) {
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

	return nil, errors.New("Cannot parse expression.")
}

func (parser *Parser) unary() (Expr, error) {
	if parser.match(EXCL, MINUS) {
		operator := parser.previous()
		expr, err := parser.unary()
		if err != nil {
			return nil, errors.WithMessage(err, "Cannot parser construct unary expression.")
		}
		return Unary{operator, expr}, err
	}
	return parser.primary()
}

func (parser *Parser) statement() (Expr, error) {
	expr, err := parser.unary()
	if parser.match(OP) {
		// first the operator
		operator := parser.advance()

		left_expr, err := parser.expression()
		if err != nil {
			return nil, errors.WithMessage(err, "Cannot parse left expresssion.")
		}
		right_expr, err := parser.expression()
		if err != nil {
			return nil, errors.WithMessage(err, "Cannot parse right expression.")
		}
		_, err = parser.consume(CP, "Expected ')' for the expression.")
		if err != nil {
			return nil, errors.WithMessage(err, "Cannot parse right expression.")
		}
		return Grouping{
			Binary{
				operator: operator,
				left:     left_expr,
				right:    right_expr,
			},
		}, err
	}

	// to be replaced
	return expr, err
}

func (parser *Parser) equality() (Expr, error) {
	return parser.statement()
}

/*

func (parser *Parser) multiplication() (Expr, error) {
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

func (parser *Parser) addition() (Expr, error) {
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

func (parser *Parser) comparison() (Expr, error) {
	expr, err := parser.addition()

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right, _err := parser.addition()
		expr = Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}
		err = _eriiiuur
	}
	return expr, err
}

func (parser *Parser) equality() (Expr, error) {
	// (or (== (/ 1 2) 2) (== 33 44))
	expr, err := parser.comparison()

	for parser.match(EXCL_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right, _err := parser.comparison()
		fmt.Println("op", operator.value)
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
*/

func (parser *Parser) expression() (Expr, error) {
	return parser.equality()
}

func test() {
	var tokens []Token

	lex := startGrinding("(== (max -3 -4) 4)(AA)")
	for {
		token := lex.nextToken()
		tokens = append(tokens, token)
		if token.tokenType == EOF || token.tokenType == ERR {
			break
		}
	}
	//fmt.Printf("%#v\n\n", tokens)

	parser := Parser{tokens: tokens, current: 0}
	expr, err := parser.expression()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	prints(expr)
}

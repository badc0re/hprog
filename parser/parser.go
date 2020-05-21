package parser

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"os"
	"strconv"
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
		return Literal{object: Object{value: "false"}}, nil
	}
	if parser.match(BOOL_TRUE) {
		return Literal{object: Object{value: "true"}}, nil
	}

	if parser.match(NIL) {
		return Literal{object: Object{value: "nil"}}, nil
	}

	if parser.match(NUMBER) {
		msg := "Invalid expression for number"
		_err := errors.New(msg)
		value := parser.previous().value
		if val, err := strconv.Atoi(value); err == nil {
			return Literal{object: Object{value: val, internalType: "int"}}, nil
		}
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			return Literal{object: Object{value: val, internalType: "f64"}}, nil
		}
		return nil, _err
	}

	if parser.match(STRING) {
		return Literal{object: Object{value: parser.previous().value}}, nil
	}

	if parser.match(OP) {
		expr, err := parser.expression()
		_, err = parser.consume(CP, "Expected ')' for the expression.")
		return Grouping{expr}, err
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

func (parser *Parser) multiplication() (Expr, error) {
	expr, err := parser.unary()

	for parser.match(SLASH, STAR) {
		operator := parser.previous()
		right, err := parser.unary()
		return Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) addition() (Expr, error) {
	expr, err := parser.multiplication()

	for parser.match(MINUS, PLUS) {
		operator := parser.previous()
		right, err := parser.multiplication()
		return Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) comparison() (Expr, error) {
	expr, err := parser.addition()

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right, err := parser.addition()
		return Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) equality() (Expr, error) {
	// (or (== (/ 1 2) 2) (== 33 44))
	expr, err := parser.comparison()

	for parser.match(EXCL_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right, err := parser.comparison()
		return Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) expression() (Expr, error) {
	return parser.equality()
}

func parserExecutor() {
	var tokens []Token

	lex := startGrinding("-2 - -4")
	for {
		token := lex.nextToken()
		tokens = append(tokens, token)
		if token.tokenType == EOF || token.tokenType == ERR {
			break
		}
	}
	//fmt.Println(spew.Sdump(tokens))
	//fmt.Printf("%#v\n\n", tokens)

	parser := Parser{tokens: tokens, current: 0}
	expr, err := parser.expression()
	fmt.Println(err)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	astPrinter(expr)
}

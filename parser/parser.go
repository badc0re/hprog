package parser

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/badc0re/hprog/et"
	"github.com/badc0re/hprog/token"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) check(ttype token.TokenType) bool {
	if ttype == token.EOF {
		return false
	}

	return parser.peek().Type == ttype
}

func (parser *Parser) advance() token.Token {
	if !parser.isEOF() {
		parser.current++
	}
	return parser.previous()
}

func (parser *Parser) isEOF() bool {
	return parser.peek().Type == token.EOF
}

func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) match(types ...token.TokenType) bool {
	for _, ttype := range types {
		if parser.check(ttype) {
			parser.advance()
			return true
		}
	}
	return false
}

func (parser *Parser) consume(Type token.TokenType, msg string) (token.Token, error) {
	if parser.check(Type) {
		return parser.advance(), nil
	}
	return parser.peek(), errors.New(msg)
}

func (parser *Parser) primary() (et.Expr, error) {
	if parser.match(token.BOOL_FALSE) {
		return et.Literal{object: Object{value: "false"}}, nil
	}
	if parser.match(token.BOOL_TRUE) {
		return et.Literal{object: Object{value: "true"}}, nil
	}

	if parser.match(token.NIL) {
		return et.Literal{object: Object{value: "nil"}}, nil
	}

	if parser.match(token.NUMBER) {
		msg := "Invalid expression for number"
		_err := errors.New(msg)
		value := parser.previous().value
		if val, err := strconv.Atoi(value); err == nil {
			return et.Literal{object: Object{value: val, internalType: "int"}}, nil
		}
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			return et.Literal{object: Object{value: val, internalType: "f64"}}, nil
		}
		return nil, _err
	}

	if parser.match(token.STRING) {
		return et.Literal{object: Object{value: parser.previous().value}}, nil
	}

	if parser.match(token.OP) {
		expr, err := parser.expression()
		_, err = parser.consume(token.CP, "Expected ')' for the expression.")
		return et.Grouping{expr}, err
	}

	return nil, errors.New("Cannot parse expression.")
}

func (parser *Parser) unary() (et.Expr, error) {
	if parser.match(token.EXCL, token.MINUS) {
		operator := parser.previous()
		expr, err := parser.unary()
		if err != nil {
			return nil, errors.WithMessage(err, "Cannot parser construct unary expression.")
		}
		return et.Unary{operator, expr}, err
	}
	return parser.primary()
}

func (parser *Parser) multiplication() (et.Expr, error) {
	expr, err := parser.unary()

	for parser.match(token.SLASH, token.STAR) {
		operator := parser.previous()
		right, err := parser.unary()
		return et.Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) addition() (et.Expr, error) {
	expr, err := parser.multiplication()

	for parser.match(token.MINUS, token.PLUS) {
		operator := parser.previous()
		right, err := parser.multiplication()
		return et.Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) comparison() (et.Expr, error) {
	expr, err := parser.addition()

	for parser.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := parser.previous()
		right, err := parser.addition()
		return et.Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) equality() (et.Expr, error) {
	// (or (== (/ 1 2) 2) (== 33 44))
	expr, err := parser.comparison()

	for parser.match(token.EXCL_EQUAL, token.EQUAL_EQUAL) {
		operator := parser.previous()
		right, err := parser.comparison()
		return et.Binary{
			operator: operator,
			left:     expr,
			right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) expression() (et.Expr, error) {
	return parser.equality()
}

func parserExecutor() {
	var tokens []token.Token

	lex := et.startGrinding("-2 - -4")
	for {
		token := lex.nextToken()
		tokens = append(tokens, token)
		if token.Type == token.EOF || token.Type == token.ERR {
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
	et.astPrinter(expr)
}

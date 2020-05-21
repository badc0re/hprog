package parser

import (
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
	//"os"
	"github.com/badc0re/hprog/et"
	"github.com/badc0re/hprog/token"
	"github.com/pkg/errors"
	"strconv"
)

type Parser struct {
	Tokens  []token.Token
	Current int
}

func (parser *Parser) peek() token.Token {
	return parser.Tokens[parser.Current]
}

func (parser *Parser) check(ttype token.TokenType) bool {
	if ttype == token.EOF {
		return false
	}

	return parser.peek().Type == ttype
}

func (parser *Parser) advance() token.Token {
	if !parser.isEOF() {
		parser.Current++
	}
	return parser.previous()
}

func (parser *Parser) isEOF() bool {
	return parser.peek().Type == token.EOF
}

func (parser *Parser) previous() token.Token {
	return parser.Tokens[parser.Current-1]
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
		return et.Literal{TypedObject: et.Object{Value: "false"}}, nil
	}
	if parser.match(token.BOOL_TRUE) {
		return et.Literal{TypedObject: et.Object{Value: "true"}}, nil
	}

	if parser.match(token.NIL) {
		return et.Literal{TypedObject: et.Object{Value: "nil"}}, nil
	}

	if parser.match(token.NUMBER) {
		msg := "Invalid expression for number"
		_err := errors.New(msg)
		value := parser.previous().Value
		if val, err := strconv.Atoi(value); err == nil {
			return et.Literal{TypedObject: et.Object{Value: val, InternalType: "int"}}, nil
		}
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			return et.Literal{TypedObject: et.Object{Value: val, InternalType: "f64"}}, nil
		}
		return nil, _err
	}

	if parser.match(token.STRING) {
		return et.Literal{TypedObject: et.Object{Value: parser.previous().Value}}, nil
	}

	if parser.match(token.OP) {
		expr, err := parser.Expression()
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
			Operator: operator,
			Left:     expr,
			Right:    right,
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
			Operator: operator,
			Left:     expr,
			Right:    right,
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
			Operator: operator,
			Left:     expr,
			Right:    right,
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
			Operator: operator,
			Left:     expr,
			Right:    right,
		}, err
	}
	return expr, err
}

func (parser *Parser) Expression() (et.Expr, error) {
	return parser.equality()
}

/*
func parserExecutor() {
	var tokens []token.Token

	lex := lexer.consume("-2 - -4")
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
	expr, err := parser.Expression()
	fmt.Println(err)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	//et.astPrinter(expr)
}
*/

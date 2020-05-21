package token

import "fmt"

type TokenType int
type TokenPos int
type TokenLine int

var Eof = rune(0)

const (
	ILLEGAL TokenType = iota

	// single char tokens
	OP
	CP
	LB
	RB
	PLUS
	SLASH
	STAR
	COMMA
	DOT
	MINUS
	SEMICOLON
	QUOTE

	COLON

	GREATER
	GREATER_EQUAL
	EXCL
	EXCL_EQUAL
	LESS
	LESS_EQUAL
	EQUAL
	EQUAL_EQUAL

	// Literals
	IDENTIFIER
	NUMBER
	NIL
	STRING

	// Keywords
	IF
	ELSE
	NOT
	PLACEHOLDER
	DEFINE
	DECLARE

	ARGS

	AND
	OR

	BOOL_FALSE
	BOOL_TRUE

	// types

	COMMENT
	COMMENT_MULTILINE
	ERR
	EOF
	EOP // end of operation
)

var Keys = map[string]TokenType{
	// single char tokens
	"(":  OP,
	")":  CP,
	"{":  LB,
	"}":  RB,
	"+":  PLUS,
	"/":  SLASH,
	"*":  STAR,
	",":  COMMA,
	".":  DOT,
	"-":  MINUS,
	";":  SEMICOLON,
	":":  COLON,
	"\"": QUOTE,

	">":  GREATER,
	">=": GREATER_EQUAL,
	"!":  EXCL,
	"!=": EXCL_EQUAL,
	"<":  LESS,
	"<=": LESS_EQUAL,
	"=":  EQUAL,
	"==": EQUAL_EQUAL,

	// Keywords
	"_":       PLACEHOLDER,
	"if":      IF,
	"else":    ELSE,
	"define":  DEFINE,
	"declare": DECLARE,

	"args": ARGS,

	"and": AND,
	"or":  OR,

	"false": BOOL_FALSE,
	"true":  BOOL_TRUE,

	"#": COMMENT,

	// reserver for dbg
	"<STRING>":     STRING,
	"<IDENTIFIER>": IDENTIFIER,
	"<NUMBER>":     NUMBER,
	"<NIL>":        NIL,
	"\\0":          EOF,
}

var ReverseKeys = reverseMap(Keys)

func reverseMap(m map[string]TokenType) map[TokenType]string {
	n := make(map[TokenType]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

type Token struct {
	Type  TokenType
	Value string
	Pos   TokenPos
	End   TokenPos
	Line  TokenLine
}

func Print(token Token) {
	tokenTypeReadable, _ := ReverseKeys[token.Type]
	printFormat := "type: %s, value: %s, start: %d, end: %d, line:%d\n"
	fmt.Printf(printFormat, tokenTypeReadable, token.Value, token.Pos, token.End, token.Line)
}

package token

type TokenType int
type Pos int
type Line int

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

var keys = map[string]TokenType{
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

var reverseKeys = reverseMap(keys)

func reverseMap(m map[string]TokenType) map[TokenType]string {
	n := make(map[TokenType]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

var eof = rune(0)

type Token struct {
	tokenType TokenType
	value     string
	pos       Pos
	end       Pos
	line      Line
}

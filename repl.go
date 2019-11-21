package main

import (
	"bufio"
	//"flag"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

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
	STRING

	// Keywords
	IF
	ELSE

	ARGS

	AND
	OR

	FALSE
	TRUE

	NIL
	EOF
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
	"if":   IF,
	"else": ELSE,

	"args": ARGS,

	"and": AND,
	"or":  OR,

	"false": FALSE,
	"true":  TRUE,

	// reserver for dbg
	"<STRING>": STRING,
	"<NUMBER>": NUMBER,
}

var reverseKey = reverseMap(keys)

func reverseMap(m map[string]TokenType) map[TokenType]string {
	n := make(map[TokenType]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

var eof = rune(0)

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\r' }

func isAlpha(ch rune) bool { return unicode.IsLetter(ch) }

func isDigit(ch rune) bool { return unicode.IsDigit(ch) }

func isAlphaNumeric(ch rune) bool { return isAlpha(ch) || isDigit(ch) }

type Token struct {
	tokenType TokenType
	value     string

	pos  Pos
	line Line
}

type Reader struct {
	runeReader io.RuneScanner
}

type Scanner struct {
	pos    Pos
	start  Pos
	reader *Reader
	line   Line
	eof    bool
	// not used stuff is commented
	//startPosition   Pos
	buf bytes.Buffer
}

type Lexer struct {
	scanner *Scanner
	tokens  chan Token // channel of detected tokens
}

func (reader *Reader) read() (ch rune) {
	ch, _, _ = reader.runeReader.ReadRune()
	return ch
}

func (reader *Reader) unread() {
	// maybe err
	reader.runeReader.UnreadRune()
}

func (scanner *Scanner) unread() {
	scanner.pos--
	scanner.reader.unread()
}

func (scanner *Scanner) read() (ch rune) {
	// propagate error
	// handle eof
	scanner.pos++
	return scanner.reader.read()
}

func (scanner *Scanner) reportError() {

}

func (scanner *Scanner) peekRune() (ch rune) {
	ch = scanner.read()
	if ch != eof {
		scanner.unread()
	}
	return ch
}

func (scanner *Scanner) futureMatch(fch rune) bool {
	ch := scanner.read()
	if ch == eof {
		// end of the road
		fmt.Println("unformatted error!")
		return false
	} else if ch == fch {
		return true
	}
	return false
}

// stet function that returns a state function
type stateFunc func(*Lexer) stateFunc

func (lex *Lexer) trimWhitespace() {
	for {
		ch := lex.scanner.read()
		if !isWhitespace(ch) {
			lex.scanner.unread()
			break
		}
	}
}

func (lex *Lexer) extractString() bool {
	lex.scanner.start = lex.scanner.pos
	for {
		ch := lex.scanner.read()
		if ch == '\n' {
			lex.scanner.line++
		} else if ch == eof {
			fmt.Println("ups")
			lex.scanner.buf.Reset()
			return false
		} else if ch == '"' {
			// TODO: can decide if we want to check the next
			// if it something wrong there
			//next := lex.scanner.peek()
			return true
		} else {
		}
		lex.scanner.buf.WriteRune(ch)
	}
}

func (lex *Lexer) extractDigit() {

}

func fullScan(lex *Lexer) stateFunc {
	for {
		switch ch := lex.scanner.read(); ch {
		case ' ':
			lex.trimWhitespace()
		case '(':
			lex.emit(OP)
		case ')':
			lex.emit(CP)
		case '{':
			lex.emit(LB)
		case '}':
			lex.emit(RB)
		case '+':
			lex.emit(PLUS)
		case '-':
			lex.emit(MINUS)
		case '*':
			lex.emit(STAR)
		case ';':
			lex.emit(SEMICOLON)
		case ',':
			lex.emit(DOT)
		case '.':
			lex.emit(DOT)
		// special cases
		case '!':
			if lex.scanner.futureMatch('=') {
				// TODO: need to handle a case wher it is not matched
				lex.emit(EXCL_EQUAL)
			} else {
				lex.emit(EXCL)
			}
			// TODO: case when it is used as NOT
			// lex.emit(NOT)
			//lex.emit(ERROR)
		case '=':
			if lex.scanner.futureMatch('=') {
				lex.emit(EQUAL_EQUAL)
			} else {
				lex.emit(EQUAL)
			}
			//lex.emit(ERROR)
		case '<':
			if lex.scanner.futureMatch('=') {
				lex.emit(LESS_EQUAL)
			} else {
				lex.emit(LESS)
			}
		case '>':
			if lex.scanner.futureMatch('=') {
				lex.emit(GREATER_EQUAL)
			} else {
				lex.emit(GREATER)
			}
		case '\n':
			lex.scanner.line++
		case '"':
			if lex.extractString() {
				lex.emit(STRING)
			} else {
				fmt.Println("ups")
			}
		default:
			if isAlpha(ch) {
				//fmt.Println("is alpha!!!")
				lex.emit(IDENTIFIER)
			}
			if isDigit(ch) {
				//fmt.Println("is digit!!!")
				lex.emit(NUMBER)
			} else {
				// error
			}
			if ch == eof {
				//fmt.Println("it is the end!!!")
				lex.emit(EOF)
			}
		}
	}
	//return nil
}

func (lex *Lexer) emit(tType TokenType) {
	// need more info

	/*
		tokenType TokenType
		value     string

		start Pos
		end   Pos
		line  Line
	*/
	lex.tokens <- Token{
		tokenType: tType,
		pos:       lex.scanner.pos,
		value:     lex.scanner.buf.String(),
	}
	lex.scanner.buf.Reset()
}

func (lex *Lexer) nextToken() Token {
	return <-lex.tokens
}

func (lex *Lexer) run() {
	// sourceText is already lex function that returns 'stateFunc'
	for state := fullScan; state != nil; {
		// asign state and wait for new state
		state = state(lex)
	}
	close(lex.tokens)
}

func startGrinding(input string) (lex *Lexer) {
	lex = &Lexer{
		scanner: &Scanner{
			reader: &Reader{
				runeReader: strings.NewReader(input),
			},
		},
		tokens: make(chan Token),
	}
	go lex.run()
	return lex
}

func readline(idet string, scanner *bufio.Scanner) bool {
	fmt.Print(idet)
	return scanner.Scan()
}

func main() {
	// direct file parser
	/*
		var inputFile string
		flag.StringVar(&inputFile, "file", "", "Input hell file.")
		flag.Parse()

		fmt.Println(inputFile)
			hFile, err := os.Open(inputFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			runeReader := io.RuneReader(bufio.NewReader(hFile))
			fmt.Println(runeReader.ReadRune())
	*/
	// repl

	var buffer []string
	const idet = "hprog> "

	fmt.Println("Hprog Version 0.0.0.0.0.0.0.0.2")
	fmt.Println("One way to escape, ctr-c to exit.")

	scanner := bufio.NewScanner(os.Stdin)

	onNewLine := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return bufio.ScanLines(data, atEOF)
	}

	scanner.Split(onNewLine)
	for {
		for readline(idet, scanner) {
			var sline = scanner.Text()
			if len(sline) > 0 {
				lex := startGrinding(sline)
				for {
					token := lex.nextToken()
					humanToken, found := reverseKey[token.tokenType]
					if found == true {
						fmt.Println(humanToken)
					}
					if token.tokenType == EOF {
						break
					}
				}

				buffer = append(buffer, sline)
			}
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("error: %s\n", scanner.Err())
	}
}

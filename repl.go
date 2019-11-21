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
	NOT
	PLACEHOLDER

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
	"\\0":      EOF,
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

func (scanner *Scanner) peek() (ch rune) {
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
		lex.scanner.pos++
	}
}

func (lex *Lexer) extractString() bool {
	lex.scanner.start = lex.scanner.pos
	for {
		ch := lex.scanner.peek()
		if ch == '\n' {
			lex.scanner.line++
		} else if ch == eof {
			lex.scanner.buf.Reset()
			return false
		} else if ch == '"' {
			// TODO: can decide if we want to check the next
			// if it something wrong there
			//next := lex.scanner.peek()
			return true
		}
		lex.scanner.buf.WriteRune(lex.scanner.read())
	}
}

func (lex *Lexer) extractNumber() bool {
	lex.scanner.start = lex.scanner.pos
	foundPoint := false
	for {
		ch := lex.scanner.read()
		fmt.Println(string(ch))
		if isDigit(ch) {
			lex.scanner.buf.WriteRune(ch)
		} else if ch == '.' {
			// allow 1..0
			if foundPoint == true {
				return false
			}
			foundPoint = true
			lex.scanner.buf.WriteRune(ch)
		} else if ch == eof {
			// TODO: we don't want to have only numbers
			// return false
			return true
		} else {
			// TODO: also after the number more char can be covered
			return true
		}
	}
}

func (lex *Lexer) extractIdentifier() bool {
	lex.scanner.start = lex.scanner.pos
	for {
		ch := lex.scanner.read()
		if isAlphaNumeric(ch) {
			lex.scanner.buf.WriteRune(ch)
		} else if ch == eof {
			// TODO: we don't want to have only numbers
			// return false
			return true
		} else {
			// TODO: also after the number more char can be covered
			return true
		}
	}
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
			if isDigit(lex.scanner.peek()) {
				fmt.Println("ups dot")
			}
			// TODO: special cases not solved: .01
			// TODO: .. range char not sovled
			/*
				if isDigit(lex.scanner.peek()) {
					lex.scanner.unread()
					if lex.extractNumber() {
						lex.emit(NUMBER)
					}
				} else if !isAlpha(lex.scanner.peek()) {
					lex.emit(DOT)
				} else {
					fmt.Println("ups dot")
				}
			*/
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
		case '\t':
		case '\r':
		case '\n':
			lex.scanner.line++
		case '"':
			if lex.extractString() {
				lex.emit(STRING)
			} else {
				fmt.Println("string ups")
			}
		default:
			if isAlpha(ch) {
				if lex.extractIdentifier() {
					lex.emit(IDENTIFIER)
				} else {
					fmt.Println("number ups")
				}
			}
			if isDigit(ch) {
				// we detected the number, now we need to go back
				// from the start to process the whole thing
				lex.scanner.unread()
				if lex.extractNumber() {
					lex.emit(NUMBER)
				} else {
					fmt.Println("number ups")
				}
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
		pos:       lex.scanner.start,
		value:     lex.scanner.buf.String(),
	}
	lex.scanner.buf.Reset()
	lex.scanner.start = lex.scanner.pos
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
						fmt.Println(token, humanToken)
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

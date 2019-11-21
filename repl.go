package main

import (
	"bufio"
	"bytes"
	"flag"
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
	DEFINE
	DECLARE

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
	"_":       PLACEHOLDER,
	"if":      IF,
	"else":    ELSE,
	"define":  DEFINE,
	"declare": DECLARE,

	"args": ARGS,

	"and": AND,
	"or":  OR,

	"false": FALSE,
	"true":  TRUE,

	// reserver for dbg
	//	"<STRING>":     STRING,
	//	"<IDENTIFIER>": IDENTIFIER,
	//	"<NUMBER>":     NUMBER,
	"\\0": EOF,
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

func wannaPrint(token Token) {
	printFormat := "type: %d, value: %s, start: %d, end: %d, line:%d\n"
	fmt.Printf(printFormat, token.tokenType, token.value, token.pos, token.end, token.line)
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\r' }

func isAlpha(ch rune) bool { return unicode.IsLetter(ch) }

func isDigit(ch rune) bool { return unicode.IsDigit(ch) }

func isAlphaNumeric(ch rune) bool { return isAlpha(ch) || isDigit(ch) }

type Token struct {
	tokenType TokenType
	value     string
	pos       Pos
	end       Pos
	line      Line
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
	ch := scanner.peek()
	if ch == eof {
		// end of the road
		fmt.Println("unformatted error!")
		return false
	} else if ch == fch {
		scanner.read()
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

func (lex *Lexer) trimComment() {
	for {
		ch := lex.scanner.read()
		if ch == eof {
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
		} else if ch == '"' {
			// TODO: can decide if we want to check the next
			// if it something wrong there
			return true
		} else if ch == eof {
			return false
		}
		lex.scanner.buf.WriteRune(ch)
	}
}

func (lex *Lexer) extractNumber(ch rune) bool {
	lex.scanner.buf.WriteRune(ch)
	lex.scanner.start = lex.scanner.pos - 1

	foundPoint := false
	for {
		ch := lex.scanner.peek()
		if isDigit(ch) {
			lex.scanner.read()
			lex.scanner.buf.WriteRune(ch)
		} else if ch == '.' {
			// allow 1..0
			if foundPoint == true {
				lex.scanner.buf.Reset()
				return false
			}
			foundPoint = true
			lex.scanner.read()
			lex.scanner.buf.WriteRune(ch)
		} else if isAlpha(ch) {
			lex.scanner.buf.Reset()
			return false
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

func (lex *Lexer) extractIdentifier(ch rune) bool {
	lex.scanner.start = lex.scanner.pos
	lex.scanner.buf.WriteRune(ch)
	for {
		ch = lex.scanner.peek()
		isAlphaNumeric := isAlphaNumeric(ch)
		if isAlphaNumeric {
			lex.scanner.buf.WriteRune(ch)
			lex.scanner.read()
		} else if !isAlphaNumeric {
			return true
		} else if ch == eof {
			lex.scanner.buf.Reset()
			return false
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
		case '_':
			lex.emit(PLACEHOLDER)
		case '/':
			lex.emit(SLASH)
		case '.':
			lex.emit(DOT)
			if isDigit(lex.scanner.peek()) {
				fmt.Println("ups dot")
			}
			// TODO: special cases not solved: .01
			// TODO: .. range char not sovled
			/*
				if isDigit(lex.scanner.peek()) {
					if lex.extractNumber(ch) {
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
				if lex.extractIdentifier(ch) {
					// TODO: the next char should not special
					lex.emit(IDENTIFIER)
				} else {
					fmt.Println("identifier ups")
				}
			}
			if isDigit(ch) {
				// we detected the number, now we need to go back
				// from the start to process the whole thing
				if lex.extractNumber(ch) {
					// TODO: the next char should not special
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
	value := ""
	if lex.scanner.buf.Len() > 0 {
		value = lex.scanner.buf.String()
		tokenType, foundType := keys[value]
		if foundType {
			tType = tokenType
		}
	}
	tokenValue, foundValue := reverseKeys[tType]
	if foundValue {
		value = tokenValue
	}
	lex.tokens <- Token{
		tokenType: tType,
		pos:       lex.scanner.start,
		end:       lex.scanner.pos,
		line:      lex.scanner.line,
		value:     value,
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
	var buffer []string
	var inputFile string
	flag.StringVar(&inputFile, "file", "", "Input hell file.")
	flag.Parse()

	fmt.Println(inputFile)
	hFile, err := os.Open(inputFile)
	if hFile != nil {
		if err != nil {
			fmt.Println(err)
			return
		}
		fileScanner := bufio.NewScanner(hFile)
		for fileScanner.Scan() {
			buffer = append(buffer, fileScanner.Text())
		}

		fmt.Println(strings.Join(buffer[:], "\n"))
		lex := startGrinding(strings.Join(buffer[:], "\n"))

		for {
			token := lex.nextToken()
			//wannaPrint(token)
			fmt.Print(token.value, " ")
			if token.tokenType == EOF {
				break
			}
		}
	} else {
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
						wannaPrint(token)
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
}

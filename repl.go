package main

import (
	"bufio"
	//"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Token int
type Pos int
type Line int

const (
	ILLEGAL Token = iota

	// single char tokens
	OB
	CB
	LB
	RB
	PLUS
	SLASH
	STAR
	COMMA
	DOT
	MINUS

	GREATER
	GREATER_EQUAL
	BANG
	BANG_EQUAL
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

var eof = rune(0)

func isWhitespace(ch rune) bool {
	// NOTE: extend this for other types
	if ch == ' ' {
		return true
	}
	return false
}

func isAlpha(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isAlphaNumeric(ch rune) bool {
	return isAlpha(ch) || isDigit(ch)
}

type Reader struct {
	runeReader io.RuneScanner
	eof        bool
	start      Pos
	current    Pos
	line       Line
}

type Scanner struct {
	reader *Reader
}

func (reader *Reader) read() (ch rune) {
	ch, _, _ = reader.runeReader.ReadRune()
	reader.current++
	return ch
}

func (reader *Reader) unread() {
	// maybe err
	reader.runeReader.UnreadRune()
	reader.current--
}

func (scanner *Scanner) scan() (ch rune) {
	// propagate error
	// handle eof
	return scanner.reader.read()
}

func (scanner *Scanner) unread() {
	scanner.reader.unread()
}

// stet function that returns a state function
type stateFunc func(*lexer) stateFunc

type item struct{}

type lexer struct {
	scanner *Scanner
	items   chan item // channel of detected items
}

func (lex *lexer) trimWhitespace() {
	for {
		ch := lex.scanner.scan()
		if !isWhitespace(ch) {
			lex.scanner.unread()
			break
		}
	}
}

func fullScan(lex *lexer) stateFunc {
	for {
		switch ch := lex.scanner.scan(); ch {
		case ' ':
			fmt.Println("is WS!!!")
			lex.trimWhitespace()
		case '(':
			fmt.Println("is OB!!!")
			lex.emit(OB)
		case ')':
			fmt.Println("is CB!!!")
			lex.emit(CB)
		default:
			if isAlpha(ch) {
				fmt.Println("is alpha!!!")
			}
			if isDigit(ch) {
				fmt.Println("is digit!!!")
			}
			if ch == eof {
				fmt.Println("it is the end!!!")
				lex.emit(EOF)
			}
		}
	}
	return nil
}

func (lex *lexer) emit(t Token) {
	lex.items <- item{}
}

func (lex *lexer) run() {
	// sourceText is already lex function that returns 'stateFunc'
	for state := fullScan; state != nil; {
		// asign state and wait for new state
		state = state(lex)
	}
	close(lex.items)
}

func startGrinding(input string) {
	lex := &lexer{
		scanner: &Scanner{
			reader: &Reader{
				runeReader: strings.NewReader(input),
			},
		},
		items: make(chan item),
	}
	go lex.run()
	lex = nil
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
				startGrinding(sline)
				buffer = append(buffer, sline)
			}
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("error: %s\n", scanner.Err())
	}
}

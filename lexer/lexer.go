package lexer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func reportError(line Line, position Pos, what string) {
	fmt.Fprintf(os.Stderr, "[line:%d, pos:%d] Error, %s\n",
		line, position, what)
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\r' }

func isAlpha(ch rune) bool { return unicode.IsLetter(ch) }

func isDigit(ch rune) bool { return unicode.IsDigit(ch) }

func isAlphaNumeric(ch rune) bool { return isAlpha(ch) || isDigit(ch) }

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

// stet function that returns a state function
type stateFunc func(*Lexer) stateFunc

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

func (scanner *Scanner) peek() (ch rune) {
	ch = scanner.read()
	if ch != eof {
		scanner.unread()
	}
	return ch
}

func (scanner *Scanner) seeFuture(fch rune) bool {
	ch := scanner.peek()
	if ch == eof {
		// end of the road
		reportError(scanner.line, scanner.pos,
			"Reached EOF.")
		return false
	} else if ch == fch {
		scanner.read()
		return true
	}
	return false
}

func (lex *Lexer) walkOnWhitespace() {
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

func (lex *Lexer) extractIdentifier(ch rune) (extracted bool, identifierType TokenType) {
	lex.scanner.start = lex.scanner.pos
	lex.scanner.buf.WriteRune(ch)
	extracted = false
	identifierType = IDENTIFIER
	for {
		ch = lex.scanner.peek()
		isAlphaNumeric := isAlphaNumeric(ch)
		if isAlphaNumeric {
			lex.scanner.buf.WriteRune(ch)
			lex.scanner.read()
		} else if !isAlphaNumeric {
			// idk about this maybe change
			//return true
			extracted = true
			break
		} else if ch == eof {
			lex.scanner.buf.Reset()
			//return false
			reportError(lex.scanner.line, lex.scanner.pos,
				"Reached EOL or OEF.")
			extracted = false
			break
		}
	}
	// get identifier type
	if lex.scanner.buf.Len() > 0 {
		tokenType, foundType := keys[lex.scanner.buf.String()]
		if foundType {
			identifierType = tokenType
		}
	}
	return extracted, identifierType
}

func fullScan(lex *Lexer) stateFunc {
	for {
		switch ch := lex.scanner.read(); ch {
		case ' ':
			lex.walkOnWhitespace()
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
		case '#':
			// TODO: extended to goto EOF
			lex.emit(COMMENT)
		case ':':
			lex.emit(COLON)
		case '!':
			if lex.scanner.seeFuture('=') {
				// TODO: need to handle a case when it is not matched
				lex.emit(EXCL_EQUAL)
			} else {
				lex.emit(EXCL)
			}
			// TODO: case when it is used as NOT
			// lex.emit(NOT)
			//lex.emit(ERROR)
		case '=':
			if lex.scanner.seeFuture('=') {
				lex.emit(EQUAL_EQUAL)
			} else {
				lex.emit(EQUAL)
			}
			// lex.emit(ERR)
		case '<':
			if lex.scanner.seeFuture('=') {
				lex.emit(LESS_EQUAL)
			} else {
				lex.emit(LESS)
			}
		case '>':
			if lex.scanner.seeFuture('=') {
				lex.emit(GREATER_EQUAL)
			} else {
				lex.emit(GREATER)
			}
			lex.scanner.line++
		case '"':
			if lex.extractString() {
				lex.emit(STRING)
			} else {
				reportError(lex.scanner.line, lex.scanner.pos,
					"Wrong string formatting.")
				lex.emit(ERR)
			}
		case '\t':
		case '\r':
		case '\n':
		default:
			if isAlpha(ch) {
				extracted, identifier := lex.extractIdentifier(ch)
				if extracted {
					// TODO: the next char should not special
					lex.emit(identifier)
				} else {
					reportError(lex.scanner.line, lex.scanner.pos,
						"Invalid identifier.")
				}
			} else if isDigit(ch) {
				// we detected the number, now we need to go back
				// from the start to process the whole thing
				if lex.extractNumber(ch) {
					// TODO: the next char should not special
					lex.emit(NUMBER)
				} else {
					reportError(lex.scanner.line, lex.scanner.pos,
						"Invalid identifer.")
					lex.emit(ERR)
				}
			} else if ch == eof {
				lex.emit(EOF)
			} else {
				reportError(lex.scanner.line, lex.scanner.pos,
					"Invalid character.")
				lex.emit(ERR)
			}
		}
	}
	//return nil
}

func (lex *Lexer) emit(tType TokenType) {

	value := lex.scanner.buf.String()
	if len(value) == 0 {
		tokenValue, foundValue := reverseKeys[tType]
		if foundValue {
			value = tokenValue
		}
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

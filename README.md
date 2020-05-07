# Hel programming language
This is a functional language.

## Philosophy
- Code should be readable
- No states, only transformations
- Everything is a function except variables
- Variables are immutable
- Functions can be referenced by aliases

## Lexer
- Change handling of error messages, instead of bools, use a wrapper for the errors
- Parsing numbers
    - No double
    - No negative
- parsing print stuff
    - not implemented
- how to use print?
- 1..10 (ranges)
    - not implemented
- .11 floats starting with dot
    - not allowed

## Parse (todo)
- not even started
- Types

## Supported features
- Conditional statements (not)
- Loops (not)
- Functions (not)
- Types (not)
- Data structures (not)
    - list
    - dict
    - set
- anonymous functions (not)
- function aliases (not)
- lambda functions (not)
- Functional transformations (not)
    - map, filter, reduce
- Structs (not)
- Concurrency (not)

Sub-features:
- Immutability (not)
- Pure functions (not)
- Lazy evaluation?


## Demo

Currently, (== (max -3 -4) 4), translates to:
``````kk`
(main.Grouping) {
 expression: (main.Binary) {
  operator: (main.Token) {
   tokenType: (main.TokenType) 21,
   value: (string) (len=2) "==",
   pos: (main.Pos) 1,
   end: (main.Pos) 3,
   line: (main.Line) 0
  },
  left: (main.Grouping) {
   expression: (main.Binary) {
    operator: (main.Token) {
     tokenType: (main.TokenType) 22,
     value: (string) (len=3) "max",
     pos: (main.Pos) 6,
     end: (main.Pos) 8,
     line: (main.Line) 0
    },
    left: (main.Unary) {
     operator: (main.Token) {
      tokenType: (main.TokenType) 10,
      value: (string) (len=1) "-",
      pos: (main.Pos) 8,
      end: (main.Pos) 10,
      line: (main.Line) 0
     },
     right: (main.Literal) {
      value: (string) (len=1) "3"
     }
    },
    right: (main.Unary) {
     operator: (main.Token) {
      tokenType: (main.TokenType) 10,
      value: (string) (len=1) "-",
      pos: (main.Pos) 11,
      end: (main.Pos) 13,
      line: (main.Line) 0
     },
     right: (main.Literal) {
      value: (string) (len=1) "4"
     }
    }
   }
  },
  right: (main.Literal) {
   value: (string) (len=1) "4"
  }
 }
}
`

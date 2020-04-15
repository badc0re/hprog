package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func readline(idet string, scanner *bufio.Scanner) bool {
	fmt.Print(idet)
	return scanner.Scan()
}

func wannaPrint(token Token) {
	tokenTypeReadable, _ := reverseKeys[token.tokenType]
	printFormat := "type: %s, value: %s, start: %d, end: %d, line:%d\n"
	fmt.Printf(printFormat, tokenTypeReadable, token.value, token.pos, token.end, token.line)
}

func main() {
	// direct file parser
	var buffer []string
	var inputFile string
	flag.StringVar(&inputFile, "file", "", "Input hell file.")
	flag.Parse()

	fmt.Println(inputFile)
	hFile, err := os.Open(inputFile)

	bla()

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
			wannaPrint(token)
			//fmt.Print(token.value, " ")
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
						if token.tokenType == EOF || token.tokenType == ERR {
							break
						}
						wannaPrint(token)
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

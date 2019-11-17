package main

import (
	"bufio"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"hprog/hcore"
	"os"
)

func readline(idet string, scanner *bufio.Scanner) bool {
	fmt.Print(idet)
	return scanner.Scan()
}

type blaListener struct {
	*HCore.BaseHCoreListener
}

func main() {
	const idet = "hprog> "

	fmt.Println("Hprog Version 0.0.0.0.0.0.2")
	fmt.Println("One way to escape, ctr-c to exit.")

	scanner := bufio.NewScanner(os.Stdin)

	onNewLine := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return bufio.ScanLines(data, atEOF)
	}

	scanner.Split(onNewLine)
	for {
		for readline(idet, scanner) {
			var _bline = scanner.Text()
			is := antlr.NewInputStream(_bline)
			lexer := HCore.NewHCoreLexer(is)
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
			p := HCore.NewHCoreParser(stream)
			antlr.ParseTreeWalkerDefault.Walk(&blaListener{}, p.Start())
			// add history
			//buffer = append(buffer, _bline)
			//fmt.Println(buffer)
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("error: %s\n", scanner.Err())
	}
}

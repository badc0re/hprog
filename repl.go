package main

import (
	"bufio"
	"fmt"
	"os"
)

func add_history() {

}

func readline(idet string, scanner *bufio.Scanner) bool {
	fmt.Print(idet)
	return scanner.Scan()
}

func main() {
	var buffer []byte
	const idet = "hprog> "

	fmt.Println("Hprog Version 0.0.0.0.0.0.1")
	fmt.Println("One way to escape, ctr-c to exit.")

	scanner := bufio.NewScanner(os.Stdin)

	onNewLine := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return bufio.ScanLines(data, atEOF)
	}

	//buffer := make([]byte, 12)
	//scanner.Buffer(buffer, bufio.MaxScanTokenSize /* 65536 */)
	scanner.Split(onNewLine)
	for {
		for readline(idet, scanner) {
			var _bline = scanner.Bytes()
			buffer = append(buffer, _bline...)
			fmt.Println(buffer)
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("error: %s\n", scanner.Err())
	}
}

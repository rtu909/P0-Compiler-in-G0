package main

import (
	"bufio"
	"os"
)

func main() {

	// Open a file for buffered reading
	f, err := os.Open("../p0-programs/Fibonacci.p")
	if err != nil {
		panic("Unable to open the requested file")
	}
	reader := bufio.NewReader(f)
	// Make a channel to connect the scanner and the parser
	tokenChannel := make(chan SourceUnit, 5)
	endChannel := make(chan int)
	// start the parser
	ScannerInit(reader, tokenChannel)
	// start the scanner
	compileFile(tokenChannel, endChannel, "wat")
	<-endChannel
}

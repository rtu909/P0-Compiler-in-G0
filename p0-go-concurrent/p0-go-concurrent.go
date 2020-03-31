package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	if strings.HasSuffix(sourceFilePath, ".p") {
		var destinationFilePath = sourceFilePath[:len(sourceFilePath)-3] + ".s"
	} else {
		fmt.Printf(".p file extension expected")
		panic(nil)
	}
	<-endChannel
}

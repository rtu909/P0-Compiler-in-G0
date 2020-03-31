package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {

	sourceFilePath := "../p0-programs/Fibonacci.p"
	var destFilePath string

	// Open a file for buffered reading
	if strings.HasSuffix(sourceFilePath, ".p") {
		destFilePath = sourceFilePath[:len(sourceFilePath)-3] + ".s"
	} else {
		panic(".p file extension expected")
	}

	f, err := os.Open(sourceFilePath)
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
	go compileFile(tokenChannel, endChannel, destFilePath, "wat")
	<-endChannel
}

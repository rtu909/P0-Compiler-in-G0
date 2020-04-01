package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	args := os.Args[1:]

	if len(args) > 1 {
		panic("ERROR: Expected one command line argument but got " + string(len(args)))
		os.Exit(1)
	}
	sourceFilePath := args[0]
	if !isAccessibleFile(sourceFilePath) {
		panic("ERROR: Accessing: " + sourceFilePath + " failed. Ensure that it exists and read & write permissions are enabled.")
		os.Exit(1)
	}

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
	tokenChannel := make(chan SourceUnit)
	endChannel := make(chan int)
	// start the parser
	go ScannerInit(reader, tokenChannel)
	// start the scanner
	go compileFile(tokenChannel, endChannel, destFilePath, "wat")
	<-endChannel
}

func isAccessibleFile(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true // Path to an existing file
	}
	var contents []byte
	err = ioutil.WriteFile(filePath, contents, 0644) // 0644 = -rw-r--r--
	if err == nil {
		os.Remove(filePath)
		return false // Path to an existing accessible directory
	}
	return false // Failed to access
}

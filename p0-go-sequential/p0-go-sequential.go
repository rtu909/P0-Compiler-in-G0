package main

import (
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		panic("ERROR: Expected one command line argument but got " + string(len(args)))
		os.Exit(1)
	}
	sourceFilePath := args[0]
	if !isAccessibleFile(sourceFilePath) {
		panic("ERROR: " + sourceFilePath + " is not the path to an existing file.")
		os.Exit(1)
	}

	//var destFilePath string
	//
	//// Open a file for buffered reading
	//if strings.HasSuffix(sourceFilePath, ".p") {
	//	destFilePath = sourceFilePath[:len(sourceFilePath)-3] + ".s"
	//} else {
	//	panic(".p file extension expected")
	//}

	compileFile(sourceFilePath, "mips")
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

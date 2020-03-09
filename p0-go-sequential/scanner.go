package main
import (
	"fmt"
)
//symbols as integer constants
var TIMES = 1
var DIV = 2
var MOD = 3
var AND = 4
var PLUS = 5
var MINUS = 6
var OR = 7
var EQ = 8
var NE = 9
var LT = 10
var GT = 11
var LE = 12
var GE = 13
var PERIOD = 14
var COMMA = 15
var COLON = 16
var RPAREN = 17
var RBRAK = 18
var OF = 19
var THEN = 20
var DO = 21
var LPAREN = 22
var LBRAK = 23
var NOT = 24
var BECOMES = 25
var NUMBER = 26
var IDENT = 27
var SEMICOLON = 28
var END = 29
var ELSE = 30
var IF = 31
var WHILE = 32
var ARRAY = 33
var RECORD = 34
var CONST = 35
var TYPE = 36
var VAR = 37
var PROCEDURE = 38
var BEGIN = 39
var PROGRAM = 40
var EOF = 41

//global variables
var line, lastline, errline int
var pos, lastpos, errpos int
var sym, val interface{}
var error bool
var source, ch string
var index int

//initialization of the scanner
//source is string
func initial(src string) {
	line, lastline, errline = 1,1,1
	pos, lastpos, errpos = 0,0,0
	sym, val, error, source, index = nil, nil, false, src, 0
	getChar(); getSym()
}

//assigns the next character in ch
//variables line, pos are updated with the current location in source
//lastline, lastpos are updated with location of previously read character
func getChar(){
	if index == len(source){
		ch = string(0) //equivalent to chr(0), converts 0 to UTF=8 string
	} else {
		ch, index = string(source[index]), index + 1
		lastpos = pos
		if ch == string('\n'){
			pos, line = 0, line + 1
		} else {
			lastline, pos = line, pos + 1
		}
	}
}
//prints error message with current location in the source
func mark(msg string){
	if (lastline > errline) || (lastpos > errpos){
		fmt.Println("error: line", lastline, "pos", lastpos, msg)
	}
	errline, errpos, error = lastline, lastpos, true
}

func number(){

}

func identKW(){

}

func comment(){

}
//recognizes the next symbol and assigns it to the variables sym and val
func getSym() {

}
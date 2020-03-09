package main
import (
	"fmt"
	"strconv"
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

//sets sym to NUMBER and assigns NUMBER to val
func number(){
	sym, val = NUMBER, 0
	for "0" <= ch && ch <= "9"{
		val, _ = strconv.Atoi(ch)
		tempVal := 10*val.(int)
		val = val.(int) + tempVal //weird stuff, check this
		getChar()
	}
	if val.(int) >= 2^31{
		mark("number too large"); val = 0
	}
}

var KEYWORDS = map[string]int{
	"div": DIV, "mod": MOD, "and": AND, "or": OR, "of": OF, "then": THEN, "do": DO, "not": NOT,
	"end": END, "else": ELSE, "if": IF, "while": WHILE, "array": ARRAY, "record": RECORD,
	"const": CONST, "type": TYPE, "var": VAR, "procedure": PROCEDURE, "begin": BEGIN, "program": PROGRAM,
}

func identKW(){
	start := index - 1
	for ("A" <= ch && ch <= "Z") || ("a" <= ch && ch <= "z") || ("0" <= ch && ch <= "9"){
		getChar()
	}
	val = source[start:index-1]
	var exists bool
	sym, exists = KEYWORDS[val.(string)]
	//if val is not in KEYWORDS dictionary, then sym is IDENT
	if (!exists){
		sym = IDENT
	}
}

func comment(){
	for (string(0) != ch) && (ch != "}"){
		getChar()
	}
	if ch == string(0){
		mark("comment not terminated")
	} else {
		getChar()
	}
}
//recognizes the next symbol and assigns it to the variables sym and val
func getSym() {
	for (string(0) < ch) && (ch <= " "){
		getChar()
	}
	if ("A" <= ch) && (ch <= "Z") || ("a" <= ch) && (ch <= "z"){
		identKW()
	} else if ("0" <= ch) && (ch <= "9"){
		number()
	} else if (ch == "{"){
		comment(); getSym()
	} else if (ch == "*"){
		getChar(); sym = TIMES
	} else if (ch == "+"){
		getChar(); sym = PLUS
	} else if (ch == "-"){
		getChar(); sym = MINUS
	} else if (ch == "="){
		getChar(); sym = EQ
	} else if (ch == "<"){
		getChar()
		if (ch == "="){
			getChar(); sym = LE
		} else if (ch == ">"){
			getChar(); sym = NE
		} else {
			sym = LT
		}
	} else if (ch == ">"){
		getChar()
		if (ch == "="){
			getChar(); sym = GE
		} else {
			sym = GT
		}
	} else if (ch == ";"){
		getChar(); sym = SEMICOLON
	} else if (ch == ","){
		getChar(); sym = COMMA
	} else if (ch == ":"){
		getChar()
		if (ch == "="){
			getChar(); sym = BECOMES
		} else {
			sym = COLON
		}
	} else if (ch == "."){
		getChar(); sym = PERIOD
	} else if (ch == "("){
		getChar(); sym = LPAREN
	} else if (ch == ")"){
		getChar(); sym = RPAREN
	} else if (ch == "["){
		getChar(); sym = LBRAK
	} else if (ch == string(0)){
		getChar(); sym = EOF
	} else{
		mark("illegal character"); getChar(); sym = nil
	}
}
package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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
var sym int
var val interface{}
var error bool
var ch string
var reader *bufio.Reader

var parserChannel chan SourceUnit

//initialization of the scanner
//source is string
func ScannerInit(r *bufio.Reader, pc chan SourceUnit) {
	parserChannel = pc
	line, lastline, errline = 1, 1, 1
	pos, lastpos, errpos = 0, 0, 0
	sym, val, error, reader = 0, nil, false, r
	getChar()

	for {
		getSym()
		if ch == string(0) {
			return
		}
	}
}

type SourceUnit struct {
	sym int
	val interface{}
}

//assigns the next character in ch
//variables line, pos are updated with the current location in source
//lastline, lastpos are updated with location of previously read character
func getChar() {
	rune, _, err := reader.ReadRune()
	if err == io.EOF {
		ch = string(0) //equivalent to chr(0), converts 0 to UTF=8 string
	} else {
		ch = string(rune)
		lastpos = pos
		if ch == string('\n') {
			pos, line = 0, line+1
		} else {
			lastline, pos = line, pos+1
		}
	}
}

//prints error message with current location in the source
func mark(msg string) {
	if (lastline > errline) || (lastpos > errpos) {
		fmt.Println("error: line", lastline, "pos", lastpos, msg)
	}
	errline, errpos, error = lastline, lastpos, true
}

//sets sym to NUMBER and assigns NUMBER to val
func number() {
	sym, val = NUMBER, 0
	for "0" <= ch && ch <= "9" {
		var tempVal int
		tempVal, _ = strconv.Atoi(ch)
		val = tempVal + 10*val.(int) //weird stuff, check this
		getChar()
	}
	//fmt.Println(int(math.Pow(2, 31)))
	//fmt.Println(val.(int))
	if val.(int) >= int(math.Pow(2, 31)) {
		mark("number too large")
		val = 0
	}
}

var KEYWORDS = map[string]int{
	"div": DIV, "mod": MOD, "and": AND, "or": OR, "of": OF, "then": THEN, "do": DO, "not": NOT,
	"end": END, "else": ELSE, "if": IF, "while": WHILE, "array": ARRAY, "record": RECORD,
	"const": CONST, "type": TYPE, "var": VAR, "procedure": PROCEDURE, "begin": BEGIN, "program": PROGRAM,
}

func identKW() {
	var valStr = ""
	for ("A" <= ch && ch <= "Z") || ("a" <= ch && ch <= "z") || ("0" <= ch && ch <= "9") {
		valStr += ch
		getChar()
	}
	val = valStr
	//fmt.Println(val)
	var exists bool
	sym, exists = KEYWORDS[val.(string)]
	//fmt.Println(sym)
	//if val is not in KEYWORDS dictionary, then sym is IDENT
	if !exists {
		//fmt.Println("didn't work")
		sym = IDENT
	}
	//fmt.Println(sym, val)
}

func comment() {
	for (string(0) != ch) && (ch != "}") {
		getChar()
	}
	if ch == string(0) {
		mark("comment not terminated")
	} else {
		getChar()
	}
}

//recognizes the next symbol and assigns it to the variables sym and val
func getSym() {
	for (string(0) < ch) && (ch <= " ") {
		getChar()
	}
	if ("A" <= ch) && (ch <= "Z") || ("a" <= ch) && (ch <= "z") {
		identKW()
	} else if ("0" <= ch) && (ch <= "9") {
		number()
	} else if ch == "{" {
		comment()
		return
	} else if ch == "*" {
		getChar()
		sym = TIMES
	} else if ch == "+" {
		getChar()
		sym = PLUS
	} else if ch == "-" {
		getChar()
		sym = MINUS
	} else if ch == "=" {
		getChar()
		sym = EQ
	} else if ch == "<" {
		getChar()
		if ch == "=" {
			getChar()
			sym = LE
		} else if ch == ">" {
			getChar()
			sym = NE
		} else {
			sym = LT
		}
	} else if ch == ">" {
		getChar()
		if ch == "=" {
			getChar()
			sym = GE
		} else {
			sym = GT
		}
	} else if ch == ";" {
		getChar()
		sym = SEMICOLON
	} else if ch == "," {
		getChar()
		sym = COMMA
	} else if ch == ":" {
		getChar()
		if ch == "=" {
			getChar()
			sym = BECOMES
		} else {
			sym = COLON
		}
	} else if ch == "." {
		getChar()
		sym = PERIOD
		//fmt.Println("period")
	} else if ch == "(" {
		getChar()
		sym = LPAREN
		//fmt.Println("parentheses")
	} else if ch == ")" {
		getChar()
		sym = RPAREN
	} else if ch == "[" {
		getChar()
		sym = LBRAK
	} else if ch == "]" {
		getChar()
		sym = RBRAK
	} else if ch == string(0) {
		getChar()
		sym = EOF
	} else {
		mark("illegal character")
		getChar()
		sym = 0
	}
	fmt.Println(sym, val)
	result := SourceUnit{sym, val}
	parserChannel <- result
}

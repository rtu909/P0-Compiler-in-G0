package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var FIRSTFACTOR = [4]int{IDENT, NUMBER, LPAREN, NOT}
var FOLLOWFACTOR = [22]int{TIMES, DIV, MOD, AND, OR, PLUS, MINUS, EQ, NE, LT, LE, GT, GE, COMMA, SEMICOLON, THEN, ELSE,
	RPAREN, RBRAK, DO, PERIOD, END}
var FIRSTEXPRESSION = [6]int{PLUS, MINUS, IDENT, NUMBER, LPAREN, NOT}
var FIRSTSTATEMENT = [4]int{IDENT, IF, WHILE, BEGIN}
var FOLLOWSTATEMENT = [3]int{SEMICOLON, END, ELSE}
var FIRSTTYPE = [4]int{IDENT, RECORD, ARRAY, LPAREN}
var FOLLOWTYPE = [1]int{SEMICOLON}
var FIRSTDECL = [4]int{CONST, TYPE, VAR, PROCEDURE}
var FOLLOWDECL = [1]int{BEGIN}
var FOLLOWPROCCALL = [3]int{SEMICOLON, END, ELSE}
var STRONGSYMS = [8]int{CONST, TYPE, VAR, PROCEDURE, WHILE, IF, BEGIN, EOF}

func selector(x Entry) Entry {
	var a = [2]int{PERIOD, LBRAK}
	for doesContain(a[:], sym) {
		if sym == PERIOD { // x.f
			getSym()
			if sym == IDENT {
				asRec, isRec := x.GetP0Type().(*P0Record)
				if isRec {
					for f := range asRec.GetFields() {
						if f == val {
							// x = CG.genSelect(x, f);
							break
						} else {
							mark("not a field")
						}
					}
					getSym()
				} else {
					mark("not a record")
				}
			} else {
				mark("identifier expected")
			}
		} else { // x[y]
			getSym()
			var y = expression()
			xAsArray, xIsArray := x.(*P0Array)
			if xIsArray {
				_, yIsInt := y.GetP0Type().(*P0Int)
				if yIsInt {
					var lowerbound = xAsArray.GetLowerBound()
					yAsConst, castSucceed := y.(*P0Const)
					if castSucceed && yAsConst.GetValue().(int) < lowerbound || yAsConst.GetValue().(int) >= lowerbound+xAsArray.GetLength() {
						mark("index out of bounds")
					} else {
						// x = CG.genIndex(x, y)
					}
				} else {
					mark("index not integer")
				}
			} else {
				mark("not an array")
			}
			if sym == RBRAK {
				getSym()
			} else {
				mark("] expected")
			}
		}
	}
	return x
}

func factor() {
	if !doesContain(FIRSTFACTOR[:], sym) {
		mark("expression expected")
		for !(doesContain(FOLLOWFACTOR[:], sym) || doesContain(STRONGSYMS[:], sym) ||
			doesContain(FIRSTFACTOR[:], sym)) {
			getSym()
		}
	}
	if sym == IDENT {
		// CONTINUE FROM HERE
	}
}

// TODO: IMPLEMENT
func expression() Entry {
	return nil
}

// P0Primitive is an enumerated type that represents one of the built-in types in P0.
// It is only meant to represent the base types; composite types are represented in P0Type
type P0Target int

const (
	Wat P0Target = iota
	Mips
)

func compileFile(sourceFilePath string, target string) {
	if strings.HasSuffix(sourceFilePath, ".p") {
		var fileData, fileOpenError = ioutil.ReadFile(sourceFilePath)
		panicIfError(fileOpenError)
		var sourceCode = string(fileData)
		var destinationFilePath = sourceFilePath[:len(sourceFilePath)-3] + ".s"
		compileString(sourceCode, destinationFilePath, toP0Target(target))
	} else {
		fmt.Printf(".p file extension expected")
		panic(nil)
	}
}

func toP0Target(target string) P0Target {
	switch target {
	case "wat":
		return Wat
	case "mips":
		return Mips
	default:
		fmt.Printf("target does not exist")
		panic(nil)
	}
}

func panicIfError(e interface{}) {
	if e != nil {
		panic(e)
	}
}

func compileString(sourceCode string, destinationFilePath string, target P0Target) {
	switch target {
	case Wat:
		// Prepare
	case Mips:
		// Prepare
	default:
		fmt.Printf("target recognized but is not supported")
		panic(nil)
	}
	ScannerInit(sourceCode)
	st := new(SliceMapSymbolTable)
	st.Init()
}

func doesContain(elements []int, e int) bool {
	for _, a := range elements {
		if a == e {
			return true
		}
	}
	return false
}

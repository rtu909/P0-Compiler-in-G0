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

// TODO: put into a struct?
var st SymbolTable
var cg CodeGenerator

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

func compoundStatement() {
	// TODO:
	// Doesn't need to return anything; the result is unused
}

func statement() Entry {
	// TODO:
	return nil
}

func typ() P0Type {
	// TODO:
	return nil
}

func typedIds(kind func(P0Type) P0Type) {
	// TODO:
}

func declarations() int {
	if !(doesContain(FIRSTDECL[:], sym) || doesContain(FOLLOWDECL[:], sym)) {
		mark("'begin' or declaration expected")
		for !(doesContain(FIRSTDECL[:], sym) || doesContain(FOLLOWDECL[:], sym) || doesContain(STRONGSYMS[:], sym)) {
			getSym()
		}
	}
	for sym == CONST {
		getSym()
		if sym == IDENT {
			ident := val.(string)
			getSym()
			getElseMark(sym == EQ, "= expected")
			x := expression()
			_, xIsConst := x.(*P0Const)
			if xIsConst {
				st.NewDecl(ident, x)
			} else {
				mark("expression not constant")
			}
		} else {
			mark("constant name expected")
		}
		getElseMark(sym == SEMICOLON, "; expected")
	}
	for sym == TYPE {
		getSym()
		if sym == IDENT {
			ident := val.(string)
			getSym()
			getElseMark(sym == EQ, "= expected")
			x := typ()
			st.NewDecl(ident, x)
			getElseMark(sym == SEMICOLON, "; expected")
		} else {
			mark("type name expected")
		}
	}
	start := len(st.TopScope())
	for sym == VAR {
		getSym()
		typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type, "", 0, "", 0, 0} })
		getElseMark(sym == SEMICOLON, "; expected")
	}
	cg.GenGlobalVars(st.TopScope(), start)
	return 0 // TODO:
}

func program() string {
	st.NewDecl("boolean", cg.GenBool(&P0Bool{}))
	st.NewDecl("integer", cg.GenBool(&P0Int{}))
	st.NewDecl("true", &P0Const{&P0Bool{}, "", 0, 1})
	st.NewDecl("false", &P0Const{&P0Bool{}, "", 0, 0})
	st.NewDecl("read", &P0StdProc{nil, "", 0, []P0Type{&P0Ref{&P0Int{}, "", 0, "", 0, 0}}})
	st.NewDecl("write", &P0StdProc{nil, "", 0, []P0Type{&P0Var{&P0Int{}, "", 0, "", 0, 0}}})
	st.NewDecl("writeln", &P0StdProc{nil, "", 0, []P0Type{}})
	cg.GenProgStart()
	getElseMark(sym == PROGRAM, "'program expected")
	// The original program actually accessed the program name here
	getElseMark(sym == IDENT, "Program name expected")
	getElseMark(sym == SEMICOLON, "; expected")
	declarations()
	cg.GenProgEntry( /*ident*/ ) // ident was passed in the og P0 compiler, but it is not used so we removed it
	compoundStatement()
	return cg.GenProgExit()
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
		cg = &WasmGenerator{}
	case Mips:
		// Prepare
	default:
		fmt.Printf("target recognized but is not supported")
		panic(nil)
	}
	ScannerInit(sourceCode)
	st = new(SliceMapSymbolTable)
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

// If the predicate is true, getSym is called. Otherwise, mark is called with the message
// I introduced to try to make the code easier to read; if it has the opposite effect let me know - David
func getElseMark(predicate bool, markMessage string) {
	if predicate {
		getSym()
	} else {
		mark(markMessage)
	}
}

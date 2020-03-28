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
					for _, f := range asRec.GetFields() {
						if f == val {
							x = cg.GenSelect(x, f)
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
						x = cg.GenIndex(x, y)
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

func factor() Entry {
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

func term() Entry {
	var x = factor()

	for sym == TIMES || sym == DIV || sym == MOD || sym == AND {
		var op = sym
		getSym()

		_, xIsConst1 := x.(*P0Const)
		if op == AND && !xIsConst1 {
			x = cg.GenUnaryOp(AND, x)
		}
		xAsConst, xIsConst := x.(*P0Const)
		_, xIsInt := x.GetP0Type().(*P0Int)
		_, xIsBool := x.GetP0Type().(*P0Bool)

		var y = factor()
		_, yIsBool := y.GetP0Type().(*P0Bool)
		_, yIsInt := y.GetP0Type().(*P0Int)
		yAsConst, yIsConst := y.(*P0Const)

		if xIsInt && yIsInt {
			if xIsConst && yIsConst {
				if op == TIMES {
					xAsConst.SetValue(xAsConst.GetValue().(int) * yAsConst.GetValue().(int))
				} else if op == DIV {
					xAsConst.SetValue(xAsConst.GetValue().(int) / yAsConst.GetValue().(int))
				} else if op == MOD {
					xAsConst.SetValue(xAsConst.GetValue().(int) % yAsConst.GetValue().(int))
				}
			} else {
				x = cg.GenBinaryOp(op, x, y)
			}
		} else if xIsBool && yIsBool {
			if xIsConst {
				// if x false, x = y
				if xAsConst.GetValue().(int) == 1 {
					xAsConst.SetValue(yAsConst.GetValue().(int))
				}
			} else {
				x = cg.GenBinaryOp(AND, x, y)
			}
		} else {
			mark("bad type")
		}
	}
	return x
}

func simpleExpression() Entry {
	var x Entry
	if sym == PLUS {
		getSym()
		x = term()
	} else if sym == MINUS {
		getSym()
		x = term()
		_, xIsInt := x.GetP0Type().(*P0Int)
		if !xIsInt {
			mark("Bad type")
		}
	} else {
		x = term()

	}

	for sym == PLUS || sym == MINUS || sym == OR {
		var op = sym
		getSym()

		_, xIsConst1 := x.(*P0Const)
		if op == OR && !xIsConst1 {
			x = cg.GenUnaryOp(op, x)
		}
		var y = term()

		_, xIsInt := x.GetP0Type().(*P0Int)
		_, yIsInt := y.GetP0Type().(*P0Int)
		_, xIsBool := x.GetP0Type().(*P0Bool)
		_, yIsBool := y.GetP0Type().(*P0Bool)
		xAsConst, xIsConst := x.(*P0Const)
		yAsConst, yIsConst := y.(*P0Const)

		if xIsInt && yIsInt && (op == PLUS || op == MINUS) {
			if xIsConst && yIsConst {
				if op == PLUS {
					// x = x + y
					xAsConst.SetValue(xAsConst.GetValue().(int) + yAsConst.GetValue().(int))
				} else if op == MINUS {
					// x = x - y
					xAsConst.SetValue(xAsConst.GetValue().(int) - yAsConst.GetValue().(int))
				}
			} else {
				x = cg.GenBinaryOp(op, x, y)
			}
		} else if xIsBool && yIsBool && op == OR {

			if xIsConst {
				// if x false, x = y
				if xAsConst.GetValue().(int) == 0 {
					xAsConst.SetValue(yAsConst.GetValue().(int))
				}

			} else {
				x = cg.GenBinaryOp(OR, x, y)
			}
		} else {
			mark("bad type")
		}
	}
	return x
}

func expression() Entry {
	x := simpleExpression()
	for sym == EQ || sym == NE || sym == LT || sym == LE || sym == GT || sym == GE {
		op := sym
		getSym()
		y := simpleExpression()
		_, xIsBool := x.GetP0Type().(*P0Bool)
		_, yIsBool := y.GetP0Type().(*P0Bool)
		_, xIsInt := x.GetP0Type().(*P0Int)
		_, yIsInt := y.GetP0Type().(*P0Int)
		if (xIsInt && yIsInt) || (xIsBool && yIsBool) {
			xAsConst, xIsConst := x.(*P0Const)
			yAsConst, yIsConst := y.(*P0Const)
			if xIsConst && yIsConst {
				// Perform some constant folding
				// Useful conversion function
				var bool2int func(bool) int = func(predicate bool) int {
					if predicate {
						return 1
					}
					return 0
				}
				switch op {
				case EQ:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) == yAsConst.GetValue().(int)))
					break
				case NE:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) != yAsConst.GetValue().(int)))
					break
				case LT:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) < yAsConst.GetValue().(int)))
					break
				case LE:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) <= yAsConst.GetValue().(int)))
					break
				case GT:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) > yAsConst.GetValue().(int)))
					break
				case GE:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) >= yAsConst.GetValue().(int)))
					break
				}
				xAsConst.p0type = cg.GenBool(&P0Bool{})
				x = xAsConst
			} else {
				x = cg.GenRelation(op, x, y)
			}
		} else {
			mark("bad type")
		}
	}
	return x
}

func compoundStatement() Entry {
	getElseMark(sym == BEGIN, "'begin' expected")
	x := statement()
	for sym == SEMICOLON || doesContain(FIRSTSTATEMENT[:], sym) {
		getElseMark(sym == SEMICOLON, "; missing")
		y := statement()
		cg.GenSeq(x, y) // This returns a value in p0 to build the AST string representation; we just set x to nil here (same effect)
		x = nil
	}
	getElseMark(sym == END, "'end' expected")
	return x
}

// statement does a few things. It's a pretty cool function, you should take a look at the source code.
func statement() Entry {
	var x Entry
	if !doesContain(FIRSTSTATEMENT[:], sym) {
		mark("statement expected")
		for !doesContain(FIRSTSTATEMENT[:], sym) && !doesContain(FOLLOWSTATEMENT[:], sym) && !doesContain(STRONGSYMS[:], sym) {
			getSym()
		}
	}
	if sym == IDENT {
		x = st.Find(val.(string))
		getSym()
		switch x.(type) {
		case *P0Var:
		case *P0Int:
			x = cg.GenVar(x)
			x = selector(x)
			if sym == BECOMES {
				getSym()
				y := expression()
				_, xIsBool := x.(*P0Bool)
				_, yIsBool := y.(*P0Bool)
				_, xIsInt := x.(*P0Int)
				_, yIsInt := y.(*P0Int)
				if (xIsBool && yIsBool) || (xIsInt && yIsInt) {
					cg.GenAssign(x, y)
					x = nil // FIXME: ?
				} else {
					mark("incompatible assignment")
				}
			} else if sym == EQ {
				mark(":= expected")
				getSym()
				_ = expression() // We parse to consume the input, but we can't use the result because the code in incorrect
			} else {
				mark(":= expected")
			}
			break
		case *P0Proc:
		case *P0StdProc:
			// This man codes 8 lines of Go in one line of python
			var fp []P0Type
			var y Entry
			xAsProc, xIsProc := x.(*P0Proc)
			if xIsProc {
				fp = xAsProc.GetParameters()
			} else {
				fp = x.(*P0StdProc).GetParameters()
			}
			ap := make([]Entry, 0)
			i := 0
			if sym == LPAREN {
				getSym()
				if doesContain(FIRSTEXPRESSION[:], sym) {
					y = expression()
					if i < len(fp) {
						if fp[i] == y.GetP0Type() { // TODO: How to do this properly in Go?
							if xIsProc {
								cg.GenActualPara(y, fp[i], i)
							}
						} else {
							mark("illegal parameter mode")
						}
					} else {
						mark("extra parameter")
					}
					i++
					for sym == COMMA {
						getSym()
						y = expression()
						if i < len(fp) {
							if fp[i] == y.GetP0Type() { // TODO: How to do this properly in Go?
								if xIsProc {
									cg.GenActualPara(y, fp[i], i)
								}
							} else {
								mark("illegal parameter mode")
							}
						} else {
							mark("extra parameter")
						}
						i++
					}
					getElseMark(sym == RPAREN, "')' expected")
				}
				if i < len(fp) {
					mark("too few parameters")
				} else if !xIsProc { // x is P0StdProc
					if x.GetName() == "read" {
						cg.GenRead(y) // TODO: continue from here
					} else if x.GetName() == "write" {

					}
				}
			}
			break
		default:
			mark("variable or procedure expected")
			break
		}
	}
	return nil
}

func typ() P0Type {
	var typeToReturn P0Type
	if !doesContain(FIRSTTYPE[:], sym) {
		mark("type expected")
		for !(doesContain(FIRSTTYPE[:], sym) || doesContain(FOLLOWTYPE[:], sym) || doesContain(STRONGSYMS[:], sym)) {
			getSym()
		}
	}
	if sym == IDENT {
		ident := val.(string)
		typeToReturn, _ = st.Find(ident).(P0Type)
		getSym()
	} else if sym == ARRAY {
		getSym()
		getElseMark(sym == LBRAK, "'[' expected")
		x := expression()
		getElseMark(sym == PERIOD, "'.' expected")
		getElseMark(sym == PERIOD, "'.' expected")
		y := expression()
		getElseMark(sym == RBRAK, "']' expected")
		getElseMark(sym == OF, "'of' expected")
		z := typ()
		xAsConst, xIsConst := x.(*P0Const)
		yAsConst, yIsConst := y.(*P0Const)
		if !xIsConst || xAsConst.GetValue().(int) < 0 {
			mark("bad lower bound")
			typeToReturn = nil
		} else if !yIsConst || yAsConst.GetValue().(int) <= xAsConst.GetValue().(int) {
			mark("bad upper bound")
			typeToReturn = nil
		} else {
			typeToReturn = cg.GenArray(&P0Array{
				base:   z,
				lower:  xAsConst.GetValue().(int),
				length: yAsConst.GetValue().(int) - xAsConst.GetValue().(int) + 1,
			})
		}
	} else if sym == RECORD {
		getSym()
		st.OpenScope()
		typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type: p0type} })
		for sym == SEMICOLON {
			getSym()
			typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type: p0type} })
		}
		getElseMark(sym == END, "'end' expected")
		r := st.TopScope()
		rCasted := make([]P0Type, 0)
		for _, val := range r {
			rCasted = append(rCasted, val.(P0Type))
		}
		st.CloseScope()
		typeToReturn = cg.GenRecord(&P0Record{fields: rCasted})
	} else {
		typeToReturn = nil
	}
	return typeToReturn
}

func typedIds(kind func(P0Type) P0Type) {
	var tid []string
	if sym == IDENT {
		tid = make([]string, 1)
		tid[0] = val.(string)
	} else {
		mark("identifier expected")
		tid = make([]string, 0)
	}
	for sym == COMMA {
		getSym()
		if sym == IDENT {
			tid = append(tid, val.(string))
			getSym()
		} else {
			mark("identifier expected")
		}
	}
	if sym == COLON {
		getSym()
		tp := typ()
		if tp != nil {
			for _, attrName := range tid {
				st.NewDecl(attrName, kind(tp))
			}
		}
	} else {
		mark("':' expected")
	}
}

func declarations(generatorFunc func(declaredVars []Entry, start int) int) int {
	var varsize int
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
	varsize = generatorFunc(st.TopScope(), start)
	for sym == PROCEDURE {
		getSym()
		getElseMark(sym == IDENT, "procedure name expected")
		ident := val.(string)
		st.NewDecl(ident, &P0Proc{nil, "", 0, nil})
		st.OpenScope()
		var fp []Entry
		if sym == LPAREN {
			getSym()
			if sym == VAR || sym == IDENT {
				if sym == VAR {
					getSym()
					typedIds(func(p0type P0Type) P0Type { return &P0Ref{p0type, "", 0, "", 0, 0} })
				} else {
					typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type, "", 0, "", 0, 0} })
				}
				for sym == SEMICOLON {
					if sym == VAR {
						getSym()
						typedIds(func(p0type P0Type) P0Type { return &P0Ref{p0type, "", 0, "", 0, 0} })
					} else {
						typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type, "", 0, "", 0, 0} })
					}
				}
			} else {
				mark("Formal parameters expected")
			}
			// The function parameters are stored in the top scope. Make a copy for the symbol table declaration
			fp = st.TopScope()
			tmp := make([]Entry, len(fp))
			copy(tmp, fp)
			st.Find(ident).(*P0Proc).parameters = tmp
			getElseMark(sym == RPAREN, ") expected")
		} else {
			fp = make([]Entry, 0)
		}
		parsize := cg.GenProcStart(ident, fp)
		getElseMark(sym == SEMICOLON, "; expected")
		localsize := declarations(cg.GenLocalVars)
		cg.GenProcEntry(ident, parsize, localsize)
		var x Entry = compoundStatement()
		cg.GenProcExit(x, parsize, localsize)
		st.CloseScope()
		getElseMark(sym == SEMICOLON, "; expected")
	}
	return varsize
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
	declarations(cg.GenGlobalVars)
	cg.GenProgEntry( /*ident*/ ) // ident was passed in the og P0 compiler, but it is not used so we removed it
	compoundStatement()
	return cg.GenProgExit()
}

func compileString(sourceCode string, destinationFilePath string, target P0Target) {
	switch target {
	case Wat:
		// Prepare
		cg = &WasmGenerator{}
	case Mips:
		// Prepare
		cg = &CGmips{}
	default:
		panic("target recognized, but it is not supported")
	}
	ScannerInit(sourceCode)
	st = new(SliceMapSymbolTable)
	st.Init()
}

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

// P0Primitive is an enumerated type that represents one of the built-in types in P0.
// It is only meant to represent the base types; composite types are represented in P0Type
type P0Target int

const (
	Wat P0Target = iota
	Mips
)

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

package main

import (
	"fmt"
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
var sourceUnitChannel chan SourceUnit
var currSym int
var currVal interface{}

// Gets the next symbol from the source unit channel
// Sets currVal and codePosition
func getNextSym() int {
	su := <-sourceUnitChannel
	currSym = su.sym
	currVal = su.val
	//fmt.Println(currSym)
	return currSym
}

func selector(x Entry) Entry {
	var a = [2]int{PERIOD, LBRAK}
	for doesContain(a[:], currSym) {
		if currSym == PERIOD { // x.f
			getNextSym()
			if currSym == IDENT {
				asRec, isRec := x.GetP0Type().(*P0Record)
				if isRec {
					fieldFound := false
					for _, f := range asRec.GetFields() {
						if f.GetName() == currVal.(string) {
							x = cg.GenSelect(x, f)
							fieldFound = true
							break
						}
					}
					if !fieldFound {
						mark("not a field")
					}
					getNextSym()
				} else {
					mark("not a record")
				}
			} else {
				mark("identifier expected")
			}
		} else { // x[y]
			getNextSym()
			var y = expression()
			xAsArray, xIsArray := x.GetP0Type().(*P0Array)
			if xIsArray {
				_, yIsInt := y.GetP0Type().(*P0Int)
				if yIsInt {
					var lowerbound = xAsArray.GetLowerBound()
					yAsConst, castSucceed := y.(*P0Const)
					if castSucceed && (yAsConst.GetValue().(int) < lowerbound || yAsConst.GetValue().(int) >= lowerbound+xAsArray.GetLength()) {
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
			if currSym == RBRAK {
				getNextSym()
			} else {
				mark("] expected")
			}
		}
	}
	return x
}

func factor() Entry {
	if !doesContain(FIRSTFACTOR[:], currSym) {
		mark("expression expected")
		for !(doesContain(FOLLOWFACTOR[:], currSym) || doesContain(STRONGSYMS[:], currSym) ||
			doesContain(FIRSTFACTOR[:], currSym)) {
			getNextSym()
		}
	}
	var x Entry
	if currSym == IDENT {
		x = st.Find(currVal.(string))

		_, xIsVar := x.(*P0Var)
		_, xIsRef := x.(*P0Ref)
		xAsConst, xIsConst := x.(*P0Const)

		if xIsVar || xIsRef {
			x = cg.GenVar(x)
			getNextSym()
		} else if xIsConst {
			x = &P0Const{
				x.GetP0Type(),
				x.GetName(),
				x.GetLevel(),
				xAsConst.GetValue().(int),
			}
			x = cg.GenConst(x)
			getNextSym()
		} else {
			mark("expression expected")
		}
		x = selector(x)
	} else if currSym == NUMBER {
		x = &P0Const{
			cg.GenInt(&P0Int{}),
			"",
			0,
			currVal,
		}
		x = cg.GenConst(x)
		getNextSym()
	} else if currSym == LPAREN {
		getNextSym()
		x = expression()
		if currSym == RPAREN {
			getNextSym()
		} else {
			mark(") expected")
		}
	} else if currSym == NOT {
		getNextSym()
		x = factor()
		_, xIsBool := x.GetP0Type().(*P0Bool)
		xAsConst2, xIsConst2 := x.(*P0Const)
		if !xIsBool {
			mark("not boolean")
		} else if xIsConst2 {
			// Toggle x
			xAsConst2.SetValue(1 - xAsConst2.GetValue().(int))
		} else {
			x = cg.GenUnaryOp(NOT, x)
		}
	} else {
		x = &P0Const{
			nil, // TODO: Check if nil is supposed to go here
			x.GetName(),
			x.GetLevel(),
			0,
		}
	}
	return x
}

func term() Entry {
	var x = factor()

	for currSym == TIMES || currSym == DIV || currSym == MOD || currSym == AND {
		var op = currSym
		getNextSym()

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
	if currSym == PLUS {
		getNextSym()
		x = term()
	} else if currSym == MINUS {
		getNextSym()
		x = term()
		_, xIsInt := x.GetP0Type().(*P0Int)
		if !xIsInt {
			mark("Bad type")
		}
	} else {
		x = term()
	}

	for currSym == PLUS || currSym == MINUS || currSym == OR {
		var op = currSym
		getNextSym()

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
	for currSym == EQ || currSym == NE || currSym == LT || currSym == LE || currSym == GT || currSym == GE {
		op := currSym
		getNextSym()
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
				case NE:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) != yAsConst.GetValue().(int)))
				case LT:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) < yAsConst.GetValue().(int)))
				case LE:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) <= yAsConst.GetValue().(int)))
				case GT:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) > yAsConst.GetValue().(int)))
				case GE:
					xAsConst.SetValue(bool2int(xAsConst.GetValue().(int) >= yAsConst.GetValue().(int)))
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
	getElseMark(currSym == BEGIN, "'begin' expected")
	x := statement()
	for currSym == SEMICOLON || doesContain(FIRSTSTATEMENT[:], currSym) {
		getElseMark(currSym == SEMICOLON, "; missing")
		y := statement()
		cg.GenSeq(x, y) // This returns a value in p0 to build the AST string representation; we just set x to nil here (same effect)
		x = nil
	}
	getElseMark(currSym == END, "'end' expected")
	return x
}

// statement does a few things. It's a pretty cool function, you should take a look at the source code.
func statement() Entry {
	var x Entry
	if !doesContain(FIRSTSTATEMENT[:], currSym) {
		mark("statement expected")
		for !doesContain(FIRSTSTATEMENT[:], currSym) && !doesContain(FOLLOWSTATEMENT[:], currSym) && !doesContain(STRONGSYMS[:], currSym) {
			getNextSym()
		}
	}
	if currSym == IDENT {
		x = st.Find(currVal.(string))
		getNextSym()
		switch x.(type) {
		case *P0Var, *P0Ref:
			x = cg.GenVar(x)
			x = selector(x)
			if currSym == BECOMES {
				getNextSym()
				y := expression()
				_, xIsBool := x.GetP0Type().(*P0Bool)
				_, yIsBool := y.GetP0Type().(*P0Bool)
				_, xIsInt := x.GetP0Type().(*P0Int)
				_, yIsInt := y.GetP0Type().(*P0Int)
				if (xIsBool && yIsBool) || (xIsInt && yIsInt) {
					cg.GenAssign(x, y)
					x = nil // FIXME: ?
				} else {
					mark("incompatible assignment")
				}
			} else if currSym == EQ {
				mark(":= expected")
				getNextSym()
				_ = expression() // We parse to consume the input, but we can't use the result because the code in incorrect
			} else {
				mark(":= expected")
			}
		case *P0Proc, *P0StdProc:
			// This man codes 8 lines of Go in one line of python
			var fp []P0Type
			var y Entry
			xAsProc, xIsProc := x.(*P0Proc)
			if xIsProc {
				fp = xAsProc.GetParameters()
			} else {
				fp = x.(*P0StdProc).GetParameters()
			}
			i := 0
			if currSym == LPAREN {
				getNextSym()
				if doesContain(FIRSTEXPRESSION[:], currSym) {
					y = expression()
					if i < len(fp) {
						if typesEqual(fp[i].GetP0Type(), y.GetP0Type()) { // TODO: How to do this properly in Go?
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
					for currSym == COMMA {
						getNextSym()
						y = expression()
						if i < len(fp) {
							if typesEqual(fp[i].GetP0Type(), y.GetP0Type()) { // TODO: How to do this properly in Go?
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
				}
				getElseMark(currSym == RPAREN, "')' expected")
				if i < len(fp) {
					mark("too few parameters")
				} else if !xIsProc { // x is P0StdProc
					if x.GetName() == "read" {
						cg.GenRead(y)
					} else if x.GetName() == "write" {
						cg.GenWrite(y)
					} else if x.GetName() == "writeln" {
						cg.GenWriteln()
					}
				} else {
					cg.GenCall(x)
				}
			} else {
				mark("'(' expected")
			}
		default:
			mark("variable or procedure expected")
		}
	} else if currSym == BEGIN {
		x = compoundStatement()
	} else if currSym == IF {
		getNextSym()
		x = expression()
		_, xIsBool := x.GetP0Type().(*P0Bool)
		if xIsBool {
			x = cg.GenThen(x)
		} else {
			mark("boolean expected")
		}
		getElseMark(currSym == THEN, "'then' expected")
		y := statement()
		if currSym == ELSE {
			_, xIsBool = x.GetP0Type().(*P0Bool)
			var label string
			if xIsBool {
				label = cg.GenElse(x, y) // TODO: GenElse needs to return something
			}
			getNextSym()
			statement()
			_, xIsBool = x.GetP0Type().(*P0Bool)
			if xIsBool {
				cg.GenIfElse(label)
				x = nil
			}
		} else {
			_, xIsBool = x.GetP0Type().(*P0Bool)
			cg.GenIfThen(x) // TODO: GenIfThen needs to return something
			x = nil
		}
	} else if currSym == WHILE {
		getNextSym()
		t := cg.GenWhile()
		x = expression()
		_, xIsBool := x.GetP0Type().(*P0Bool)
		if xIsBool {
			x = cg.GenDo(x)
		} else {
			mark("boolean expected")
		}
		getElseMark(currSym == DO, "'do' expected")
		y := statement()
		_, xIsBool = x.GetP0Type().(*P0Bool)
		if xIsBool {
			cg.GenWhileDo(t, x, y)
			x = nil
		}
	} else {
		x = nil
	}
	return x
}

func typ() P0Type {
	var typeToReturn P0Type
	if !doesContain(FIRSTTYPE[:], currSym) {
		mark("type expected")
		for !(doesContain(FIRSTTYPE[:], currSym) || doesContain(FOLLOWTYPE[:], currSym) || doesContain(STRONGSYMS[:], currSym)) {
			getNextSym()
		}
	}
	if currSym == IDENT {
		ident := currVal.(string)
		typeToReturn, _ = st.Find(ident).(P0Type)
		getNextSym()
	} else if currSym == ARRAY {
		getNextSym()
		getElseMark(currSym == LBRAK, "'[' expected")
		x := expression()
		getElseMark(currSym == PERIOD, "'.' expected")
		getElseMark(currSym == PERIOD, "'.' expected")
		y := expression()
		getElseMark(currSym == RBRAK, "']' expected")
		getElseMark(currSym == OF, "'of' expected")
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
	} else if currSym == RECORD {
		getNextSym()
		st.OpenScope()
		typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type: p0type} })
		for currSym == SEMICOLON {
			getNextSym()
			typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type: p0type} })
		}
		getElseMark(currSym == END, "'end' expected")
		r := st.TopScope()
		st.CloseScope()
		typeToReturn = cg.GenRecord(&P0Record{fields: r})
	} else {
		typeToReturn = nil
	}
	return typeToReturn
}

func typedIds(kind func(P0Type) P0Type) {
	var tid []string
	if currSym == IDENT {
		tid = make([]string, 1)
		tid[0] = currVal.(string)
		getNextSym()
	} else {
		mark("identifier expected")
		tid = make([]string, 0)
	}
	for currSym == COMMA {
		getNextSym()
		if currSym == IDENT {
			tid = append(tid, currVal.(string))
			getNextSym()
		} else {
			mark("identifier expected")
		}
	}
	if currSym == COLON {
		getNextSym()
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
	if !(doesContain(FIRSTDECL[:], currSym) || doesContain(FOLLOWDECL[:], currSym)) {
		mark("'begin' or declaration expected")
		for !(doesContain(FIRSTDECL[:], currSym) || doesContain(FOLLOWDECL[:], currSym) || doesContain(STRONGSYMS[:], currSym)) {
			getNextSym()
		}
	}
	for currSym == CONST {
		getNextSym()
		if currSym == IDENT {
			ident := currVal.(string)
			getNextSym()
			getElseMark(currSym == EQ, "= expected")
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
		getElseMark(currSym == SEMICOLON, "; expected")
	}
	for currSym == TYPE {
		getNextSym()
		if currSym == IDENT {
			ident := currVal.(string)
			getNextSym()
			getElseMark(currSym == EQ, "= expected")
			x := typ()
			st.NewDecl(ident, x)
			getElseMark(currSym == SEMICOLON, "; expected")
		} else {
			mark("type name expected")
		}
	}
	st.OpenScope()
	for currSym == VAR {
		getNextSym()
		typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type, "", 0, "", 0, 0} })
		getElseMark(currSym == SEMICOLON, "; expected")
	}
	localVarDecls := st.TopScope()
	st.CloseScope()
	for i := 0; i < len(localVarDecls); i++ {
		localVarDecls[i].SetLevel(localVarDecls[i].GetLevel() - 1)
		st.NewDecl(localVarDecls[i].GetName(), localVarDecls[i])
	}
	varsize = generatorFunc(localVarDecls, 0)
	for currSym == PROCEDURE {
		getNextSym()
		getElseMark(currSym == IDENT, "procedure name expected")
		ident := currVal.(string)
		st.NewDecl(ident, &P0Proc{nil, "", 0, nil})
		st.OpenScope()
		var fp []Entry
		if currSym == LPAREN {
			getNextSym()
			if currSym == VAR || currSym == IDENT {
				if currSym == VAR {
					getNextSym()
					typedIds(func(p0type P0Type) P0Type { return &P0Ref{p0type, "", 0, "", 0, 0} })
				} else {
					typedIds(func(p0type P0Type) P0Type { return &P0Var{p0type, "", 0, "", 0, 0} })
				}
				for currSym == SEMICOLON {
					getNextSym()
					if currSym == VAR {
						getNextSym()
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
			tmp := make([]P0Type, len(fp))
			for i, item := range fp {
				tmp[i] = item.(P0Type)
			}
			st.Find(ident).(*P0Proc).parameters = tmp
			getElseMark(currSym == RPAREN, ") expected")
		} else {
			fp = make([]Entry, 0)
		}
		parsize := cg.GenProcStart(ident, fp)
		getElseMark(currSym == SEMICOLON, "; expected")
		localsize := declarations(cg.GenLocalVars)
		cg.GenProcEntry(ident, parsize, localsize)
		var x Entry = compoundStatement()
		cg.GenProcExit(x, parsize, localsize)
		st.CloseScope()
		getElseMark(currSym == SEMICOLON, "; expected")
	}
	return varsize
}

func program() string {
	st.NewDecl("boolean", cg.GenBool(&P0Bool{}))
	st.NewDecl("integer", cg.GenInt(&P0Int{}))
	st.NewDecl("true", &P0Const{cg.GenBool(&P0Bool{}), "", 0, 1})
	st.NewDecl("false", &P0Const{cg.GenBool(&P0Bool{}), "", 0, 0})
	st.NewDecl("read", &P0StdProc{nil, "", 0, []P0Type{&P0Ref{cg.GenInt(&P0Int{}), "", 0, "", 0, 0}}})
	st.NewDecl("write", &P0StdProc{nil, "", 0, []P0Type{&P0Var{cg.GenInt(&P0Int{}), "", 0, "", 0, 0}}})
	st.NewDecl("writeln", &P0StdProc{nil, "", 0, []P0Type{}})
	cg.GenProgStart()
	getNextSym()
	getElseMark(currSym == PROGRAM, "'program expected")
	// The original program actually accessed the program name here
	getElseMark(currSym == IDENT, "Program name expected")
	getElseMark(currSym == SEMICOLON, "; expected")
	declarations(cg.GenGlobalVars)
	cg.GenProgEntry( /*ident*/ ) // ident was passed in the og P0 compiler, but it is not used so we removed it
	compoundStatement()
	return cg.GenProgExit()
}

func compileString(destinationFilePath string, target P0Target) {
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
	st = new(SliceMapSymbolTable)
	st.Init()
	p := program()
	if p != "" && !error {
	} else {
		panic("Something went wrong in the parser and its not good :(")
	}
}

func compileFile(tokenChannel chan SourceUnit, endChannel chan int, destFilePath string, target string) {
	sourceUnitChannel = tokenChannel
	compileString(destFilePath, toP0Target(target))
	endChannel <- 0
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

// If the predicate is true, getNextSym is called. Otherwise, mark is called with the message
// I introduced to try to make the code easier to read; if it has the opposite effect let me know - David
func getElseMark(predicate bool, markMessage string) {
	if predicate {
		getNextSym()
	} else {
		mark(markMessage)
	}
}

// typesEqual checks if two P0 types are equal. It uses duck typing I guess ;-;
func typesEqual(a, b P0Type) bool {
	switch t := a.(type) {
	case *P0Int:
		_, bIsInt := b.(*P0Int)
		return bIsInt
	case *P0Bool:
		_, bIsBool := b.(*P0Bool)
		return bIsBool
	case *P0Record:
		aAsRec := a.(*P0Record)
		bAsRec, bIsRec := b.(*P0Record)
		if !bIsRec || len(aAsRec.GetFields()) != len(bAsRec.GetFields()) {
			return false
		}
		for i := 0; i < len(aAsRec.GetFields()); i++ {
			if !typesEqual(aAsRec.GetFields()[i].GetP0Type(), bAsRec.GetFields()[i].GetP0Type()) {
				return false
			}
		}
		return true
	case *P0Array:
		aAsArray := a.(*P0Array)
		bAsArray, bIsArray := b.(*P0Array)
		return bIsArray && typesEqual(aAsArray.GetElementType(), bAsArray.GetElementType())
	default:
		panic(fmt.Sprint("Unrecognized type %v", t))
	}
}

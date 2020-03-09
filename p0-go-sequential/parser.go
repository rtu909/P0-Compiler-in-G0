package main

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
				if x.GetP0Type() == Record {
					for f := range x.GetFieldNames() {
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
			if x.GetP0Type() == Array {
				if y.GetP0Type() == Int {
					var value = y.GetValue()
					var lowerbound = x.GetLowerBound()
					if y.IsConstant() && value < lowerbound || y.GetValue() >= lowerbound+x.GetLength() {
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

// TODO: IMPLEMENT
func expression() Entry {
	return nil
}

func doesContain(elements []int, e int) bool {
	for _, a := range elements {
		if a == e {
			return true
		}
	}
	return false
}

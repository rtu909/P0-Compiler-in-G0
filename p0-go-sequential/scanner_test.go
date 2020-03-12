package main

import "testing"

func TestTimes(t *testing.T) {
	ScannerInit("*")
	if sym != 1 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestDiv(t *testing.T) {
	ScannerInit("div")
	if sym != 2 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestMod(t *testing.T) {
	ScannerInit("mod")
	if sym != 3 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestAnd(t *testing.T) {
	ScannerInit("and")
	if sym != 4 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestPlus(t *testing.T) {
	ScannerInit("+")
	if sym != 5 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestMinus(t *testing.T) {
	ScannerInit("-")
	if sym != 6 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestOr(t *testing.T) {
	ScannerInit("or")
	if sym != 7 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestEq(t *testing.T) {
	ScannerInit("=")
	if sym != 8 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestNe(t *testing.T) {
	ScannerInit("<>")
	if sym != 9 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestLT(t *testing.T) {
	ScannerInit("<")
	if sym != 10 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestGT(t *testing.T) {
	ScannerInit(">")
	if sym != 11 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestLE(t *testing.T) {
	ScannerInit("<=")
	if sym != 12 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestGE(t *testing.T) {
	ScannerInit(">=")
	if sym != 13 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestPeriod(t *testing.T) {
	ScannerInit(".")
	if sym != 14 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestComma(t *testing.T) {
	ScannerInit(",")
	if sym != 15 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestColon(t *testing.T) {
	ScannerInit(":")
	if sym != 16 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestRParen(t *testing.T) {
	ScannerInit(")")
	if sym != 17 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestRBRak(t *testing.T) {
	ScannerInit("]")
	if sym != 18 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestOf(t *testing.T) {
	ScannerInit("of")
	if sym != 19 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestThen(t *testing.T) {
	ScannerInit("then")
	if sym != 20 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestDo(t *testing.T) {
	ScannerInit("do")
	if sym != 21 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestLParen(t *testing.T) {
	ScannerInit("(")
	if sym != 22 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestLBrak(t *testing.T) {
	ScannerInit("[")
	if sym != 23 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestNot(t *testing.T) {
	ScannerInit("not")
	if sym != 24 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestBecomes(t *testing.T) {
	ScannerInit(":=")
	if sym != 25 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestNumber(t *testing.T) {
	ScannerInit("23456")
	if sym != 26 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestIdent(t *testing.T) {
	ScannerInit("potato")
	if sym != 27 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestSemicolon(t *testing.T) {
	ScannerInit(";")
	if sym != 28 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestEnd(t *testing.T) {
	ScannerInit("end")
	if sym != 29 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestElse(t *testing.T) {
	ScannerInit("else")
	if sym != 30 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestIf(t *testing.T) {
	ScannerInit("if")
	if sym != 31 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestWhile(t *testing.T) {
	ScannerInit("while")
	if sym != 32 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestArray(t *testing.T) {
	ScannerInit("array")
	if sym != 33 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestRecord(t *testing.T) {
	ScannerInit("record")
	if sym != 34 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestConst(t *testing.T) {
	ScannerInit("const")
	if sym != 35 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestType(t *testing.T) {
	ScannerInit("type")
	if sym != 36 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestVar(t *testing.T) {
	ScannerInit("var")
	if sym != 37 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestProcedure(t *testing.T) {
	ScannerInit("procedure")
	if sym != 38 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestBegin(t *testing.T) {
	ScannerInit("begin")
	if sym != 39 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestProgram(t *testing.T) {
	ScannerInit("program")
	if sym != 40 {
		t.Errorf("The symbol found is incorrect")
	}
}

func TestEOF(t *testing.T) {
	ScannerInit("00")
	getSym()
	if sym != 41 {
		t.Errorf("The symbol found is incorrect")
	}
}
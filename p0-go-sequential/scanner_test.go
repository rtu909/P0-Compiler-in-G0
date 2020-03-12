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
	ScannerInit("}")
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
	ScannerInit(")")
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
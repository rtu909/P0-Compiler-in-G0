package main

import "testing"

func TestInit(t *testing.T) {
	// a) Not necessarily, but often it is
	// b) If the requirement disallows zero, but zero is actually a valid input that is easy to account for
	// c) No process could always prevent late detection of faults, but some are better than others at detecting them earlier
	var ast ArraySymbolTable
	ast.Init()
	ast.NewDecl("potate" /*TODO*/)
}

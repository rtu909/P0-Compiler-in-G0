package main

import "testing"

func TestVariableFind(t *testing.T) {
	var st ArraySymbolTable
	var entry P0Var
	st.Init()
	st.NewDecl("potate", entry)
	found := st.Find("potate")
	if found.GetP0Type() != Int {
		t.Errorf("The type of variable potate was incorrect")
	}
}

func TestConstantFind(t *testing.T) {
	var st ArraySymbolTable
	myConst := P0Const{Int, 67}
	st.Init()
	st.NewDecl("i", myConst)
	found := st.Find("i")
	castFound, castOk := found.(P0Const)
	if !castOk {
		t.Errorf("The declaration couldn't be retrieved as a const")
	}
	if castFound.p0type != Int || castFound.value != 67 {
		t.Errorf("The type or value of the const was incorrect")
	}
}

func TestEmptyFind(t *testing.T) {
	var st ArraySymbolTable
	st.Init()
	found := st.Find("item")
	if found.GetP0Type() != None {
		t.Errorf("The type was supposed to be nil, but was %v", found.GetP0Type())
	}
	foundConst, castOk := found.(P0Const)
	if !castOk {
		t.Errorf("A P0Const was expected, but not returned")
	}
	if foundConst.p0type != None || foundConst.value != 0 {
		t.Errorf("The value ")
	}
}

package main

import "testing"

func TestVariableFind(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myVar := &P0Var{&P0Bool{}, "", 0, "", 0, 0}
	st.Init()
	st.NewDecl("potato", myVar)
	foundValue := st.Find("potato")
	asVar, isVar := foundValue.(*P0Var)
	if !isVar {
		t.Errorf("The found declaration was not a variable")
	}
	_, isBool := asVar.GetP0Type().(*P0Bool)
	if !isBool {
		t.Errorf("The found declaration was not of type bool")
	}
}

func TestConstantFind(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst := &P0Const{&P0Int{}, "", 0, 67}
	st.Init()
	st.NewDecl("i", myConst)
	found := st.Find("i")
	asConst, isConst := found.(*P0Const)
	if !isConst {
		t.Errorf("The declaration couldn't be retrieved as a const")
	}
	_, isInt := asConst.GetP0Type().(*P0Int)
	if !isInt {
		t.Errorf("The declaration is not an integer")
	}
	valAsInt, valIsInt := asConst.GetValue().(int)
	if !valIsInt {
		t.Errorf("Unable to cast the constant value to an integer")
	}
	if valAsInt != 67 {
		t.Errorf("Incorrect value of the constant")
	}
}

func TestEmptyFind(t *testing.T) {
	st := new(SliceMapSymbolTable)
	st.Init()
	found := st.Find("item")
	asConst, isConst := found.(*P0Const)
	if !isConst {
		t.Error("Undeclared value is not given as a const")
	}
	valAsInt, valIsInt := asConst.GetValue().(int)
	if !valIsInt {
		t.Errorf("Unable to cast the constant value to an integer")
	}
	if valAsInt != 0 {
		t.Errorf("Incorrect value of the constant")
	}
}

func TestDeclDroppedOutOfScope(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myVar := &P0Var{&P0Int{}, "", 0, "", 0, 0}
	st.Init()
	st.OpenScope()
	st.NewDecl("potato", myVar)
	st.CloseScope()
	found := st.Find("potato")
	asConst, isConst := found.(*P0Const)
	if !isConst {
		t.Error("A declaration was found, but it was supposed to drop out of scope")
	}
	valAsInt, valIsInt := asConst.GetValue().(int)
	if !valIsInt {
		t.Error("A declaration was found, but it was supposed to drop out of scope")
	}
	if valAsInt != 0 {
		t.Error("A declaration was found, but it was supposed to drop out of scope")
	}
}

func TestOuterDeclarationFound(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst1 := &P0Const{&P0Int{}, "", 0, 42}
	myConst2 := &P0Const{&P0Int{}, "", 0, 68}
	st.Init()
	st.NewDecl("potato", myConst1)
	st.OpenScope()
	st.NewDecl("potato", myConst2)
	found := st.Find("potato").(*P0Const)
	if found.GetValue().(int) != 68 {
		t.Errorf("Picked up a value of %v, expect 68", found.value)
	}
	st.CloseScope()
	found = st.Find("potato").(*P0Const)
	if found.GetValue() != 42 {
		t.Errorf("Picked up a value of %v, expected 42", found.GetValue())
	}
}

func TestFindDeclarationFromInnerLabel(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst := &P0Const{&P0Int{}, "", 0, 42}
	st.Init()
	st.NewDecl("potato", myConst)
	st.OpenScope()
	found, castOk := st.Find("potato").(*P0Const)
	if !castOk {
		t.Error("Found something of an unexpected type")
	}
	if found.GetValue() != 42 {
		t.Errorf("Found a value of %v, but expected 42", myConst.GetValue())
	}
}

func TestTopScope(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst := &P0Const{&P0Int{}, "", 0, 42}
	st.Init()
	st.NewDecl("potato", myConst)
	st.OpenScope()
	myOtherConst := &P0Const{&P0Int{}, "", 0, 32}
	st.NewDecl("cilantro", myOtherConst)
	for _, val := range st.TopScope() {
		asConst, isConst := val.(*P0Const)
		if !isConst {
			t.Error("Only consts were declared, but a non-const declaration was found")
		}
		if asConst.GetValue() != 32 || asConst.GetName() != "cilantro" {
			t.Error("Found something that wasn't supposed to be in the top declaration")
		}
	}
	st.CloseScope()
	for _, val := range st.TopScope() {
		asConst, isConst := val.(*P0Const)
		if !isConst {
			t.Error("Only consts were declared, but a non-const declaration was found")
		}
		if asConst.GetValue() != 42 || asConst.GetName() != "potato" {
			t.Error("Found something that wasn't supposed to be in the top declaration")
		}
	}
}

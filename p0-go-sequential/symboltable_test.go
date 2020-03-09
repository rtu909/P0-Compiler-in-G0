package main

import "testing"

func TestVariableFind(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myVar := P0Var{Int, nil}
	st.Init()
	st.NewDecl("potato", myVar)
	found := st.Find("potato")
	if found.GetP0Type().p0primitive != Int {
		t.Errorf("The type of variable potate was incorrect")
	}
}

func TestConstantFind(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst := P0Const{P0Type{Int, nil}, 67}
	st.Init()
	st.NewDecl("i", myConst)
	found := st.Find("i")
	castFound, castOk := found.(P0Const)
	if !castOk {
		t.Errorf("The declaration couldn't be retrieved as a const")
	}
	if castFound.p0type.p0primitive != Int || castFound.value != 67 {
		t.Errorf("The type or value of the const was incorrect")
	}
}

func TestEmptyFind(t *testing.T) {
	st := new(SliceMapSymbolTable)
	st.Init()
	found := st.Find("item")
	if found.GetP0Type().p0primitive != None {
		t.Errorf("The type was supposed to be nil, but was %v", found.GetP0Type().p0primitive)
	}
	foundConst, castOk := found.(P0Const)
	if !castOk {
		t.Errorf("A P0Const was expected, but not returned")
	}
	if foundConst.p0type.p0primitive != None || foundConst.value != 0 {
		t.Errorf("The value was %v, but was expected to be 0", found.GetP0Type().p0primitive)
	}
}

func TestDeclDroppedOutOfScope(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst := P0Var{Int, nil}
	st.Init()
	st.OpenScope()
	st.NewDecl("potato", myConst)
	st.CloseScope()
	found := st.Find("potato")
	castFound, castOk := found.(P0Const)
	if !castOk || castFound.GetP0Type().p0primitive != None || castFound.value != 0 {
		t.Errorf("A declaration of type %v, value %v was found, but there should be no variables in scope", castFound.GetP0Type().p0primitive, castFound.value)
	}
}

func TestOuterDeclarationFound(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst1 := P0Const{P0Type{Int, nil}, 42}
	myConst2 := P0Const{P0Type{Int, nil}, 68}
	st.Init()
	st.NewDecl("potato", myConst1)
	st.OpenScope()
	st.NewDecl("potato", myConst2)
	found := st.Find("potato").(P0Const)
	if found.value != 68 {
		t.Errorf("Picked up a value of %v, expect 68", found.value)
	}
	st.CloseScope()
	found = st.Find("potato").(P0Const)
	if found.value != 42 {
		t.Errorf("Picked up a value of %v, expected 42", found.value)
	}
}

func TestFindDeclarationFromInnerLabel(t *testing.T) {
	st := new(SliceMapSymbolTable)
	myConst := P0Const{P0Type{Int, nil}, 42}
	st.Init()
	st.NewDecl("potato", myConst)
	st.OpenScope()
	found, castOk := st.Find("potato").(P0Const)
	if !castOk {
		t.Error("Found something of an unexpected type")
	}
	if found.value != 42 {
		t.Errorf("Found a value of %v, but expected 42", myConst)
	}
}

package main

// This is the generic interface for a code generator

type CodeGenerator interface {
	GenProgStart()                                     //
	GenBool(P0Type) P0Type                             //
	GenInt(P0Type) P0Type                              //
	GenRecord(P0Type) P0Type                           //
	GenArray(P0Type) P0Type                            //
	GenGlobalVars(declaredVars []Entry, start int) int //no return statement
	GenLocalVars(declaredVars []Entry, start int) int  // doesn't exist
	GenVar(Entry) Entry                                //
	GenConst(Entry) Entry                              //
	GenUnaryOp(op int, entry Entry) Entry              //issues with entry
	GenBinaryOp(op int, x, y Entry) Entry
	GenRelation(op int, x, y Entry) Entry
	GenSelect(record Entry, field Entry) Entry
	GenIndex(array Entry, index Entry) Entry
	GenAssign(x, y Entry)
	GenProgEntry()       //
	GenProgExit() string //
	GenProcStart(identity string, functionParameters []Entry) int
	GenProcEntry(identity string, parsize, localsize int)
	GenProcExit(x Entry, parsize, localsize int)
	GenActualPara(actualparameter, formalparameter Entry, parameterNumber int)
	GenCall(procedure Entry)
	GenRead(Entry)
	GenWrite(Entry)
	GenWriteln()
	GenSeq(x, y Entry)
	GenThen(Entry) Entry
	GenIfThen(Entry)
	GenElse(x, y Entry) string
	GenIfElse(label string)
	GenWhile() string
	GenDo(Entry) Entry
	GenWhileDo(t string, y, z Entry)
}

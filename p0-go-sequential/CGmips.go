package main

import (
	"container/list"
	"go/types"
	"strconv"
)

//triple tuple data structure
type Triple struct {
	a, b, c interface{}
}

//global variables
var curlev int
var label int
var GPRegs = []string{"$t0", "$t1", "$t2", "$t3", "$t4", "$t5", "$t6", "$t7", "$t8" }
var regs [] string
var asm [] Triple

//reserved registers
var R0 = "$0"
var FP = "$fp"
var SP = "$sp"
var LNK = "$a3"
var A0 = "$a0"
var A1 = "$a1"
var A2 = "$a2"
var A3 = "$a3"

func genProgStart(){
	curlev, label = 0, 0
	regs = GPRegs
	putInstr(".data", "")
}

func obtainReg() string{
	if len(regs) == 0{
		mark("out of registers");
		return R0
	} else {
		var popped = regs[9]
		regs = regs[0:len(regs)-1]
		return popped
	}
}

func releaseReg(r string){
	for i := 0; i < len(GPRegs); i++{
		if r == GPRegs[i]{
			regs = append(regs, r)
		}
	}
}

func putLab(lab []string, instr string){

	if len(lab) == 1{
		tuple := Triple{lab[0], instr, "" }
		asm = append(asm, tuple)
	} else {
		for i := 0; i < len(lab)-1; i++{
			tuple := Triple{lab[i], "", "" }
			asm = append(asm, tuple)
		}
		tuple := Triple{lab[len(lab)-1], instr, "" }
		asm = append(asm, tuple)
	}
}

func putInstr(instr string, target string){
	tuple := Triple{"", instr, target}
	asm = append(asm, tuple)
}

func putOp(op string, a string, b string, c string){
	putInstr(op + " " + a + ", " + b + ", " + c, "")
}

func putBranchOp(op string, a string, b string, c string){
	putInstr(op + " " + a + ", " + b, c)
}

func putMemOp(op string, a string, b string, c string){
	if b == R0{
		putInstr(op + " " + a + ", " + c, "")
	} else {
		putInstr(op + " " + a + ", " + c + "(" + b + ")", "")
	}
}

//size - not sure what's going on here in the regular code
func genBool(b P0Bool) P0Bool{
	b.SetSize(4)
	return b
}

func genInt(i P0Int) P0Int{
	i.SetSize(4)
	return i
}

func genRec(r P0Record) P0Record{
	s:= 0
	fields := r.GetFields()
	for f := 0; f < len(fields); f++{

	}
	r.SetSize(s)
	return r
}

func genArray(a P0Array) P0Array {
	size := a.GetLength() + a.GetElementType().GetSize()
	a.SetSize(size)
	return a
}

//todo
func genGlobalVars(sc []P0Var, start int ){
	for i:= len(sc) -1; i > start -1; i--{

	}
	putInstr(".text", "")
}

func genProgEntry(){
	putInstr(".globl main", "")
	putInstr(".ent main", "")
	lab := []string
	putLab(lab, "main" )
}

func assembly(l string, i string, t string) string{
	string1 := ""
	if l != ""{
		string1 = l + ":\t"
	} else {
		string1 = "\t"
	}
	string2 := ""
	if t != ""{
		string2 = ", " + t
	} else {
		string2 = ""
	}
	string3 := string1 + i + string2
	return string3
}

func genProgExit() string{
	putInstr("li $v0, 10", "")
	putInstr("syscall", "")
	putInstr(".end main", "")
	returnStr := ""
	for i := 0; i < len(asm); i++{
		asm_l := asm[i].a
		asm_i := asm[i].b
		asm_t := asm[i].c
		returnStr = returnStr + assembly(asm_l.(string), asm_i.(string), asm_t.(string)) + "\n"
	}
	return (returnStr)
}

func newLabel() string{
	label = label + 1
	return ("L" + strconv.Itoa(label))
}

type Reg struct {
	tp interface{}
	reg string
}

type Cond struct{
	tp interface{}
	cond string
	left, right interface{}
	labA []string
	labB []string
}

func NewCond(tp interface{}, cond string, left interface{}, right interface{}) Cond{
	var labA []string
	var labB []string
	labA = append(labA, newLabel())
	labB = append(labB, newLabel())
	c := Cond{
		tp:    tp,
		cond:  cond,
		left:  left,
		right: right,
		labA:  labA,
		labB:  labB,
	}
	return (c)
}

func testRange(x P0Type){
	if (x.GetLevel() >= 0x8000) || (x.GetLevel() < -0x8000){
		mark("value too large")
	}
}

//todo
func loadItemReg(x P0Type, r string){
	if {
		putMemOp("lw", r, x.GetName(), r)
		releaseReg(x.GetName())
	}else if {
		testRange(x)
		putOp("addi", r, R0, strconv.Itoa(x.GetLevel()))
	}else if{
		putOp("add", r, x.GetName(), R0)
	} else{
		panic("loadItemReg has problems")
	}
}

//todo
func loadItem(x P0Type) Reg{
	if {
		r := R0
	} else{
		r := obtainReg()
		loadItemReg(x, r)
	}
	//Reg(x.GetP0Type(), r)
}

//todo
func loadBool(x P0Type) Cond{
	if {
		r := R0
	} else{
		r := obtainReg()
		loadItemReg(x, r)
	}
}

func put(){

}

func genVar(){

}

func genConst(){

}

func negate(){

}

func condOp(){

}

func genUnaryOp(){

}

func genBinaryOp(){

}

func genRelation(){

}

func genSelect(){

}

func genIndex(){

}

func genAssign(){

}

func genLocalVars(){

}

func genProcStart(){

}

func genProcEntry(){

}

func genProcExit(){

}

func genActualPara(){

}

func genCall(){

}

func genRead(){

}

func genWrite(){

}

func genWriteln(){

}

func genSeq(){

}

func genThen(){

}

func genIfThen(){

}

func genElse(){

}

func genIfElse(){

}

func genWhile(){
	lab := newLabel()
	putLab(lab, "")
}

func genDo(){
	genThen()
}

func genWhileDo(lab string, x Cond){
	putInstr("b", lab)
	putLab(x.labA, "")
}
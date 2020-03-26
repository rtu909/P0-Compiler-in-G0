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

func NewReg(tp interface{}, reg string) Reg{
	r := Reg{
		tp:  tp,
		reg: reg,
	}
	return r
}

type Cond struct{
	tp interface{}
	cond string
	left, right interface{}
	labA []string
	labB []string
}

func (condition *Cond) SetCondition(newCond string){
	(*condition).cond = newCond
}

func (condition *Cond) SetLabA(newCond[]string){
	(*condition).labA = newCond
}

func (condition *Cond) SetLabB(newCond[]string){
	(*condition).labB = newCond
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

func testRange(x P0Const){
	if (x.GetValue().(int) >= 0x8000) || (x.GetValue().(int) < -0x8000){
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

func put(cd string, x interface{}, y interface{}) interface{} {
	_, xisReg := x.(*Reg)
	r := ""
	if !xisReg{
		x = loadItem(x.(P0Type))
	} else {
		var regList []string
		regList = append(regList, R0)
		regList = append(regList, A0)
		regList = append(regList, A1)
		regList = append(regList, A2)
		regList = append(regList, A3)
		x := x.(Reg)
		var regFound bool
		regFound = false

		for i := 0; i < len(regList); i++ {
			if regList[i] == x.reg {
				regFound = true
			}
		}
		if regFound {
			r = x.reg
			x.reg = obtainReg()
		} else {
			r = x.reg
		}
	}
		_, yisConst := y.(*P0Const)
		if yisConst{
			y := y.(P0Const)
			testRange(y)
			putOp(cd, x.(Reg).reg, r, y.GetValue().(string))
		} else{
			_, yisReg := y.(*Reg)
			if !yisReg{
				y = loadItem(y.(P0Type))
			}
			putOp(cd, x.(Reg).reg, r, y.(Reg).reg)
			releaseReg(y.(Reg).reg)
		}
	return x
}

func genVar(x P0Type) interface{}{
	var y interface{}

	if (0 < x.GetLevel()) && (x.GetLevel() < curlev){
		mark("level")
	}
	_, xisRef := x.(*P0Ref)
	_, xisVar := x.(*P0Var)
	if xisRef{
		y := P0Var{p0type:x.GetP0Type()}
		y.SetLevel(x.GetLevel())

		var regList []string
		regList = append(regList, R0)
		regList = append(regList, A0)
		regList = append(regList, A1)
		regList = append(regList, A2)
		regList = append(regList, A3)

		var xinReg bool
		for i := 0; i < len(regList); i++{
			if x.(P0Ref).GetRegister() == regList[i]{
				xinReg = true
			}
		}
		if xinReg{
			y.SetRegister(x.(P0Ref).GetRegister())
			y.SetAddress(0)
		} else{
			y.SetRegister(obtainReg())
			y.SetAddress(0)

			putMemOp("lw", y.GetRegister(), x.(P0Ref).GetRegister(), x.(P0Ref).GetAddress())
		}
	} else if xisVar {
		var regList []string
		regList = append(regList, R0)
		regList = append(regList, A0)
		regList = append(regList, A1)
		regList = append(regList, A2)
		regList = append(regList, A3)

		var xinReg bool
		for i := 0; i < len(regList); i++{
			if x.(P0Ref).GetRegister() == regList[i]{
				xinReg = true
			}
		}
		if xinReg{
			y := Reg{
				tp:  x.GetP0Type(),
				reg: x.(P0Var).GetRegister(),
			}
		} else{
			y := P0Var{p0type:x.GetP0Type()}
			y.SetLevel(x.GetLevel())
			y.SetRegister(x.(P0Var).GetRegister())
			y.SetAddress(x.(P0Var).GetAddress())
		}

	} else{
		panic("nothing is working")
	}

	return y
}

func genConst(x P0Const) P0Const{
	return x
}

func negate(cd int)int{
	var dict = map[int]int{
		EQ: NE, NE: EQ, LT: GE, LE: GT, GT: LE, GE: LT,
	}
	return dict[cd]
}

func condOp(cd int) string{
	var dict = map[int]string{
		EQ: "beq", NE: "bne", LT: GE, LE: GT, GT: LE, GE: LT,
	}
	return dict[cd]
}

func genUnaryOp(op int, x interface{}) interface{}{
	_, xisVar := x.(*P0Var)
	_, xisCond := x.(*Cond)
	if op == MINUS{
		if xisVar{
			x = loadItem(x.(P0Type))
		}
		putOp("sub", x.(P0Var).GetRegister(), R0, x.(P0Var).GetRegister())
	}else if op == NOT{
		if !xisCond{
			x = loadBool(x.(P0Type))
		}
		x.(Cond).SetCondition(negate(x.(Cond).cond))
		x.(Cond).SetLabA(x.(Cond).labB)
		x.(Cond).SetLabB(x.(Cond).labA)
	} else if op == AND{
		if !xisCond{
			x = loadBool(x.(P0Type))
		}
		str1 := x.(Cond).cond
		str2, _ := strconv.Atoi(str1)
		putBranchOp(condOp(negate(str2)), x.(Cond).left.(string), x.(Cond).right.(string), x.(Cond).labA[0])
		releaseReg(x.(Cond).left.(string))
		releaseReg(x.(Cond).right.(string))
		putLab(x.(Cond).labB, "")
	} else if op == OR{
		if !xisCond{
			x = loadBool(x.(P0Type))
		}
		str1 := x.(Cond).cond
		str2, _ := strconv.Atoi(str1)
		putBranchOp(condOp(str2), x.(Cond).left.(string), x.(Cond).right.(string), x.(Cond).labB[0])
		releaseReg(x.(Cond).left.(string))
		releaseReg(x.(Cond).right.(string))
		putLab(x.(Cond).labA, "")
	} else{
		panic("get unary op failed")
	}
	return x
}

func genBinaryOp(op int, x Cond, y interface{}) interface{}{
	if op == PLUS{
		y = put("add", x, y)
	} else if op == MINUS{
		y = put("sub", x, y)
	} else if op == TIMES{
		y = put("mul", x, y)
	} else if op == DIV{
		y = put("div", x, y)
	} else if op == MOD{
		y = put("mod", x, y)
	} else if op == AND{
		_, yisCond := y.(*Cond)
		if !yisCond{
			y = loadBool(y.(P0Type))
		}
		//todo
	} else if op == OR {
		_, yisCond := y.(*Cond)
		if !yisCond{
			y = loadBool(y.(P0Type))
		}
		//todo
	} else{
		panic("genBinaryOp failed")
	}

	return y
}

func genRelation(op int, x interface{}, y interface{}) Cond{
	_, xisReg := x.(*Reg)
	_, yisReg := y.(*Reg)
	if !xisReg{
		x = loadItem(x.(P0Type))
	}
	if !yisReg{
		y = loadItem(y.(P0Type))
	}
	return NewCond(op, x.(Reg).reg,y.(Reg).reg, ""  )
}

func genSelect(x P0Ref, f P0Var) P0Ref{
	x.p0type = f.p0type
	x.adr = x.adr + f.offset
	return x
}

func genIndex(x interface{}, y interface{}) interface{}{
	_, yisConst := y.(*P0Const)
	if yisConst{
		offset := (y.(P0Const).GetValue().(int) - x.(P0Var).p0type.(int)) * x.(P0Var).GetSize()
		x.(P0Var).SetAddress(x.(P0Var).GetAddress() + offset)
	} else{
		_, yisReg := y.(*Reg)
		if !yisReg{
			y = loadItem(y.(P0Type))
		}
		putOp("sub", y.(Reg).reg, y.(Reg).reg, x.(P0Var).GetP0Type())
		putOp("mul", y.(Reg).reg, y.(Reg).reg, x.(P0Var).GetSize())
		if x.(P0Var).GetRegister() != R0{
			putOp("sub", y.(Reg).reg, x.(P0Var).reg, y.(Reg).reg)
			releaseReg(x.(Reg).reg)
		}
		x.(P0Var).SetRegister(y.(Reg).reg)
	}
	x.(P0Var).p0type = x.(P0Var).GetSize() //idk what to do here
	return x
}

func genAssign(x interface{}, y interface{}){
	
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

func genCall(pr P0Proc){
	putInstr("jal", pr.GetName())
}

func genRead(x P0Var){
	putInstr("li $v0, 5", "")
	putInstr("syscall", "")
	adr := strconv.Itoa(x.GetAddress())
	putMemOp("sw", "$v0", x.GetRegister(), adr)
}

func genWrite(x P0Type){
	loadItemReg(x, "$a0")
	putInstr("li $v0, 1", "")
	putInstr("syscall", "")
}

func genWriteln(){
	putInstr("li $v0, 11", "")
	putInstr("li $a0, '\\n'", "")
	putInstr("syscall", "")
}

func genSeq(){
	//pass
}

func genThen(x interface{}) interface{}{
	_, xisCond := x.(*Cond)
	if !xisCond{
		x = loadBool(x.(P0Type))
	}
	str1 := x.(Cond).cond
	str2, _ := strconv.Atoi(str1)
	putBranchOp(condOp(negate(str2)), x.(Cond).left.(string), x.(Cond).right.(string), x.(Cond).labA[0])
	releaseReg(x.(Cond).left.(string))
	releaseReg(x.(Cond).right.(string))
	putLab(x.(Cond).labB, "")
	return x
}

func genIfThen(x Cond){
	putLab(x.labA, "")
}

func genElse(x Cond) string{
	lab := newLabel()
	putInstr("b", lab)
	putLab(x.labA, "")
	return lab
}

func genIfElse(y[]string){
	putLab(y, "")
}

func genWhile(){
	lab := newLabel()
	var lab1 []string
	lab1 = append(lab1, lab)
	putLab(lab1, "")
}

func genDo(x interface{}) interface{}{
	return genThen(x)
}

func genWhileDo(lab string, x Cond){
	putInstr("b", lab)
	putLab(x.labA, "")
}
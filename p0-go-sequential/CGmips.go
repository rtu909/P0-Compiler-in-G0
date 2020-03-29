package main

import (
	"fmt"
	"strconv"
)

//triple tuple data structure
type Triple struct {
	a, b, c interface{}
}

//global variables
type CGmips struct {
	curlev int
	label  int
	regs   []string
	asm    []Triple
}

var GPRegs = []string{"$t0", "$t1", "$t2", "$t3", "$t4", "$t5", "$t6", "$t7", "$t8"}

//reserved registers
var R0 = "$0"
var FP = "$fp"
var SP = "$sp"
var LNK = "$a3"
var A0 = "$a0"
var A1 = "$a1"
var A2 = "$a2"
var A3 = "$a3"

func (cg *CGmips) GenProgStart() {
	cg.curlev, cg.label = 0, 0
	cg.regs = GPRegs
	cg.putInstr(".data", "")
}

func (cg *CGmips) obtainReg() string {
	if len(cg.regs) == 0 {
		mark("out of registers")
		return R0
	} else {
		var popped = cg.regs[len(cg.regs)-1]
		cg.regs = cg.regs[0 : len(cg.regs)-1]
		return popped
	}
}

func (cg *CGmips) releaseReg(r string) {
	for i := 0; i < len(GPRegs); i++ {
		if r == GPRegs[i] {
			cg.regs = append(cg.regs, r)
		}
	}
}

func (cg *CGmips) putLab(lab []string, instr string) {

	if len(lab) == 1 {
		tuple := Triple{lab[0], instr, ""}
		cg.asm = append(cg.asm, tuple)
	} else {
		fmt.Print("length of lab ", len(lab), "\n")
		for i := 0; i < len(lab)-1; i++ {
			tuple := Triple{lab[i], "", ""}
			cg.asm = append(cg.asm, tuple)
		}
		tuple := Triple{lab[len(lab)-1], instr, ""}
		cg.asm = append(cg.asm, tuple)
	}
}

func (cg *CGmips) putInstr(instr string, target string) {
	tuple := Triple{"", instr, target}
	cg.asm = append(cg.asm, tuple)
}

func (cg *CGmips) putOp(op string, a string, b string, c string) {
	cg.putInstr(op+" "+a+", "+b+", "+c, "")
}

func (cg *CGmips) putBranchOp(op string, a string, b string, c string) {
	cg.putInstr(op+" "+a+", "+b, c)
}

func (cg *CGmips) putMemOp(op string, a string, b string, c string) {
	if b == R0 {
		cg.putInstr(op+" "+a+", "+c, "")
	} else {
		cg.putInstr(op+" "+a+", "+c+"("+b+")", "")
	}
}

//size - not sure what's going on here in the regular code
func (cg *CGmips) GenBool(b P0Type) P0Type {
	b.SetSize(4)
	return b
}

func (cg *CGmips) GenInt(i P0Type) P0Type {
	i.SetSize(4)
	return i
}

//todo
func (cg *CGmips) GenRecord(r P0Type) P0Type {
	s := 0
	fields := r.(*P0Record).GetFields()
	for f := 0; f < len(fields); f++ {
		fields[f].(*P0Var).offset = s
		s = s + fields[f].GetP0Type().GetSize()
	}
	r.SetSize(s)
	return r
}

func (cg *CGmips) GenArray(a P0Type) P0Type {
	size := a.(*P0Array).GetLength() + a.(*P0Array).GetElementType().GetSize()
	a.SetSize(size)
	return a
}

//todo
func (cg *CGmips) GenGlobalVars(declaredVars []Entry, start int) int {
	for i := len(declaredVars) - 1; i > start-1; i-- {
		_, scisVar := declaredVars[i].(*P0Var)
		if scisVar {
			declaredVars[i].(*P0Var).SetRegister(R0)
			declaredVars[i].(*P0Var).SetAddress(declaredVars[i].GetSize())
			var labs []string
			labs = append(labs, strconv.Itoa(declaredVars[i].(*P0Var).GetAddress()))

			cg.putLab(labs, ".space"+strconv.Itoa(declaredVars[i].GetSize()))
		}
	}
	cg.putInstr(".text", "")
	return 0
}

func (cg *CGmips) GenProgEntry() {
	cg.putInstr(".globl main", "")
	cg.putInstr(".ent main", "")
	var lab []string
	lab = append(lab, "main")
	cg.putLab(lab, "")
}

func (cg *CGmips) assembly(l string, i string, t string) string {
	string1 := ""
	if l != "" {
		string1 = l + ":\t"
	} else {
		string1 = "\t"
	}
	string2 := ""
	if t != "" {
		string2 = ", " + t
	} else {
		string2 = ""
	}
	string3 := string1 + i + string2
	return string3
}

func (cg *CGmips) GenProgExit() string {
	cg.putInstr("li $v0, 10", "")
	cg.putInstr("syscall", "")
	cg.putInstr(".end main", "")
	returnStr := ""
	for i := 0; i < len(cg.asm); i++ {
		asm_l := cg.asm[i].a
		asm_i := cg.asm[i].b
		asm_t := cg.asm[i].c
		returnStr = returnStr + cg.assembly(asm_l.(string), asm_i.(string), asm_t.(string)) + "\n"
	}
	return (returnStr)
}

func (cg *CGmips) newLabel() string {
	cg.label = cg.label + 1
	return ("L" + strconv.Itoa(cg.label))
}

// Reg is used like an Entry in the symbol table, so Reg needs to implement the same interface
type Reg struct {
	tp  P0Type
	reg string
}

func (reg *Reg) GetP0Type() P0Type {
	return reg.tp
}

func (reg *Reg) GetName() string {
	return ""
}

func (reg *Reg) SetName(string) {
}

func (reg *Reg) GetLevel() int {
	return -2 // Using this b/c this is what level variables on stack are in the wasm code
}

func (reg *Reg) SetLevel(int) {

}

func (reg *Reg) GetSize() int {
	return 0 // TODO: reconsider?
}

func NewReg(tp P0Type, reg string) Reg {
	r := Reg{
		tp:  tp,
		reg: reg,
	}
	return r
}

// Cond is used like a symbol table Entry, so it implements the same interface
type Cond struct {
	tp          interface{}
	cond        string
	left, right interface{}
	labA        []string
	labB        []string
}

func (cond *Cond) GetP0Type() P0Type {
	return nil
}

func (cond *Cond) GetName() string {
	return ""
}

func (cond *Cond) SetName(string) {

}

func (cond *Cond) GetLevel() int {
	return -69 // TODO: change to something sensible
}

func (cond *Cond) SetLevel(int) {

}

func (cond *Cond) GetSize() int {
	return 0 // If this doesn't make sense, chagne ti to something that does make sense
}

func (condition *Cond) SetCondition(newCond string) {
	(*condition).cond = newCond
}

func (condition *Cond) SetLabA(newCond []string) {
	(*condition).labA = newCond
}

func (condition *Cond) SetLabB(newCond []string) {
	(*condition).labB = newCond
}

func (cg *CGmips) NewCond(tp interface{}, cond string, left interface{}, right interface{}) Cond {
	var labA []string
	var labB []string
	labA = append(labA, cg.newLabel())
	labB = append(labB, cg.newLabel())
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

func (cg *CGmips) testRange(x interface{}) {
	if (x.(*P0Const).GetValue().(int) >= 0x8000) || (x.(*P0Const).GetValue().(int) < -0x8000) {
		mark("value too large")
	}
}

//todo
func (cg *CGmips) loadItemReg(x interface{}, r string) {
	_, xisVar := x.(*P0Var)
	_, xisConst := x.(*P0Const)
	_, xisReg := x.(*Reg)
	if xisVar {
		cg.putMemOp("lw", r, x.(*P0Var).GetRegister(), strconv.Itoa(x.(*P0Var).GetAddress()))
		cg.releaseReg(x.(*P0Var).GetRegister())
	} else if xisConst {
		cg.testRange(x)
		cg.putOp("addi", r, R0, strconv.Itoa(x.(*P0Const).GetValue().(int)))
	} else if xisReg {
		cg.putOp("add", r, x.(Reg).reg, R0)
	} else {
		panic("loadItemReg has problems")
	}
}

//todo
func (cg *CGmips) loadItem(x interface{}) *Reg {
	_, xisConst := x.(*P0Const)
	r := ""
	if xisConst && x.(*P0Const).GetValue() == 0 {
		r = R0
	} else {
		r = cg.obtainReg()
		cg.loadItemReg(x, r)
	}
	return &Reg{x.(P0Type).GetP0Type(), r}
}

//todo
func (cg *CGmips) loadBool(x interface{}) interface{} {
	_, xisConst := x.(*P0Const)
	r := ""
	if xisConst && x.(*P0Const).GetValue() == 0 {
		r = R0
	} else {
		r := cg.obtainReg()
		cg.loadItemReg(x, r)
	}
	return cg.NewCond(NE, r, R0, "")
}

func (cg *CGmips) put(cd string, x interface{}, y interface{}) interface{} {
	_, xisReg := x.(*Reg)
	r := ""
	if !xisReg {
		x = cg.loadItem(x.(P0Type))
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
			x.reg = cg.obtainReg()
		} else {
			r = x.reg
		}
	}
	_, yisConst := y.(*P0Const)
	if yisConst {
		cg.testRange(y)
		cg.putOp(cd, x.(*Reg).reg, r, strconv.Itoa(y.(*P0Const).GetValue().(int)))
	} else {
		_, yisReg := y.(*Reg)
		if !yisReg {
			y = cg.loadItem(y.(P0Type))
		}
		cg.putOp(cd, x.(*Reg).reg, r, y.(*Reg).reg)
		cg.releaseReg(y.(*Reg).reg)
	}
	return x
}

func (cg *CGmips) GenVar(x Entry) Entry {
	var y interface{}

	if (0 < x.GetLevel()) && (x.GetLevel() < cg.curlev) {
		mark("level")
	}
	_, xisRef := x.(*P0Ref)
	_, xisVar := x.(*P0Var)
	if xisRef {
		y = &P0Var{p0type: x.GetP0Type()}
		y.(*P0Var).SetLevel(x.GetLevel())

		var regList []string
		regList = append(regList, R0)
		regList = append(regList, A0)
		regList = append(regList, A1)
		regList = append(regList, A2)
		regList = append(regList, A3)

		var xinReg bool
		for i := 0; i < len(regList); i++ {
			if x.(*P0Ref).GetRegister() == regList[i] {
				xinReg = true
			}
		}
		if xinReg {
			y.(*P0Var).SetRegister(x.(*P0Ref).GetRegister())
			y.(*P0Var).SetAddress(0)
		} else {
			y.(*P0Var).SetRegister(cg.obtainReg())
			y.(*P0Var).SetAddress(0)

			cg.putMemOp("lw", y.(*P0Var).GetRegister(), x.(*P0Ref).GetRegister(), strconv.Itoa(x.(*P0Ref).GetAddress()))
		}
	} else if xisVar {
		var regList []string
		regList = append(regList, R0)
		regList = append(regList, A0)
		regList = append(regList, A1)
		regList = append(regList, A2)
		regList = append(regList, A3)

		var xinReg bool
		for i := 0; i < len(regList); i++ {
			if x.(*P0Var).GetRegister() == regList[i] {
				xinReg = true
			}
		}
		if xinReg {
			y = &Reg{
				tp:  x.GetP0Type(),
				reg: x.(*P0Var).GetRegister(),
			}
		} else {
			y = &P0Var{p0type: x.GetP0Type()}
			y.(*P0Var).SetLevel(x.GetLevel())
			y.(*P0Var).SetRegister(x.(*P0Var).GetRegister())
			y.(*P0Var).SetAddress(x.(*P0Var).GetAddress())
		}

	} else {
		panic("nothing is working")
	}

	return y.(Entry)
}

func (cg *CGmips) GenConst(x Entry) Entry {
	return x
}

func negate(cd int) int {
	var dict = map[int]int{
		EQ: NE, NE: EQ, LT: GE, LE: GT, GT: LE, GE: LT,
	}
	return dict[cd]
}

func condOp(cd int) string {
	var dict = map[int]string{
		EQ: "beq", NE: "bne", LT: "blt", LE: "ble", GT: "bgt", GE: "bge",
	}
	return dict[cd]
}

func (cg *CGmips) GenUnaryOp(op int, entry Entry) Entry {
	_, xisVar := entry.(*P0Var)
	_, xisCond := entry.(*Cond)
	if op == MINUS {
		if xisVar {
			entry = cg.loadBool(entry).(Entry)
		}
		cg.putOp("sub", entry.(*P0Var).GetRegister(), R0, entry.(*P0Var).GetRegister())
	} else if op == NOT {
		if !xisCond {
			entry = cg.loadBool(entry).(Entry)
		}
		str1 := entry.(*Cond).cond
		str2, _ := strconv.Atoi(str1)
		entry.(*Cond).SetCondition(strconv.Itoa(negate(str2)))
		entry.(*Cond).SetLabA(entry.(*Cond).labB)
		entry.(*Cond).SetLabB(entry.(*Cond).labA)
	} else if op == AND {
		if !xisCond {
			entry = cg.loadBool(entry).(Entry)
		}
		str1 := entry.(*Cond).cond
		str2, _ := strconv.Atoi(str1)
		cg.putBranchOp(condOp(negate(str2)), entry.(*Cond).left.(string), entry.(*Cond).right.(string), entry.(*Cond).labA[0])
		cg.releaseReg(entry.(*Cond).left.(string))
		cg.releaseReg(entry.(*Cond).right.(string))
		cg.putLab(entry.(*Cond).labB, "")
	} else if op == OR {
		if !xisCond {
			entry = cg.loadBool(entry).(Entry)
		}
		str1 := entry.(*Cond).cond
		str2, _ := strconv.Atoi(str1)
		cg.putBranchOp(condOp(str2), entry.(*Cond).left.(string), entry.(*Cond).right.(string), entry.(*Cond).labB[0])
		cg.releaseReg(entry.(*Cond).left.(string))
		cg.releaseReg(entry.(*Cond).right.(string))
		cg.putLab(entry.(*Cond).labA, "")
	} else {
		panic("get unary op failed")
	}
	return entry
}

func (cg *CGmips) GenBinaryOp(op int, x Entry, y Entry) Entry {
	if op == PLUS {
		y = cg.put("add", x, y).(Entry)
	} else if op == MINUS {
		y = cg.put("sub", x, y).(Entry)
	} else if op == TIMES {
		y = cg.put("mul", x, y).(Entry)
	} else if op == DIV {
		y = cg.put("div", x, y).(Entry)
	} else if op == MOD {
		y = cg.put("mod", x, y).(Entry)
	} else if op == AND {
		_, yisCond := y.(*Cond)
		if !yisCond {
			y = cg.loadBool(y).(Entry)
		}
		for i := 0; i < len(x.(*Cond).labA); i++ {
			y.(*Cond).SetLabA(append(y.(*Cond).labA, x.(*Cond).labA[i])) // FIXME:
		}
	} else if op == OR {
		_, yisCond := y.(*Cond)
		if !yisCond {
			y = cg.loadBool(y).(Entry)
		}
		for i := 0; i < len(x.(*Cond).labB); i++ {
			y.(*Cond).SetLabB(append(y.(*Cond).labB, x.(*Cond).labB[i])) // FIXME:
		}
	} else {
		panic("genBinaryOp failed")
	}

	return y
}

func (cg *CGmips) GenRelation(op int, x Entry, y Entry) Entry {
	_, xisReg := x.(*Reg)
	_, yisReg := y.(*Reg)
	if !xisReg {
		x = cg.loadItem(x.(P0Type))
	}
	if !yisReg {
		y = cg.loadItem(y.(P0Type))
	}
	x_Entry := cg.NewCond(op, x.(*Reg).reg, y.(*Reg).reg, "")
	return &x_Entry
}

func (cg *CGmips) GenSelect(record Entry, field Entry) Entry {
	record.(*P0Var).p0type = field.(*P0Var).p0type
	record.(*P0Var).SetAddress(record.(*P0Var).GetAddress() + record.(*P0Var).GetOffset())
	return record
}

func (cg *CGmips) GenIndex(x Entry, y Entry) Entry {
	_, yisConst := y.(*P0Const)
	if yisConst {
		offset := (y.(*P0Const).GetValue().(int) - x.(*P0Var).GetP0Type().(*P0Array).lower) * x.(*P0Var).GetSize()
		x.(*P0Var).SetAddress(x.(*P0Var).GetAddress() + offset)
	} else {
		_, yisReg := y.(*Reg)
		if !yisReg {
			y = cg.loadItem(y.(P0Type))
		}
		cg.putOp("sub", y.(*Reg).reg, y.(*Reg).reg, strconv.Itoa(x.(*P0Var).GetP0Type().(*P0Array).lower))
		cg.putOp("mul", y.(*Reg).reg, y.(*Reg).reg, strconv.Itoa(x.(*P0Var).GetSize()))
		if x.(*P0Var).GetRegister() != R0 {
			cg.putOp("sub", y.(*Reg).reg, x.(*P0Var).reg, y.(*Reg).reg)
			cg.releaseReg(x.(*P0Var).GetRegister())
		}
		x.(*P0Var).SetRegister(y.(*Reg).reg)
	}
	//p_0type := x.(*P0Array).GetElementType()
	x = &P0Ref{x.(*P0Array).GetElementType(), x.GetName(), x.GetLevel(), "", 0, 0}
	return x
}

func (cg *CGmips) GenAssign(x, y Entry) {
	_, xisVar := x.(*P0Var)
	_, xisReg := x.(*Reg)
	r := ""
	if xisVar {
		_, yisVar := y.(*Cond)
		_, yisReg := y.(*Reg)
		if yisVar {
			str1 := y.(*Cond).cond
			str2, _ := strconv.Atoi(str1)
			cg.putBranchOp(condOp(str2), y.(*Cond).left.(string), y.(*Cond).right.(string), y.(*Cond).labA[0])
			cg.releaseReg(y.(*Cond).left.(string))
			cg.releaseReg(y.(*Cond).right.(string))
			r = cg.obtainReg()
			cg.putLab(y.(*Cond).labB, "")
			cg.putOp("addi", r, R0, strconv.Itoa(1))
			var lab_list []string
			lab := cg.newLabel()
			lab_list = append(lab_list, lab)
			cg.putInstr("b", lab)
			cg.putLab(y.(*Cond).labA, "")
			cg.putOp("addi", r, R0, strconv.Itoa(0))
			cg.putLab(lab_list, "")
		} else if !yisReg {
			y = cg.loadItem(y.(P0Type))
			r = y.(*Reg).reg
		} else {
			r = y.(*Reg).reg
		}
		cg.putMemOp("sw", r, x.(*P0Var).GetRegister(), strconv.Itoa(x.(*P0Var).GetAddress()))
		cg.releaseReg(r)
	} else if xisReg {
		_, yisVar := y.(*Cond)
		_, yisReg := y.(*Reg)
		if yisVar {
			str1 := y.(*Cond).cond
			str2, _ := strconv.Atoi(str1)
			cg.putBranchOp(condOp(str2), y.(*Cond).left.(string), y.(*Cond).right.(string), y.(*Cond).labA[0])
			cg.releaseReg(y.(*Cond).left.(string))
			cg.releaseReg(y.(*Cond).right.(string))
			cg.putLab(y.(*Cond).labB, "")
			cg.putOp("addi", x.(*Reg).reg, R0, strconv.Itoa(1))
			var lab_list []string
			lab := cg.newLabel()
			lab_list = append(lab_list, lab)
			cg.putInstr("b", lab)
			cg.putLab(y.(*Cond).labA, "")
			cg.putOp("addi", x.(*Reg).reg, R0, strconv.Itoa(0))
			cg.putLab(lab_list, "")
		} else if !yisReg {
			cg.loadItemReg(y.(P0Type), x.(*Reg).reg)
		} else {
			cg.putOp("addi", x.(*Reg).reg, y.(*Reg).reg, strconv.Itoa(0))
		}
	} else {
		panic("genAssign not working")
	}
}

func (cg *CGmips) GenLocalVars(sc []Entry, start int) int {
	s := 0
	for i := start; i < len(sc); i++ {
		_, scIsVar := sc[i].(*P0Var)
		if scIsVar {
			s = s + sc[i].(*P0Var).GetSize()
			sc[i].(*P0Var).SetRegister(FP)
			sc[i].(*P0Var).SetAddress(-s - 0)
		}
	}
	return s
}

func (cg *CGmips) GenProcStart(unused string, fp []Entry) int {
	cg.curlev = cg.curlev + 1
	n := len(fp)
	for i := 0; i < n; i++ {
		_, fpisInt := fp[i].GetP0Type().(*P0Int)
		_, fpisBool := fp[i].GetP0Type().(*P0Bool)
		_, fpisRef := fp[i].(*P0Ref)
		if fpisInt || fpisBool || fpisRef {
			if fpisInt || fpisBool {
				if i < 4 {
					fp[i].(*P0Var).SetRegister("$a" + strconv.Itoa(i))
					fp[i].(*P0Var).SetAddress(0)
				} else {
					fp[i].(*P0Var).SetRegister(FP)
					fp[i].(*P0Var).SetAddress((n - i - 1) * 4)
				}
			} else if fpisRef {
				if i < 4 {
					fp[i].(*P0Ref).SetRegister("$a" + strconv.Itoa(i))
					fp[i].(*P0Ref).SetAddress(0)
				} else {
					fp[i].(*P0Ref).SetRegister(FP)
					fp[i].(*P0Ref).SetAddress((n - i - 1) * 4)
				}
			}
		} else {
			mark("no structured value parameters")
		}
	}
	if (n-4)*4 > 0 {
		return (n - 4) * 4
	} else {
		return 0
	}
}

func (cg *CGmips) GenProcEntry(ident string, parsize int, localsize int) {
	cg.putInstr(".globl"+ident, "")
	cg.putInstr(".ent"+ident, "")
	var lab_list []string
	lab_list = append(lab_list, ident)
	cg.putLab(lab_list, "")
	cg.putMemOp("sw", FP, SP, strconv.Itoa(-parsize-4))
	cg.putMemOp("sw", LNK, SP, strconv.Itoa(-parsize-8))
	cg.putOp("sub", FP, SP, strconv.Itoa(parsize))
	cg.putOp("sub", SP, FP, strconv.Itoa(localsize+8))
}

func (cg *CGmips) GenProcExit(unused Entry, parsize, localsize int) {
	cg.curlev = cg.curlev - 1
	cg.putOp("add", SP, FP, strconv.Itoa(parsize))
	cg.putMemOp("lw", LNK, FP, strconv.Itoa(-8))
	cg.putMemOp("lw", FP, FP, strconv.Itoa(-4))
	cg.putInstr("jr $ra", "")
}

func (cg *CGmips) GenActualPara(ap, fp Entry, n int) {
	_, fpisRef := fp.(*P0Ref)
	r := ""
	if fpisRef {
		if ap.(*P0Var).GetAddress() == 0 {
			if n < 4 {
				cg.putOp("sw", ap.(*P0Var).GetRegister(), SP, strconv.Itoa(-4*(n+1-4)))
			} else {
				cg.putMemOp("sw", ap.(*P0Var).GetRegister(), SP, strconv.Itoa(-4*(n+1-4)))
			}
			cg.releaseReg(ap.(*P0Var).GetRegister())
		} else {
			if n < 4 {
				cg.putMemOp("la", "$a"+strconv.Itoa(n), ap.(*P0Var).GetRegister(), strconv.Itoa(ap.(*P0Var).GetAddress()))
			} else {
				r = cg.obtainReg()
				cg.putMemOp("la", r, ap.(*P0Var).GetRegister(), strconv.Itoa(ap.(*P0Var).GetAddress()))
				cg.putMemOp("sw", r, SP, strconv.Itoa(-4*(n+1-4)))
				cg.releaseReg(r)
			}
		}
	} else {
		_, apisCond := ap.(*Cond)
		_, apisReg := ap.(*Reg)
		if !apisCond {
			if n < 4 {
				cg.loadItemReg(ap, "$a"+strconv.Itoa(n))
			} else {
				if !apisReg {
					ap = cg.loadItem(ap)
				}
				cg.putMemOp("sw", ap.(*P0Var).GetRegister(), SP, strconv.Itoa(-4*(n+1-4)))
				cg.releaseReg(ap.(*Reg).reg)
			}
		} else {
			mark("unsupported parameter type")
		}

	}
}

func (cg *CGmips) GenCall(procedure Entry) {
	cg.putInstr("jal", procedure.GetName())
}

func (cg *CGmips) GenRead(x Entry) {
	cg.putInstr("li $v0, 5", "")
	cg.putInstr("syscall", "")
	adr := strconv.Itoa(x.(*P0Var).GetAddress())
	cg.putMemOp("sw", "$v0", x.(*P0Var).GetRegister(), adr)
}

func (cg *CGmips) GenWrite(x Entry) {
	cg.loadItemReg(x, "$a0")
	cg.putInstr("li $v0, 1", "")
	cg.putInstr("syscall", "")
}

func (cg *CGmips) GenWriteln() {
	cg.putInstr("li $v0, 11", "")
	cg.putInstr("li $a0, '\\n'", "")
	cg.putInstr("syscall", "")
}

func (cg *CGmips) GenSeq(Entry, Entry) {
	//pass
}

func (cg *CGmips) GenThen(x Entry) Entry {
	_, xisCond := x.(*Cond)
	if !xisCond {
		val := cg.loadBool(x).(Cond)
		x = &val
	}
	str1 := x.(*Cond).cond
	str2, _ := strconv.Atoi(str1)
	cg.putBranchOp(condOp(negate(str2)), x.(*Cond).left.(string), x.(*Cond).right.(string), x.(*Cond).labA[0])
	cg.releaseReg(x.(*Cond).left.(string))
	cg.releaseReg(x.(*Cond).right.(string))
	cg.putLab(x.(*Cond).labB, "")
	return x
}

func (cg *CGmips) GenIfThen(x Entry) {
	cg.putLab(x.(*Cond).labA, "")
}

func (cg *CGmips) GenElse(x, y Entry) string {
	lab := cg.newLabel()
	cg.putInstr("b", lab)
	cg.putLab(x.(*Cond).labA, "")
	return lab
}

func (cg *CGmips) GenIfElse(label string) {
	arr := make([]string, 1)
	arr[0] = label
	cg.putLab(arr, "")
}

func (cg *CGmips) GenWhile() string {
	lab := cg.newLabel()
	var lab1 []string
	lab1 = append(lab1, lab)
	cg.putLab(lab1, "")
	return lab
}

func (cg *CGmips) GenDo(x Entry) Entry {
	return cg.GenThen(x)
}

func (cg *CGmips) GenWhileDo(lab string, x, y Entry) {
	cg.putInstr("b", lab)
	cg.putLab(x.(*Cond).labA, "")
}

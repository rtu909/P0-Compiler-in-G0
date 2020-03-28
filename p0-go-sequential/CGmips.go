package main

import (
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

func (cg * CGmips) obtainReg() string {
	if len(cg.regs) == 0 {
		mark("out of registers")
		return R0
	} else {
		var popped = cg.regs[9]
		cg.regs = cg.regs[0 : len(cg.regs)-1]
		return popped
	}
}

func (cg * CGmips) releaseReg(r string) {
	for i := 0; i < len(GPRegs); i++ {
		if r == GPRegs[i] {
			cg.regs = append(cg.regs, r)
		}
	}
}

func (cg * CGmips) putLab(lab []string, instr string) {

	if len(lab) == 1 {
		tuple := Triple{lab[0], instr, ""}
		cg.asm = append(cg.asm, tuple)
	} else {
		for i := 0; i < len(lab)-1; i++ {
			tuple := Triple{lab[i], "", ""}
			cg.asm = append(cg.asm, tuple)
		}
		tuple := Triple{lab[len(lab)-1], instr, ""}
		cg.asm = append(cg.asm, tuple)
	}
}

func (cg * CGmips) putInstr(instr string, target string) {
	tuple := Triple{"", instr, target}
	cg.asm = append(cg.asm, tuple)
}

func (cg * CGmips) putOp(op string, a string, b string, c string) {
	cg.putInstr(op+" "+a+", "+b+", "+c, "")
}

func (cg * CGmips) putBranchOp(op string, a string, b string, c string) {
	cg.putInstr(op+" "+a+", "+b, c)
}

func (cg * CGmips) putMemOp(op string, a string, b string, c string) {
	if b == R0 {
		cg.putInstr(op+" "+a+", "+c, "")
	} else {
		cg.putInstr(op+" "+a+", "+c+"("+b+")", "")
	}
}

//size - not sure what's going on here in the regular code
func (cg * CGmips) genBool(b P0Bool) P0Bool {
	b.SetSize(4)
	return b
}

func (cg * CGmips) genInt(i P0Int) P0Int {
	i.SetSize(4)
	return i
}

func genRec(r P0Record) P0Record {
	s := 0
	fields := r.GetFields()
	for f := 0; f < len(fields); f++ {

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
func genGlobalVars(sc []P0Var, start int) {
	for i := len(sc) - 1; i > start-1; i-- {

	}
	putInstr(".text", "")
}

func genProgEntry() {
	putInstr(".globl main", "")
	putInstr(".ent main", "")
	var lab []string
	putLab(lab, "main")
}

func assembly(l string, i string, t string) string {
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

func genProgExit() string {
	putInstr("li $v0, 10", "")
	putInstr("syscall", "")
	putInstr(".end main", "")
	returnStr := ""
	for i := 0; i < len(asm); i++ {
		asm_l := asm[i].a
		asm_i := asm[i].b
		asm_t := asm[i].c
		returnStr = returnStr + assembly(asm_l.(string), asm_i.(string), asm_t.(string)) + "\n"
	}
	return (returnStr)
}

func newLabel() string {
	label = label + 1
	return ("L" + strconv.Itoa(label))
}

type Reg struct {
	tp  interface{}
	reg string
}

func NewReg(tp interface{}, reg string) Reg {
	r := Reg{
		tp:  tp,
		reg: reg,
	}
	return r
}

type Cond struct {
	tp          interface{}
	cond        string
	left, right interface{}
	labA        []string
	labB        []string
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

func NewCond(tp interface{}, cond string, left interface{}, right interface{}) Cond {
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

func testRange(x P0Const) {
	if (x.GetValue().(int) >= 0x8000) || (x.GetValue().(int) < -0x8000) {
		mark("value too large")
	}
}

//todo
func loadItemReg(x interface{}, r string) {
	_, xisVar := x.(*P0Var)
	_, xisConst := x.(*P0Const)
	_, xisReg := x.(*Reg)
	if xisVar {
		putMemOp("lw", r, x.(P0Var).GetRegister(), strconv.Itoa(x.(P0Var).GetAddress()))
		releaseReg(x.(P0Var).GetRegister())
	} else if xisConst {
		testRange(x.(P0Const))
		putOp("addi", r, R0, strconv.Itoa(x.(P0Const).GetValue().(int)))
	} else if xisReg {
		putOp("add", r, x.(Reg).reg, R0)
	} else {
		panic("loadItemReg has problems")
	}
}

//todo
func loadItem(x interface{}) Reg {
	_, xisConst := x.(*P0Const)
	r := ""
	if xisConst && x.(P0Const).GetValue() == 0 {
		r = R0
	} else {
		r = obtainReg()
		loadItemReg(x, r)
	}
	return Reg{x.(P0Const).GetP0Type(), r}
}

//todo
func loadBool(x interface{}) Cond {
	_, xisConst := x.(*P0Const)
	r := ""
	if xisConst && x.(P0Const).GetValue() == 0 {
		r = R0
	} else {
		r := obtainReg()
		loadItemReg(x, r)
	}
	return NewCond(NE, r, R0, "")
}

func put(cd string, x interface{}, y interface{}) interface{} {
	_, xisReg := x.(*Reg)
	r := ""
	if !xisReg {
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
	if yisConst {
		y := y.(P0Const)
		testRange(y)
		putOp(cd, x.(Reg).reg, r, y.GetValue().(string))
	} else {
		_, yisReg := y.(*Reg)
		if !yisReg {
			y = loadItem(y.(P0Type))
		}
		putOp(cd, x.(Reg).reg, r, y.(Reg).reg)
		releaseReg(y.(Reg).reg)
	}
	return x
}

func genVar(x Entry) interface{} {
	var y interface{}

	if (0 < x.GetLevel()) && (x.GetLevel() < curlev) {
		mark("level")
	}
	_, xisRef := x.(*P0Ref)
	_, xisVar := x.(*P0Var)
	if xisRef {
		y := P0Var{p0type: x.GetP0Type()}
		y.SetLevel(x.GetLevel())

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
			y.SetRegister(x.(*P0Ref).GetRegister())
			y.SetAddress(0)
		} else {
			y.SetRegister(obtainReg())
			y.SetAddress(0)

			putMemOp("lw", y.GetRegister(), x.(*P0Ref).GetRegister(), strconv.Itoa(x.(*P0Ref).GetAddress()))
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
			if x.(*P0Ref).GetRegister() == regList[i] {
				xinReg = true
			}
		}
		if xinReg {
			y = Reg{
				tp:  x.GetP0Type(),
				reg: x.(*P0Var).GetRegister(),
			}
		} else {
			y := P0Var{p0type: x.GetP0Type()}
			y.SetLevel(x.GetLevel())
			y.SetRegister(x.(*P0Var).GetRegister())
			y.SetAddress(x.(*P0Var).GetAddress())
		}

	} else {
		panic("nothing is working")
	}

	return y
}

func genConst(x P0Const) P0Const {
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

func genUnaryOp(op int, x interface{}) interface{} {
	_, xisVar := x.(*P0Var)
	_, xisCond := x.(*Cond)
	if op == MINUS {
		if xisVar {
			x = loadItem(x.(P0Type))
		}
		putOp("sub", x.(P0Var).GetRegister(), R0, x.(P0Var).GetRegister())
	} else if op == NOT {
		if !xisCond {
			x = loadBool(x.(P0Type))
		}
		str1 := x.(Cond).cond
		str2, _ := strconv.Atoi(str1)
		x.(Cond).SetCondition(strconv.Itoa(negate(str2)))
		x.(Cond).SetLabA(x.(Cond).labB)
		x.(Cond).SetLabB(x.(Cond).labA)
	} else if op == AND {
		if !xisCond {
			x = loadBool(x.(P0Type))
		}
		str1 := x.(Cond).cond
		str2, _ := strconv.Atoi(str1)
		putBranchOp(condOp(negate(str2)), x.(Cond).left.(string), x.(Cond).right.(string), x.(Cond).labA[0])
		releaseReg(x.(Cond).left.(string))
		releaseReg(x.(Cond).right.(string))
		putLab(x.(Cond).labB, "")
	} else if op == OR {
		if !xisCond {
			x = loadBool(x.(P0Type))
		}
		str1 := x.(Cond).cond
		str2, _ := strconv.Atoi(str1)
		putBranchOp(condOp(str2), x.(Cond).left.(string), x.(Cond).right.(string), x.(Cond).labB[0])
		releaseReg(x.(Cond).left.(string))
		releaseReg(x.(Cond).right.(string))
		putLab(x.(Cond).labA, "")
	} else {
		panic("get unary op failed")
	}
	return x
}

func genBinaryOp(op int, x Cond, y interface{}) interface{} {
	if op == PLUS {
		y = put("add", x, y)
	} else if op == MINUS {
		y = put("sub", x, y)
	} else if op == TIMES {
		y = put("mul", x, y)
	} else if op == DIV {
		y = put("div", x, y)
	} else if op == MOD {
		y = put("mod", x, y)
	} else if op == AND {
		_, yisCond := y.(*Cond)
		if !yisCond {
			y = loadBool(y.(P0Type))
		}
		for i := 0; i < len(x.labA); i++ {
			y.(Cond).SetLabA(append(y.(Cond).labA, x.labA[i]))
		}
	} else if op == OR {
		_, yisCond := y.(*Cond)
		if !yisCond {
			y = loadBool(y.(P0Type))
		}
		for i := 0; i < len(x.labB); i++ {
			y.(Cond).SetLabB(append(y.(Cond).labB, x.labB[i]))
		}
	} else {
		panic("genBinaryOp failed")
	}

	return y
}

func genRelation(op int, x interface{}, y interface{}) Cond {
	_, xisReg := x.(*Reg)
	_, yisReg := y.(*Reg)
	if !xisReg {
		x = loadItem(x.(P0Type))
	}
	if !yisReg {
		y = loadItem(y.(P0Type))
	}
	return NewCond(op, x.(Reg).reg, y.(Reg).reg, "")
}

func genSelect(x P0Ref, f P0Var) P0Ref {
	x.p0type = f.p0type
	x.adr = x.adr + f.offset
	return x
}

func genIndex(x Entry, y interface{}) interface{} {
	_, yisConst := y.(*P0Const)
	if yisConst {
		offset := (y.(P0Const).GetValue().(int) - x.(*P0Var).GetP0Type().(*P0Array).lower) * x.(*P0Var).GetSize()
		x.(*P0Var).SetAddress(x.(*P0Var).GetAddress() + offset)
	} else {
		_, yisReg := y.(*Reg)
		if !yisReg {
			y = loadItem(y.(P0Type))
		}
		putOp("sub", y.(Reg).reg, y.(Reg).reg, strconv.Itoa(x.(*P0Var).GetP0Type().(*P0Array).lower))
		putOp("mul", y.(Reg).reg, y.(Reg).reg, strconv.Itoa(x.(*P0Var).GetSize()))
		if x.(*P0Var).GetRegister() != R0 {
			putOp("sub", y.(Reg).reg, x.(*P0Var).reg, y.(Reg).reg)
			releaseReg(x.(*P0Var).GetRegister())
		}
		x.(*P0Var).SetRegister(y.(Reg).reg)
	}
	//p_0type := x.(*P0Array).GetElementType()
	x = &P0Ref{x.(*P0Array).GetElementType(), x.GetName(), x.GetLevel(), "", 0, 0}
	return x
}

func genAssign(x interface{}, y interface{}) {
	_, xisVar := x.(*P0Var)
	_, xisReg := x.(*Reg)
	r := ""
	if xisVar {
		_, yisVar := y.(*Cond)
		_, yisReg := y.(*Reg)
		if yisVar {
			str1 := y.(Cond).cond
			str2, _ := strconv.Atoi(str1)
			putBranchOp(condOp(str2), y.(Cond).left.(string), y.(Cond).right.(string), y.(Cond).labA[0])
			releaseReg(y.(Cond).left.(string))
			releaseReg(y.(Cond).right.(string))
			r = obtainReg()
			putLab(y.(Cond).labB, "")
			putOp("addi", r, R0, strconv.Itoa(1))
			var lab_list []string
			lab := newLabel()
			lab_list = append(lab_list, lab)
			putInstr("b", lab)
			putLab(y.(Cond).labA, "")
			putOp("addi", r, R0, strconv.Itoa(0))
			putLab(lab_list, "")
		} else if !yisReg {
			y = loadItem(y.(P0Type))
			r = y.(Reg).reg
		} else {
			r = y.(Reg).reg
		}
		putMemOp("sw", r, x.(P0Var).GetRegister(), strconv.Itoa(x.(P0Var).GetAddress()))
		releaseReg(r)
	} else if xisReg {
		_, yisVar := y.(*Cond)
		_, yisReg := y.(*Reg)
		if yisVar {
			str1 := y.(Cond).cond
			str2, _ := strconv.Atoi(str1)
			putBranchOp(condOp(str2), y.(Cond).left.(string), y.(Cond).right.(string), y.(Cond).labA[0])
			releaseReg(y.(Cond).left.(string))
			releaseReg(y.(Cond).right.(string))
			putLab(y.(Cond).labB, "")
			putOp("addi", x.(Reg).reg, R0, strconv.Itoa(1))
			var lab_list []string
			lab := newLabel()
			lab_list = append(lab_list, lab)
			putInstr("b", lab)
			putLab(y.(Cond).labA, "")
			putOp("addi", x.(Reg).reg, R0, strconv.Itoa(0))
			putLab(lab_list, "")
		} else if !yisReg {
			loadItemReg(y.(P0Type), x.(Reg).reg)
		} else {
			putOp("addi", x.(Reg).reg, y.(Reg).reg, strconv.Itoa(0))
		}
	} else {
		panic("genAssign not working")
	}
}

func genLocalVars(sc []Entry, start int) int {
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

func genProcStart(fp []Entry) int {
	curlev = curlev + 1
	n := len(fp)
	for i := 0; i < n; i++ {
		_, fpisInt := fp[i].(*P0Int)
		_, fpisBool := fp[i].(*P0Bool)
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

func genProcEntry(ident string, parsize int, localsize int) {
	putInstr(".globl"+ident, "")
	putInstr(".ent"+ident, "")
	var lab_list []string
	lab_list = append(lab_list, ident)
	putLab(lab_list, "")
	putMemOp("sw", FP, SP, strconv.Itoa(-parsize-4))
	putMemOp("sw", LNK, SP, strconv.Itoa(-parsize-8))
	putOp("sub", FP, SP, strconv.Itoa(parsize))
	putOp("sub", SP, FP, strconv.Itoa(localsize+8))
}

func genProcExit(parsize int, localsize int) {
	curlev = curlev - 1
	putOp("add", SP, FP, strconv.Itoa(parsize))
	putMemOp("lw", LNK, FP, strconv.Itoa(-8))
	putMemOp("lw", FP, FP, strconv.Itoa(-4))
	putInstr("jr $ra", "")
}

func genActualPara(ap interface{}, fp Entry, n int) {
	_, fpisRef := fp.(*P0Ref)
	r := ""
	if fpisRef {
		if ap.(P0Var).GetAddress() == 0 {
			if n < 4 {
				putOp("sw", ap.(P0Var).GetRegister(), SP, strconv.Itoa(-4*(n+1-4)))
			} else {
				putMemOp("sw", ap.(P0Var).GetRegister(), SP, strconv.Itoa(-4*(n+1-4)))
			}
			releaseReg(ap.(P0Var).GetRegister())
		} else {
			if n < 4 {
				putMemOp("la", "$a"+strconv.Itoa(n), ap.(P0Var).GetRegister(), strconv.Itoa(ap.(P0Var).GetAddress()))
			} else {
				r = obtainReg()
				putMemOp("la", r, ap.(P0Var).GetRegister(), strconv.Itoa(ap.(P0Var).GetAddress()))
				putMemOp("sw", r, SP, strconv.Itoa(-4*(n+1-4)))
				releaseReg(r)
			}
		}
	} else {
		_, apisCond := ap.(*Cond)
		_, apisReg := ap.(*Reg)
		if !apisCond {
			if n < 4 {
				loadItemReg(ap, "$a"+strconv.Itoa(n))
			} else {
				if !apisReg {
					ap = loadItem(ap)
				}
				putMemOp("sw", ap.(P0Var).GetRegister(), SP, strconv.Itoa(-4*(n+1-4)))
				releaseReg(ap.(Reg).reg)
			}
		} else {
			mark("unsupported parameter type")
		}

	}
}

func genCall(pr P0Proc) {
	putInstr("jal", pr.GetName())
}

func genRead(x P0Var) {
	putInstr("li $v0, 5", "")
	putInstr("syscall", "")
	adr := strconv.Itoa(x.GetAddress())
	putMemOp("sw", "$v0", x.GetRegister(), adr)
}

func genWrite(x P0Type) {
	loadItemReg(x, "$a0")
	putInstr("li $v0, 1", "")
	putInstr("syscall", "")
}

func genWriteln() {
	putInstr("li $v0, 11", "")
	putInstr("li $a0, '\\n'", "")
	putInstr("syscall", "")
}

func genSeq() {
	//pass
}

func genThen(x interface{}) interface{} {
	_, xisCond := x.(*Cond)
	if !xisCond {
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

func genIfThen(x Cond) {
	putLab(x.labA, "")
}

func genElse(x Cond) string {
	lab := newLabel()
	putInstr("b", lab)
	putLab(x.labA, "")
	return lab
}

func genIfElse(y []string) {
	putLab(y, "")
}

func genWhile() {
	lab := newLabel()
	var lab1 []string
	lab1 = append(lab1, lab)
	putLab(lab1, "")
}

func genDo(x interface{}) interface{} {
	return genThen(x)
}

func genWhileDo(lab string, x Cond) {
	putInstr("b", lab)
	putLab(x.labA, "")
}

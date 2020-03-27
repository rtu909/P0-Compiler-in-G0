package main

type WasmGenerator struct {
	currentLevel int
	memorySize   int
	asm          []string
}

func (wg *WasmGenerator) GenProgStart() {
	wg.currentLevel = 0
	wg.memorySize = 0
	wg.asm = make([]string, 0)
	wg.asm = append(wg.asm, "(module")
	wg.asm = append(wg.asm, "(import \"P0lib\" \"write\" (func $write (param i32)))")
	wg.asm = append(wg.asm, "(import \"P0lib\" \"writeln\" (func $writeln))")
	wg.asm = append(wg.asm, "(import \"P0lib\" \"print\" (func $print (param i32)))")
	wg.asm = append(wg.asm, "(import \"P0lib\" \"read\" (func $read (result i32)))")
}

func (wg *WasmGenerator) GenBool(p0Type P0Type) P0Type {
	p0Type.SetSize(1)
	return p0Type
}

func (wg *WasmGenerator) GenInt(p0Type P0Type) P0Type {
	p0Type.SetSize(4)
	return p0Type
}

func (wg *WasmGenerator) GenRecord(p0Type P0Type) P0Type {
	p0record := p0Type.(*P0Record)
	var sum int = 0
	for _, param := range p0record.GetFields() {
		sum += param.GetSize()
	}
	p0record.SetSize(sum)
	return p0record
}

func (wg *WasmGenerator) GenArray(p0Type P0Type) P0Type {
	p0array := p0Type.(*P0Array)
	p0array.SetSize(p0array.GetElementType().GetSize() * p0array.GetLength())
	return p0array
}

// GenGlobalVars creates the code for declaring the global variables at the start of the file.
// sc is a map of names to types. It represents all the global variable declarations.
// start Represents TODO:
func (wg *WasmGenerator) GenGlobalVars(sc []Entry, start int) {
	for _, entry := range sc {
		asVar, isVar := entry.(*P0Var)
		if isVar {
			switch asVar.GetP0Type().(type) {
			case *P0Bool:
			case *P0Int:
				(*wg).asm = append((*wg).asm, "(global $"+entry.GetName()+" (mut i32) i32.const 0)")
				break
			case *P0Array:
			case *P0Record:
				asVar.SetLevel(-2)
				asVar.SetAddress(wg.memorySize)
				wg.memorySize += entry.GetP0Type().GetSize()
				break
			default:
				mark("WASM type?")
			}
		}
	}
}

func (wg *WasmGenerator) GenLocalVars(sc []Entry, start int) {
	for _, entry := range sc {
		asVar, isVar := entry.(*P0Var)
		if isVar {
			switch asVar.GetP0Type().(type) {
			case *P0Int:
			case *P0Bool:
				(*wg).asm = append((*wg).asm, "(local $"+entry.GetName()+" i32)")
				break
			case *P0Array:
			case *P0Record:
				mark("WASM: no local arrays, records")
				break
			default:
				mark("WASM type?")
			}
		}
	}
}

// The returned value is nil
func (wg *WasmGenerator) LoadItem(entry Entry) Entry {
	asVar, isVar := entry.(*P0Var)
	if isVar {
		if asVar.GetLevel() == 0 {
			wg.asm = append(wg.asm, "global get $"+asVar.GetName())
		} else if asVar.GetLevel() == wg.currentLevel {
			wg.asm = append(wg.asm, "local.get $"+asVar.GetName())
		} else if asVar.GetLevel() == -2 {
			wg.asm = append(wg.asm, "i32.const "+string(asVar.GetAddress()))
			wg.asm = append(wg.asm, "i32.load")
		}
	} else {
		asRef, isRef := entry.(*P0Ref)
		if isRef {
			if asRef.GetLevel() == -1 {
				wg.asm = append(wg.asm, "i32.load")
			} else if entry.GetLevel() == wg.currentLevel {
				wg.asm = append(wg.asm, "i32.local $"+asRef.GetName())
				wg.asm = append(wg.asm, "i32.load")
			} else {
				mark("WASM: ref level!")
			}
		} else {
			asConst, isConst := entry.(*P0Const)
			if isConst {
				wg.asm = append(wg.asm, "i32.const "+string(asConst.GetValue().(int)))
			}
		}
	}
	return nil
}

func (wg *WasmGenerator) GenVar(entry Entry) Entry {
	if 0 < entry.GetLevel() && entry.GetLevel() < wg.currentLevel {
		mark("WASM: level!")
	}
	var newEntry Entry
	_, isRef := entry.(*P0Ref)
	if isRef {
		newEntry = &P0Ref{entry.GetP0Type(), entry.GetName(), entry.GetLevel(), "", 0, 0}
	} else {
		asVar, isVar := entry.(*P0Var)
		if isVar {
			newEntry = &P0Var{entry.GetP0Type(), entry.GetName(), entry.GetLevel(), "", 0, 0}
			if entry.GetLevel() == -2 {
				newEntry.(*P0Var).SetAddress(asVar.GetAddress())
			}
		}
	}
	return newEntry
}

func (wg *WasmGenerator) GenConst(entry Entry) Entry {
	return entry
}

func (wg *WasmGenerator) GenUnaryOp(op int, entry Entry) Entry {
	wg.LoadItem(entry)
	switch op {
	case MINUS:
		wg.asm = append(wg.asm, "i32.const -1")
		wg.asm = append(wg.asm, "i32.mul")
		entry = &P0Var{&P0Int{}, "", -1, "", 0, 0}
		break
	case NOT:
		wg.asm = append(wg.asm, "i32.eqz")
		entry = &P0Var{&P0Bool{}, "", -1, "", 0, 0}
		break
	case AND:
		wg.asm = append(wg.asm, "if (result i32)")
		entry = &P0Var{&P0Bool{}, "", -1, "", 0, 0}
		break
	case OR:
		wg.asm = append(wg.asm, "if (result i32)")
		wg.asm = append(wg.asm, "i32.const 1")
		wg.asm = append(wg.asm, "else")
		entry = &P0Var{&P0Bool{}, "", -1, "", 0, 0}
		break
	default:
		mark("WASM: unary operator?")
	}
	return entry
}

func (wg *WasmGenerator) GenBinaryOp(op int, x, y Entry) Entry {
	switch op {
	case PLUS:
	case MINUS:
	case TIMES:
	case DIV:
	case MOD:
		wg.LoadItem(x)
		wg.LoadItem(y)
		switch op {
		case PLUS:
			wg.asm = append(wg.asm, "i32.add")
			break
		case MINUS:
			wg.asm = append(wg.asm, "i32.sub")
			break
		case TIMES:
			wg.asm = append(wg.asm, "i32.mul")
			break
		case DIV:
			wg.asm = append(wg.asm, "i32.div_s")
			break
		case MOD:
			wg.asm = append(wg.asm, "i32.rem_s")
			break
		}
		x = &P0Var{&P0Int{}, "", -1, "", 0, 0}
		break
	case AND:
		wg.LoadItem(y)
		wg.asm = append(wg.asm, "else")
		wg.asm = append(wg.asm, "i32.const 0")
		wg.asm = append(wg.asm, "end")
		x = &P0Var{&P0Bool{}, "", -1, "", 0, 0}
		break
	case OR:
		// x should already be on the stack b/c magic
		wg.LoadItem(y)
		wg.asm = append(wg.asm, "end")
		x = &P0Var{&P0Bool{}, "", -1, "", 0, 0}
		break
	default:
		panic("Unrecognized binary operator")
	}
	return x
}

func (wg *WasmGenerator) GenRelation(op int, x Entry, y Entry) Entry {
	wg.LoadItem(x)
	wg.LoadItem(y)
	switch op {
	case EQ:
		wg.asm = append(wg.asm, "i32.eq")
		break
	case NE:
		wg.asm = append(wg.asm, "i32.ne")
		break
	case LT:
		wg.asm = append(wg.asm, "i32.lt_s")
		break
	case GT:
		wg.asm = append(wg.asm, "i32.gt_s")
		break
	case LE:
		wg.asm = append(wg.asm, "i32.le_s")
		break
	case GE:
		wg.asm = append(wg.asm, "i32.ge_s")
		break
	default:
		panic("Unrecognized relational operator")
	}
	x = &P0Var{&P0Bool{}, "", -1, "", 0, 0}
	return x
}

// Assuming x is a Record and f is a field of x
func (wg *WasmGenerator) GenSelect(x Entry, f Entry) Entry {
	asVar, isVar := x.(*P0Var)
	if isVar {
		asVar.SetAddress(asVar.GetAddress() + f.(*P0Var).GetOffset()) // TODO: make sure that parameters are vars
		asVar.p0type = f.GetP0Type()
		return asVar
	} else {
		asRef, isRef := x.(*P0Ref)
		if isRef {
			if x.GetLevel() > 0 {
				wg.asm = append(wg.asm, "local.get $"+x.GetName())
			}
			wg.asm = append(wg.asm, "i32.const "+string(f.(*P0Var).GetOffset()))
			wg.asm = append(wg.asm, "i32.add")
			asRef.SetLevel(-1)
			asRef.p0type = f.GetP0Type()
			return asRef
		}
	}
	panic("Should only call with variable of reference")
}

func (wg *WasmGenerator) GenIndex(x Entry, y Entry) Entry {
	xAsVar, xIsVar := x.(*P0Var)
	arrayType := x.GetP0Type().(*P0Array)
	if xIsVar {
		yAsConst, yIsConst := y.(*P0Const)
		if yIsConst {
			xAsVar.SetAddress(xAsVar.GetAddress() +
				(yAsConst.GetValue().(int)-arrayType.GetLowerBound())*
					arrayType.GetElementType().GetSize())
			xAsVar.p0type = xAsVar.GetP0Type().(*P0Array).GetElementType()
			return xAsVar
		} else {
			wg.LoadItem(y)
			if arrayType.GetLowerBound() != 0 {
				wg.asm = append(wg.asm, "i32.const "+string(arrayType.GetLowerBound()))
				wg.asm = append(wg.asm, "i32.sub")
			}
			wg.asm = append(wg.asm, "i32.const "+string(arrayType.GetElementType().GetSize()))
			wg.asm = append(wg.asm, "i32.mul")
			wg.asm = append(wg.asm, "i32.const "+string(xAsVar.GetAddress()))
			wg.asm = append(wg.asm, "i32.add")
			x = &P0Ref{arrayType.GetElementType(), "", -1, "", 0, 0}
		}
	} else {
		if x.GetLevel() == wg.currentLevel {
			wg.LoadItem(x)
			x.SetLevel(-1)
		}
		yAsConst, yIsConst := y.(*P0Const)
		if yIsConst {
			wg.asm = append(wg.asm, "i32.const "+
				string((yAsConst.GetValue().(int)-arrayType.GetLowerBound())*arrayType.GetElementType().GetSize()))
			wg.asm = append(wg.asm, "i32.add")
		} else {
			wg.LoadItem(y)
			wg.asm = append(wg.asm, "i32.const "+string(arrayType.GetLowerBound()))
			wg.asm = append(wg.asm, "i32.sub")
			wg.asm = append(wg.asm, "i32.const "+string(arrayType.GetElementType().GetSize()))
			wg.asm = append(wg.asm, "i32.mul")
			wg.asm = append(wg.asm, "i32.add")
		}
		x = &P0Ref{arrayType.GetElementType(), x.GetName(), x.GetLevel(), "", 0, 0}
	}
	return x
}

func (wg *WasmGenerator) GenAssign(x, y Entry) {
	xAsVar, xIsVar := x.(*P0Var)
	xAsRef, xIsRef := x.(*P0Ref)
	if xIsVar {
		if xAsVar.GetLevel() == -2 {
			wg.asm = append(wg.asm, "i32.const "+string(xAsVar.GetAddress()))
		}
		wg.LoadItem(y)
		if xAsVar.GetLevel() == 0 {
			wg.asm = append(wg.asm, "global.set $"+x.GetName())
		} else if xAsVar.GetLevel() == wg.currentLevel {
			wg.asm = append(wg.asm, "local.set $"+x.GetName())
		} else if xAsVar.GetLevel() == -2 {
			wg.asm = append(wg.asm, "i32.store")
		} else {
			mark("WASM: level!")
		}
	} else if xIsRef {
		if xAsRef.GetLevel() == wg.currentLevel {
			wg.asm = append(wg.asm, "local.get $"+x.GetName())
		}
		wg.LoadItem(y)
		wg.asm = append(wg.asm, "i32.store")
	} else {
		panic("The generator was passed code to assign to an unchangeable value")
	}
}

func (wg *WasmGenerator) GenProgEntry() { // NOTE: originally had an unused parameter `ident`
	wg.asm = append(wg.asm, "(func $program")
}

func (wg *WasmGenerator) GenProgExit() string {
	wg.asm = append(wg.asm, ")\n(memory "+string(wg.memorySize/(2<<16)+1)+")\n(start $program)\n")
	var theCode string = ""
	for _, line := range wg.asm {
		theCode += line
	}
	return theCode
}

// GenProcStart generates the beginning of a procedure declaration
// ident is the name of the procedure
// fp is a slice holding the formal parameters of the procedure
func (wg *WasmGenerator) GenProcStart(ident string, fp []Entry) int {
	if wg.currentLevel > 0 {
		mark("WASM: no nested procedures")
	}
	wg.currentLevel++
	funcDecl := "(func $" + ident + " "
	for _, entry := range fp {
		_, isVar := entry.(*P0Var)
		_, isRef := entry.(*P0Ref)
		switch entry.GetP0Type().(type) {
		case *P0Int:
		case *P0Bool:
			if isRef {
				mark("WASM: Only array and record reference parameters")
			}
			break
		case *P0Array:
		case *P0Record:
			if isVar {
				mark("WASM: no structured valued parameters")
			}
			break
		}
		funcDecl += "(param $" + entry.GetName() + " i32) "
	}
	wg.asm = append(wg.asm, funcDecl)
	return 0
}

func (wg *WasmGenerator) GenProcEntry(ident string, parsize, localsize int) {
}

func (wg *WasmGenerator) GenProcExit(x Entry, parsize, localsize int) {
	wg.currentLevel--
	wg.asm = append(wg.asm, ")")
}

// ap is the actual parameter
// fp is the formal parameter
func (wg *WasmGenerator) GenActualPara(ap, fp Entry, parameterNumber int) {
	_, asRef := fp.(*P0Ref)
	if asRef {
		// Assume that ap is a Var
		if ap.GetLevel() == -2 {
			wg.asm = append(wg.asm, "i32.const "+string(ap.(*P0Var).GetAddress()))
		}
	} else {
		switch ap.(type) {
		case *P0Var:
		case *P0Ref:
		case *P0Const:
			wg.LoadItem(ap)
			break
		default:
			mark("Unsupported parameter type")
		}
	}
}

func (wg *WasmGenerator) GenCall(pr Entry) {
	wg.asm = append(wg.asm, "call $"+pr.GetName())
}

func (wg *WasmGenerator) GenRead(x Entry) {
	wg.asm = append(wg.asm, "call $read")
	// Dr. Sekerinski's 'hack' from the email I sent him
	y := &P0Var{&P0Int{}, "", -1, "", 0, 0}
	wg.GenAssign(x, y)
}

func (wg *WasmGenerator) GenWrite(x Entry) {
	wg.LoadItem(x)
	wg.asm = append(wg.asm, "call $write")
}

func (wg *WasmGenerator) GenWriteln() {
	wg.asm = append(wg.asm, "call $writeln")
}

func (wg *WasmGenerator) GenSeq(x, y Entry) {
}

func (wg *WasmGenerator) GenThen(x Entry) Entry {
	wg.LoadItem(x)
	wg.asm = append(wg.asm, "if")
	return x
}

func (wg *WasmGenerator) GenIfThen(x Entry) {
	wg.asm = append(wg.asm, "end")
}

func (wg *WasmGenerator) GenElse(x, y Entry) {
	wg.asm = append(wg.asm, "else")
}

func (wg *WasmGenerator) GenIfElse(x, y, z Entry) {
	wg.asm = append(wg.asm, "end")
}

func (wg *WasmGenerator) GenWhile() {
	wg.asm = append(wg.asm, "loop")
}

func (wg *WasmGenerator) GenDo(x Entry) Entry {
	wg.LoadItem(x)
	wg.asm = append(wg.asm, "if")
	return x
}

func (wg *WasmGenerator) GenWhileDo(t, x, y Entry) {
	wg.asm = append(wg.asm, "br 1")
	wg.asm = append(wg.asm, "end")
	wg.asm = append(wg.asm, "end")
}

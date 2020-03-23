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
func (wg *WasmGenerator) GenGlobalVars(sc map[string]Entry, start int) {
	for _, entry := range sc {
		asVar, isVar := entry.(*P0Var)
		if isVar {
			switch _ := asVar.GetP0Type().(type) {
			case *P0Bool:
			case *P0Int:
				(*wg).asm = append((*wg).asm, "(global $"+entry.GetName()+" (mut i32) i32.const 0)")
				break
			case *P0Array:
			case *P0Record:
				// TODO: need to add an adr instance variable to the declaration
				asVar.SetLevel(-2)
				wg.memorySize += entry.GetP0Type().GetSize()
				break
			default:
				mark("WASM type?")
			}
		}
	}
}

func (wg *WasmGenerator) GenLocalVars(sc map[string]Entry, start int) {
	for declName, entry := range sc {
		asVar, isVar := entry.(*P0Var)
		if isVar {
			switch _ := asVar.GetP0Type().(type) {
			case *P0Int:
			case *P0Bool:
				(*wg).asm = append((*wg).asm, "(local $"+declName+" i32)")
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

func (wg *WasmGenerator) LoadItem(entry Entry) {
	asVar, isVar := entry.(*P0Var)
	if isVar {
		if asVar.GetLevel() == 0 {
			wg.asm = append(wg.asm, "global get $"+asVar.GetName())
		} else if asVar.GetLevel() == wg.currentLevel {
			wg.asm = append(wg.asm, "local.get $"+asVar.GetName())
		} else if asVar.GetLevel() == -2 {
			wg.asm = append(wg.asm, "i32.const " /* TODO: +asVar.GetAddress() */)
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
}

func (wg *WasmGenerator) GenVar(entry Entry) Entry {
	if 0 < entry.GetLevel() && entry.GetLevel() < wg.currentLevel {
		mark("WASM: level!")
	}
	var newEntry Entry
	_, isRef := entry.(*P0Ref)
	if isRef {
		newEntry = &P0Ref{entry.GetP0Type(), entry.GetName(), entry.GetLevel()}
	} else {
		_, isVar := entry.(*P0Var)
		if isVar {
			newEntry = &P0Var{entry.GetP0Type(), entry.GetName(), entry.GetLevel()}
			if entry.GetLevel() == -2 {
				// TODO: copy the address of the old entry to the new entry
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
		entry = &P0Var{&P0Int{}, "", -1} // WHY? I don't know why this is done
		break
	case NOT:
		wg.asm = append(wg.asm, "i32.eqz")
		entry = &P0Var{&P0Bool{}, "", -1}
		break
	case AND:
		wg.asm = append(wg.asm, "if (result i32)")
		entry = &P0Var{&P0Bool{}, "", -1}
		break
	case OR:
		wg.asm = append(wg.asm, "if (result i32)")
		wg.asm = append(wg.asm, "i32.const 1")
		wg.asm = append(wg.asm, "else")
		entry = &P0Var{&P0Bool{}, "", -1}
		break
	default:
		mark("WASM: unary operator?")
	}
	return entry
}

func (wg *WasmGenerator) GenBinaryOP(op int, x Entry, y Entry) Entry {
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
		x = &P0Var{&P0Int{}, "", -1}
		break
	case AND:
		wg.LoadItem(y)
		wg.asm = append(wg.asm, "else")
		wg.asm = append(wg.asm, "i32.const 0")
		wg.asm = append(wg.asm, "end")
		x = &P0Var{&P0Bool{}, "", -1}
		break
	case OR:
		// x should already be on the stack b/c magic
		wg.LoadItem(y)
		wg.asm = append(wg.asm, "end")
		x = &P0Var{&P0Bool{}, "", -1}
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
	x = &P0Var{&P0Bool{}, "", -1}
	return x
}

// Assuming x is a Record and f is a field of x
func (wg *WasmGenerator) GenSelect(x Entry, f Entry) Entry {
	asVar, isVar := x.(*P0Var)
	if isVar {
		// TODO: x.SetAddress(x.GetAddress() + f.GetOffset())
		asVar.p0type = f.GetP0Type()
		return asVar
	} else {
		asRef, isRef := x.(*P0Ref)
		if isRef {
			if x.GetLevel() > 0 {
				wg.asm = append(wg.asm, "local.get $"+x.GetName())
			}
			wg.asm = append(wg.asm, "i32.const " /* TODO: + f.GetOffset()*/)
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
			// TODO:
			// x.SetAddress(x.GetAddress() +
			// (yAsConst.GetValue() - arrayType.GetElementType().GetLowerBound()) *
			// arrayType.GetElementType().GetSize())
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
			wg.asm = append(wg.asm, "i32.const " /* TODO: + string(x.GetAddress()) */)
			wg.asm = append(wg.asm, "i32.add")
			x = &P0Ref{arrayType.GetElementType(), "", -1}
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
		x = &P0Ref{arrayType.GetElementType(), x.GetName(), x.GetLevel()}
	}
	return x
}

func (wg *WasmGenerator) GenAssign(x, y Entry) {
	xAsVar, xIsVar := x.(*P0Var)
	xAsRef, xIsRef := x.(*P0Ref)
	if xIsVar {
		if xAsVar.GetLevel() == -2 {
			wg.asm = append(wg.asm, "i32.const " /* TODO: x.GetAddress()*/)
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

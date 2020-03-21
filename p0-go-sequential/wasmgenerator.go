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
	s := 0
	for _, field := range p0Type.typeComponents {
		s += field.GetSize()
	}
	p0Type.SetSize(s)
	return p0Type
}

func (wg *WasmGenerator) GenArray(p0Type P0Type) P0Type {
	p0Type.SetSize(p0Type.typeComponents[0].GetSize() * p0Type.GetArrayLength())
	return p0Type
}

// GenGlobalVars creates the code for declaring the global variables at the start of the file.
// sc is a map of names to types. It represents all the global variable declarations.
// start Represents TODO:
func (wg *WasmGenerator) GenGlobalVars(sc map[string]Entry, start int) {
	for declName, entry := range sc {
		_, isVar := entry.(P0Var)
		if isVar {
			if entry.GetP0Type().p0primitive == Int || entry.GetP0Type().p0primitive == Bool {
				(*wg).asm = append((*wg).asm, "(global $"+declName+" (mut i32) i32.const 0)")
			} else if entry.GetP0Type().p0primitive == Array || entry.GetP0Type().p0primitive == Record {
				// TODO: a bunch of new instance variables are added to the entry in the original code, but I haven't
				// added them to the entry yet, so we can't assign them. Also don't understand how they are used.
				(*wg).memorySize += entry.GetP0Type().GetSize()
			} else {
				mark("WASM type?")
			}
		}
	}
}

func (wg *WasmGenerator) GenLocalVars(sc map[string]Entry, start int) {
	for declName, entry := range sc {
		_, isVar := entry.(P0Var)
		if isVar {
			if entry.GetP0Type().p0primitive == Int || entry.GetP0Type().p0primitive == Bool {
				(*wg).asm = append((*wg).asm, "(local $"+declName+" i32)")
			} else if entry.GetP0Type().p0primitive == Record || entry.GetP0Type().p0primitive == Array {
				mark("WASM: no local arrays, records")
			} else {
				mark("WASM type?")
			}
		}
	}
}

// LoadItem loads an item TODO: how, why, what are all these variables?
func (wg *WasmGenerator) LoadItem(declaration map[string]Entry) {
	// TODO: figure out the level of the entry, then load it
	asVar, isVar := item.(P0Var)
	if isVar {
		if asVar.lev == 0 {
			(*wg).asm = append((*wg).asm, "global get $"+asVar.name)
		} else if asVar.lev == curlev {
			(*wg).asm = append((*wg).asm, "local.get $"+asVar.name)
		} else if asVar.lev == -2 {
			(*wg).asm = append((*wg).asm, "i32.const "+asVar.adr)
			(*wg).asm = append((*wg).asm, "i32.load")
		}
	} else {
		asRef, isRef := item.(P0Ref)
		if isRef {
			if asRef.lev == -1 {
				(*wg).asm = append((*wg).asm, "i32.load")
			} else if x.lev == curlev {
				(*wg).asm = append((*wg).asm, "i32.local $"+name)
				(*wg).asm = append((*wg).asm, "i32.load")
			} else {
				mark("WASM: ref level!")
			}
		} else {
			asConst, isConst := item.(P0Const)
			if isConst {
				(*wg).asm = append((*wg).asm, "i32.const "+string(asConst.value.(int)))
			}
		}
	}
}

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
	//p0Type.size = 1 // TODO: add this to P0Type
	return p0Type
}

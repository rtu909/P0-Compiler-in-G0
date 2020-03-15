package main

import (

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

}

func putLab(lab []string){

}

func putInstr(instr string, target string){

}

func putOp(op string, a string, b string, c string){

}

func putBranchOp(op string, a string, b string, c string){

}

func putMemOp(op string, a string, b string, c string){

}

//size

func genBool(){

}

func genInt(){

}

func genRec(){

}

func genArray(){
	
}


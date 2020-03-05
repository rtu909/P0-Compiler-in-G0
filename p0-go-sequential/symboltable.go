package main

// SymbolTable keeps track of declarations that are made in P0 source code.
type SymbolTable interface {
	Init()
	NewDecl(string, Entry)
	Find(string)
	OpenScope()
	TopScope()
	CloseScope()
}

// DeclType represents the different values that identifiers can be declared as in P0.
// It is an enumerated type
type DeclType int

const (
	Var DeclType = iota
	Ref
	Const
	Type
	Proc
	StdProc
)

// Entry represents items that can be put into the SymbolTable. All entries must have a DeclType (but it can be nil)
// Some entries contain more data
type Entry interface {
	GetDeclType() DeclType
}

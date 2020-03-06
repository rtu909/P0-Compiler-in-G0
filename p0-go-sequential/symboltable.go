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

// P0Type represents one of the built-in types in P0
// Current implementation is an enumerated type; this will have to be changed, because Arrays also consist of a type
// that is held, and records are sequences of types
type P0Type int

const (
	Int P0Type = iota
	Bool
	Enum
	Record
	Array
)

// Entry represents items that can be put into the SymbolTable. All entries must have a DeclType (but it can be nil)
// Some entries contain more data
type Entry interface {
	GetP0Type() P0Type
}

// Var represents an entry in the symbol table for a P0 variable
type P0Var P0Type

func (p0var P0Var) GetP0Type() P0Type {
	return P0Type(p0var)
}

// ArraySymbolTable implements the symbol table as an array of maps from strings to the Entry
type ArraySymbolTable []map[string]Entry

// Init initialized an ArraySymbolTable for use
// It sets up the outermost context for use
func (st *ArraySymbolTable) Init() {
	*st = make([]map[string]Entry, 1)
	(*st)[0] = make(map[string]Entry)
}

// NewDecl adds a new declaration to the symbol table at the current level
func (st *ArraySymbolTable) NewDecl(name string, entry Entry) {
	_, present = (*st)[len(*st)-1][name]
	if !present {
		(*st)[len(*st)-1][name] = entry
	} else {
		println("Multiple definition")
	}
}

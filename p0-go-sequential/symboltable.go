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
	None
)

// Entry represents items that can be put into the SymbolTable. All entries must have a DeclType (but it can be nil)
// Some entries contain more data
type Entry interface {
	GetP0Type() P0Type
}

// P0Var represents an entry in the symbol table for a P0 variable
type P0Var P0Type

func (p0var P0Var) GetP0Type() P0Type {
	return P0Type(p0var)
}

// P0Const represents an identifier that s linked to a constant value
type P0Const struct {
	p0type P0Type
	value  interface{} //TODO: what needs to go here?
}

func (p0const P0Const) GetP0Type() P0Type {
	return p0const.p0type
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
	_, present := (*st)[len(*st)-1][name]
	if !present {
		(*st)[len(*st)-1][name] = entry
	} else {
		println("Multiple definition")
	}
}

// Find attempts to find the innermost declaration of the symbol `name` in the symbol table.
// If it is found, the entry is returned
// If it is not found, a P0Const with value 0 and P0Type None is returned
func (st *ArraySymbolTable) Find(name string) Entry {
	for i := len(*st) - 1; i >= 0; i-- {
		entry, present := (*st)[i][name]
		if present {
			return entry
		}
	}
	println("Cannot find symbol")
	return P0Const{None, 0}
}

// OpenScope opens a new (innermost) declaration scoping.
// This means that symbols defined in the current scope could be redefined
func (st *ArraySymbolTable) OpenScope() {
	*st = append(*st, make(map[string]Entry))
}

// CloseScope closes the innermost scope of the symbol table.
// Any declarations made in this scope are deleted.
// The new innermost scope becomes the old second most inner scope.
func (st *ArraySymbolTable) CloseScope() {
	*st = (*st)[0 : len(*st)-1]
}

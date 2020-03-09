package main

// SymbolTable keeps track of declarations that are made in P0 source code.
// It keeps track of the current scope, and queries can be made against it in order to find what declaration a symbol refers to.
type SymbolTable interface {
	Init()
	NewDecl(string, Entry)
	Find(string) Entry
	OpenScope()
	CloseScope()
	TopScope() // TODO: how to implement, b/c the parser needs to work with this directly
}

// P0Primitive is an enumerated type that represents one of the built-in types in P0.
// It is only meant to represent the base types; composite types are represented in P0Type
type P0Primitive int

const (
	Int P0Primitive = iota
	Bool
	Record
	Array
	None
)

// P0Type is a representation of composite data types in P0.
// It consists of the base type, p0primitive, combined with the constituent types, typeComponents.
// If the base type is one of Int, Bool, or None, typeComponents can be nil and should not be accessed.
// If the base type is Array, typeComponents must be of length 1 and contain the type that the array holds.
// If the base type is Record, typeComponents must be of length 1 or greater.
// The values represent the types of the fields, in the order that they appear in the Record.
type P0Type struct {
	p0primitive    P0Primitive
	typeComponents []P0Type
}

// Entry represents items that can be put into the SymbolTable.
// All entries have a type.
type Entry interface {
	GetP0Type() P0Type
	GetFieldNames() []string
	IsConstant() bool
	GetValue() int
	GetLowerBound() int
	GetLength() int
}

// P0Var represents an entry in the symbol table for a P0 variable
type P0Var P0Type

// P0Var implements the Entry interface so that it can be stored in the symbol table.
func (p0var P0Var) GetP0Type() P0Type {
	return P0Type(p0var)
}

// P0Const represents an identifier that is linked to a constant value.
type P0Const struct {
	p0type P0Type
	value  interface{} //TODO: what needs to go here?
}

// P0Const implements the Entry interface so that it can be stored in the symbol table
func (p0const P0Const) GetP0Type() P0Type {
	return p0const.p0type
}

// SliceMapSymbolTable implements the symbol table as a slice of maps from string to Entry
type SliceMapSymbolTable []map[string]Entry

// Init initialized a SliceMapSymbolTable for use.
func (st *SliceMapSymbolTable) Init() {
	*st = make([]map[string]Entry, 1)
	(*st)[0] = make(map[string]Entry)
}

// NewDecl adds a new declaration to the symbol table at the current level.
func (st *SliceMapSymbolTable) NewDecl(name string, entry Entry) {
	_, present := (*st)[len(*st)-1][name]
	if !present {
		(*st)[len(*st)-1][name] = entry
	} else {
		println("Multiple definition")
	}
}

// Find attempts to find the innermost declaration of the symbol `name` in the symbol table.
// If it is found, the corresponding entry is returned.
// If it is not found, a P0Const with value 0 and P0Type None (eg. P0Type{None, nil}) is returned.
func (st *SliceMapSymbolTable) Find(name string) Entry {
	for i := len(*st) - 1; i >= 0; i-- {
		entry, present := (*st)[i][name]
		if present {
			return entry
		}
	}
	println("Cannot find symbol")
	return P0Const{P0Type{None, nil}, 0}
}

// OpenScope opens a new (innermost) declaration scope.
func (st *SliceMapSymbolTable) OpenScope() {
	*st = append(*st, make(map[string]Entry))
}

// CloseScope closes the innermost scope of the symbol table.
// Any declarations made in the scope are deleted.
// The new innermost scope becomes the old second most inner scope.
func (st *SliceMapSymbolTable) CloseScope() {
	*st = (*st)[0 : len(*st)-1]
}

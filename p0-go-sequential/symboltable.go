package main

// SymbolTable keeps track of declarations that are made in P0 source code.
// It keeps track of the current scope, and queries can be made against it in order to find what declaration a symbol refers to.
type SymbolTable interface {
	Init()
	NewDecl(string, Entry)
	Find(string) Entry
	OpenScope()
	CloseScope()
	TopScope() map[string]Entry
}

// Entry represents items that can be put into the SymbolTable.
// All entries have a type, name, and level they are declared at.
// TODO: add methods for:
//  * adr
//  * reg
//  * offset
type Entry interface {
	GetP0Type() P0Type // This is the type related to the entry
	GetName() string   // This is the name of the entry as a string
	SetName(string)    // The symbol table needs to set the name when it is putting the entry in the table
	GetLevel() int     // This is the level on the symbol table where the Entry was declared
	SetLevel(int)      // The generators do some funky stuff with the level of entries, which requires them to modify it
	GetSize() int      // How many bytes this type takes in storage; the generator will need to calculate this
}

// P0Type is a representation of types in P0.
// P0Type implements Entry, since new types can be declared in a P0 program.
// Because of this, it needs to be possible to store P0Types on the symbol table.
type P0Type interface {
	GetP0Type() P0Type // Returns nil, since this itself is a P0Type
	GetName() string
	SetName(string)
	GetLevel() int
	SetLevel(int)
	GetSize() int
	SetSize(int)
}

// Represents the integer type in P0
type P0Int struct {
	name  string
	size  int
	level int
}

func (p0int *P0Int) GetP0Type() P0Type {
	return nil
}

func (p0int *P0Int) GetName() string {
	return p0int.name
}

func (p0int *P0Int) SetName(newName string) {
	p0int.name = newName
}

func (p0int *P0Int) GetLevel() int {
	return p0int.level
}

func (p0int *P0Int) SetLevel(newLevel int) {
	p0int.level = newLevel
}

func (p0int *P0Int) GetSize() int {
	return p0int.size
}

func (p0int *P0Int) SetSize(newSize int) {
	p0int.size = newSize
}

// Represents the boolean type in P0
type P0Bool struct {
	name  string
	size  int
	level int
}

func (p0bool *P0Bool) GetP0Type() P0Type {
	return nil
}

func (p0bool *P0Bool) GetName() string {
	return p0bool.name
}

func (p0bool *P0Bool) SetName(newName string) {
	p0bool.name = newName
}

func (p0bool *P0Bool) GetLevel() int {
	return p0bool.level
}

func (p0bool *P0Bool) SetLevel(newLevel int) {
	p0bool.level = newLevel
}

func (p0bool *P0Bool) GetSize() int {
	return p0bool.size
}

func (p0bool *P0Bool) SetSize(newSize int) {
	p0bool.size = newSize
}

// Represents the enum type in P0
type P0Enum struct {
	name  string
	size  int
	level int
}

func (p0enum *P0Enum) GetP0Type() P0Type {
	return nil
}

func (p0enum *P0Enum) GetName() string {
	return p0enum.name
}

func (p0enum *P0Enum) SetName(newName string) {
	p0enum.name = newName
}

func (p0enum *P0Enum) GetLevel() int {
	return p0enum.level
}

func (p0enum *P0Enum) SetLevel(newLevel int) {
	p0enum.level = newLevel
}

func (p0enum *P0Enum) GetSize() int {
	return p0enum.size
}

func (p0enum *P0Enum) SetSize(newSize int) {
	p0enum.size = newSize
}

// P0Record represents declared record types in P0
type P0Record struct {
	name   string
	size   int
	level  int
	fields []P0Type
}

func (p0record *P0Record) GetP0Type() P0Type {
	return nil
}

func (p0record *P0Record) GetName() string {
	return p0record.name
}

func (p0record *P0Record) SetName(newName string) {
	p0record.name = newName
}

func (p0record *P0Record) GetLevel() int {
	return (*p0record).level
}

func (p0record *P0Record) SetLevel(newLevel int) {
	(*p0record).level = newLevel
}

func (p0record *P0Record) GetSize() int {
	return (*p0record).size
}

func (p0record *P0Record) SetSize(newSize int) {
	(*p0record).size = newSize
}

func (p0record *P0Record) GetFields() []P0Type {
	return (*p0record).fields
}

type P0Array struct {
	name   string
	size   int
	level  int
	base   P0Type
	lower  int // The lowest index that you can reference the array by
	length int // The number of items in the array
}

func (p0array *P0Array) GetP0Type() P0Type {
	return nil
}

func (p0array *P0Array) GetName() string {
	return p0array.name
}

func (p0array *P0Array) SetName(newName string) {
	(*p0array).name = newName
}

func (p0array *P0Array) GetLevel() int {
	return (*p0array).level
}

func (p0array *P0Array) SetLevel(newLevel int) {
	(*p0array).level = newLevel
}

func (p0array *P0Array) GetSize() int {
	return (*p0array).size
}

func (p0array *P0Array) SetSize(newSize int) {
	(*p0array).size = newSize
}

// GetArray gets the length of the array (number of elements in the array).
func (p0array *P0Array) GetLength() int {
	return (*p0array).length
}

// GetLowerBound gets the number that can be used to address the first element in the array.
func (p0array *P0Array) GetLowerBound() int {
	return (*p0array).lower
}

// GetElementType gets the type of the elements that the array contains.
func (p0array *P0Array) GetElementType() P0Type {
	return (*p0array).base
}

// P0Var represents an entry in the symbol table for a P0 variable
// It implements the Entry interface so that it can be stored in the symbol table
type P0Var struct {
	p0type P0Type
	name   string
	level  int
	reg    string // The resgister it is currently stored in
	adr    int    // the location in heap memory it is stored
	offset int    // The offset of the value from the beginning of the structure/array that it is at
}

func (p0var *P0Var) GetP0Type() P0Type {
	return p0var.p0type
}

func (p0var *P0Var) GetName() string {
	return p0var.name
}

func (p0var *P0Var) SetName(newName string) {
	(*p0var).name = newName
}

func (p0var *P0Var) GetSize() int {
	return p0var.p0type.GetSize()
}

func (p0var *P0Var) SetSize(newSize int) {
	print("Changing stored size of variable directly instead of changing size of underlying type; probably bad")
	p0var.GetP0Type().SetSize(newSize)
}

func (p0var *P0Var) GetLevel() int {
	return p0var.level
}

func (p0var *P0Var) SetLevel(newLevel int) {
	(*p0var).level = newLevel
}

func (p0var *P0Var) GetRegister() string {
	return p0var.reg
}

func (p0var *P0Var) SetRegister(newReg string) {
	p0var.reg = newReg
}

func (p0var *P0Var) GetAddress() int {
	return p0var.adr
}

func (p0var *P0Var) SetAddress(newAddress int) {
	p0var.adr = newAddress
}

func (p0var *P0Var) GetOffset() int {
	return p0var.offset
}

func (p0var *P0Var) SetOffset(newOffset int) {
	p0var.offset = newOffset
}

// P0Ref represents an entry in the symbol table for a reference (pointer-like construct) in P0.
// It implements the Entry interface so that it can be stored in the symbol table
type P0Ref struct {
	p0type P0Type
	name   string
	level  int
	reg    string
	adr    int
	offset int
}

func (p0ref *P0Ref) GetP0Type() P0Type {
	return p0ref.p0type
}

func (p0ref *P0Ref) GetName() string {
	return p0ref.name
}

func (p0ref *P0Ref) SetName(newName string) {
	(*p0ref).name = newName
}

func (p0ref *P0Ref) GetSize() int {
	return p0ref.p0type.GetSize()
}

func (p0ref *P0Ref) SetSize(newSize int) {
	print("Changing stored size of variable directly instead of changing size of underlying type; probably bad")
	p0ref.GetP0Type().SetSize(newSize)
}

func (p0ref *P0Ref) GetLevel() int {
	return p0ref.level
}

func (p0ref *P0Ref) SetLevel(newLevel int) {
	(*p0ref).level = newLevel
}

func (p0ref *P0Ref) GetRegister() string {
	return p0ref.reg
}

func (p0ref *P0Ref) SetRegister(newReg string) {
	p0ref.reg = newReg
}

func (p0ref *P0Ref) GetAddress() int {
	return p0ref.adr
}

func (p0ref *P0Ref) SetAddress(newAddress int) {
	p0ref.adr = newAddress
}

func (p0ref *P0Ref) GetOffset() int {
	return p0ref.offset
}

func (p0ref *P0Ref) SetOffset(newOffset int) {
	p0ref.offset = newOffset
}

// P0Const represents an entry in the symbol table for a P0 constant
// It implements the Entry interface so that it can be stored in the symbol table.
// On top of all the information that's necessary for storing it in the symbol table, it stores the associated constant value, as an empty interface type.
type P0Const struct {
	p0type P0Type
	name   string
	level  int
	value  interface{}
}

func (p0const *P0Const) GetP0Type() P0Type {
	return p0const.p0type
}

func (p0const *P0Const) GetName() string {
	return p0const.name
}

func (p0const *P0Const) SetName(newName string) {
	(*p0const).name = newName
}

func (p0const *P0Const) GetSize() int {
	return p0const.p0type.GetSize()
}

func (p0const *P0Const) SetSize(newSize int) {
	print("Changing stored size of constant directly instead of changing size of underlying type; probably bad")
	p0const.GetP0Type().SetSize(newSize)
}

func (p0const *P0Const) GetLevel() int {
	return p0const.level
}

func (p0const *P0Const) SetLevel(newLevel int) {
	(*p0const).level = newLevel
}

func (p0const *P0Const) GetValue() interface{} {
	return p0const.value
}

// P0Proc represents a user-declared procedure in a P0 program.
// It implements Entry so it can be stored on the symbol table.
// It also has methods for accessing the list of parameters that need to be passed.
type P0Proc struct {
	p0type     P0Type
	name       string
	level      int
	parameters []P0Type
}

func (p0proc *P0Proc) GetP0Type() P0Type {
	return nil
}

func (p0proc *P0Proc) GetName() string {
	return p0proc.name
}

func (p0proc *P0Proc) SetName(newName string) {
	p0proc.name = newName
}

// TODO: probably don't need?
func (*P0Proc) GetSize() int {
	return 0
}

// TODO: probably also don't need?
func (*P0Proc) SetSize(int) {
}

func (p0proc *P0Proc) GetLevel() int {
	return p0proc.level
}

func (p0proc *P0Proc) SetLevel(newLevel int) {
	p0proc.level = newLevel
}

func (p0proc *P0Proc) GetParameters() []P0Type {
	return p0proc.parameters
}

// P0Proc represents a user-declared procedure in a P0 program.
// It implements Entry so it can be stored on the symbol table.
// It also has methods for accessing the list of parameters that need to be passed.
type P0StdProc struct {
	p0type     P0Type
	name       string
	level      int
	parameters []P0Type
}

func (p0stdproc *P0StdProc) GetP0Type() P0Type {
	return nil
}

func (p0stdproc *P0StdProc) GetName() string {
	return p0stdproc.name
}

func (p0stdproc *P0StdProc) SetName(newName string) {
	p0stdproc.name = newName
}

// TODO: probably don't need?
func (*P0StdProc) GetSize() int {
	return 0
}

// TODO: probably also don't need?
func (*P0StdProc) SetSize(int) {
}

func (p0stdproc *P0StdProc) GetLevel() int {
	return p0stdproc.level
}

func (p0stdproc *P0StdProc) SetLevel(newLevel int) {
	p0stdproc.level = newLevel
}

func (p0stdproc *P0StdProc) GetParameters() []P0Type {
	return p0stdproc.parameters
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
		entry.SetName(name)
		entry.SetLevel(len(*st) - 1)
		(*st)[len(*st)-1][name] = entry
	} else {
		mark("Multiple definition of " + name)
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
	mark("Cannot find symbol " + name)
	return &P0Const{&P0Int{"int", 0, 0}, "error", 0, 0}
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

// TopScope returns the top scope, as a map from strings to Entries
func (st *SliceMapSymbolTable) TopScope() map[string]Entry {
	return (*st)[0]
}

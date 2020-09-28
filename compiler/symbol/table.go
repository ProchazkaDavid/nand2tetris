package symbol

import (
	"errors"
)

// Identifier represents types used in symbol table.
type Identifier int

// Identifier possibilities.
const (
	Unknown Identifier = iota
	Static
	Field
	Arg
	Var
)

var (
	errUnknownIdentifier = errors.New("unknown identifier")
)

type tableEntry struct {
	varType string
	kind    Identifier
	index   int
}

type symbolTable map[string]tableEntry

// Table represents symbol table.
type Table struct {
	class       symbolTable
	subroutine  symbolTable
	staticIndex int
	fieldIndex  int
	argIndex    int
	varIndex    int
}

// NewSymbolTable creates a new symbol table.
func NewSymbolTable() *Table {
	return &Table{
		class:      make(symbolTable),
		subroutine: make(symbolTable),
	}
}

// NewSubroutine starts a new subroutine scope.
func (t *Table) NewSubroutine() {
	t.argIndex = 0
	t.varIndex = 0
	t.subroutine = make(symbolTable)
}

// Define defines a new identifier of the given name, type and kind,
// and assigns it a running index. KeywordTypes static and field have a class scope,
// while arg and var have a subroutine scope.
func (t *Table) Define(name, varType string, kind Identifier) {
	var index *int = nil

	switch kind {
	case Static:
		index = &t.staticIndex
	case Field:
		index = &t.fieldIndex
	case Arg:
		index = &t.argIndex
	case Var:
		index = &t.varIndex
	default:
		panic(errUnknownIdentifier)
	}

	switch kind {
	case Static, Field:
		t.class[name] = tableEntry{varType, kind, *index}
	case Arg, Var:
		t.subroutine[name] = tableEntry{varType, kind, *index}
	}

	*index++
}

// VariableCount returns the number of variables of the given kind already
// defined in the current scope.
func (t *Table) VariableCount(kind Identifier) int {
	switch kind {
	case Static:
		return t.staticIndex
	case Field:
		return t.fieldIndex
	case Arg:
		return t.argIndex
	case Var:
		return t.varIndex
	default:
		panic(errUnknownIdentifier)
	}
}

// KindOf returns the kind of the names identifier in the current scope.
// If the identifier is unknown in the current scope, returns None.
func (t *Table) KindOf(name string) Identifier {
	if entry, ok := t.subroutine[name]; ok {
		return entry.kind
	}

	if entry, ok := t.class[name]; ok {
		return entry.kind
	}

	return Unknown
}

// TypeOf returns the type of the named identifier in the current scope.
func (t *Table) TypeOf(name string) (varType string, ok bool) {
	if entry, ok := t.subroutine[name]; ok {
		return entry.varType, true
	}

	if entry, ok := t.class[name]; ok {
		return entry.varType, true
	}

	return "", false
}

// IndexOf returns the index assigned to the named identifier.
func (t *Table) IndexOf(name string) int {
	if entry, ok := t.subroutine[name]; ok {
		return entry.index
	}

	if entry, ok := t.class[name]; ok {
		return entry.index
	}

	return -1
}

package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	outer          *SymbolTable
	store          map[string]Symbol
	numDefinations int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)

	return &SymbolTable{store: s}
}

func NewEnclosedSymbolTable(st *SymbolTable) *SymbolTable {
	symbolTable := NewSymbolTable()
	symbolTable.outer = st
	return symbolTable
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinations, Scope: GlobalScope}
	if s.outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinations++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.outer != nil {
		obj, ok = s.outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}

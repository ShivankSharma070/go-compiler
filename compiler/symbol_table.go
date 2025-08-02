package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinations int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)

	return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinations, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDefinations++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}

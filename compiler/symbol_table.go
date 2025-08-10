package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
	FreeScope    SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
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
	FreeSymbols    []Symbol
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	free:= []Symbol{}
	return &SymbolTable{store: s, FreeSymbols: free}
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

		// If Not resolved
		if !ok {
			return obj, ok
		}

		// If resolved
		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		free := s.defineFree(obj)
		return free, true
	}
	return obj, ok
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols)- 1, Scope:FreeScope}
	s.store[original.Name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}
	s.store[name] = symbol
	return symbol
}


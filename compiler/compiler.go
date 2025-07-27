package compiler

import (
	"github.com/ShivankSharma070/go-compiler/ast"
	"github.com/ShivankSharma070/go-compiler/code"
	"github.com/ShivankSharma070/go-compiler/object"
)


type Compiler struct {
instructions code.Instructions
	constants	[]object.Object
}

func New() *Compiler {
return &Compiler{
		instructions: code.Instructions{},
		constants: []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil 
}

func (c *Compiler) Bytecode() *Bytecode{
	return &Bytecode{
		instructions: c.instructions,
		constants: c.constants,
	}
}


type Bytecode struct {
	instructions code.Instructions
	constants	[]object.Object
}

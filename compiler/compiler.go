package compiler

import (
	"fmt"
	"sort"

	"github.com/ShivankSharma070/go-compiler/ast"
	"github.com/ShivankSharma070/go-compiler/code"
	"github.com/ShivankSharma070/go-compiler/object"
)

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction // Why we need previous Instruction when we have last instructions ?? That because, when we remove last instruction, we need to keep track of the last instruction in stack
}

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable

	scope      []CompilationScope
	scopeIndex int
}

type EmittedInstruction struct {
	op       code.Opcode
	position int
}

func NewWithState(symTab *SymbolTable, cons []object.Object) *Compiler {
	comp := New()
	comp.constants = cons
	comp.symbolTable = symTab
	return comp
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	return &Compiler{
		constants:   []object.Object{},
		symbolTable: NewSymbolTable(),

		scope:      []CompilationScope{mainScope},
		scopeIndex: 0,
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.BlockStatement:
		for _, st := range node.Statements {
			err := c.Compile(st)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable: %s", node.Value)
		}

		c.emit(code.OpGetGlobal, symbol.Index)

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("Unkown operator: %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		idx := c.addConstant(integer)
		c.emit(code.OpConstant, idx)

	case *ast.StringLiteral:
		st := &object.String{Value: node.Value}
		idx := c.addConstant(st)
		c.emit(code.OpConstant, idx)

	case *ast.BoolExpression:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "-":
			c.emit(code.OpMinus)
		case "!":
			c.emit(code.OpBang)
		default:
			return fmt.Errorf("unkown operator: %s", node.Operator)
		}
	case *ast.IfElseExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value.
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		// Remove OpPop if its there
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		jumpPos := c.emit(code.OpJump, 999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			// If alternative is nill, insert a alternative block containing a instruction of generating null value
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))
	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}

		}

		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)
	case *ast.FunctionExpression:
		c.enterScope()
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		// Replace last instruction if opPop with a opReturnValue
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}

		instruction := c.leaveScope()
		compiledFunction := &object.CompiledFunction{Instructions: instruction}

		c.emit(code.OpConstant, c.addConstant(compiledFunction))
	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scope[c.scopeIndex].lastInstruction.op == op
}

func (c *Compiler) removeLastPop() {
	c.scope[c.scopeIndex].instructions = c.currentInstructions()[:c.scope[c.scopeIndex].lastInstruction.position]
	c.scope[c.scopeIndex].lastInstruction = c.scope[c.scopeIndex].previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := range len(newInstruction) {
		c.currentInstructions()[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scope[c.scopeIndex].lastInstruction
	last := EmittedInstruction{op: op, position: pos}

	c.scope[c.scopeIndex].lastInstruction = last
	c.scope[c.scopeIndex].previousInstruction = previous
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	c.scope[c.scopeIndex].instructions = append(c.currentInstructions(), ins...)
	return posNewInstruction
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scope[c.scopeIndex].instructions
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scope = append(c.scope, scope)
	c.scopeIndex++
}

func (c *Compiler) leaveScope() code.Instructions {
	instruction := c.currentInstructions()

	c.scope = c.scope[:len(c.scope)-1]
	c.scopeIndex--
	return instruction
}

func (c *Compiler) replaceLastPopWithReturn() {
	position := c.scope[c.scopeIndex].lastInstruction.position
	c.replaceInstruction(position, code.Make(code.OpReturnValue))
	c.scope[c.scopeIndex].lastInstruction.op = code.OpReturnValue
}

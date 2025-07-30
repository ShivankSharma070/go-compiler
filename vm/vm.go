package vm

import (
	"fmt"

	"github.com/ShivankSharma070/go-compiler/code"
	"github.com/ShivankSharma070/go-compiler/compiler"
	"github.com/ShivankSharma070/go-compiler/object"
)

const StackSize = 2048

type VM struct {
	instructions code.Instructions
	constants    []object.Object

	stack []object.Object
	sp    int //Stackpointer, points to next free slot.
}

func New(bc *compiler.Bytecode) *VM {
	return &VM{
		constants:    bc.Constants,
		instructions: bc.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd :
			right := vm.pop()
			left := vm.pop()
			rightValue := right.(*object.Integer).Value
			leftValue := left.(*object.Integer).Value
			vm.push(&object.Integer{Value: rightValue + leftValue})
		}
	}
	return nil
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("Stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp += 1

	return nil
}

func (vm *VM) pop()object.Object {
	o := vm.stack[vm.sp - 1]
	vm.sp--
	return o
}

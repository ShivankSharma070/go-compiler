package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	OpConstant Opcode = iota
)

type Instructions []byte

func (ins Instructions) String() string{
	var out bytes.Buffer
	i := 0
	for i < len(ins) {
		def, err := Lookup(Opcode(ins[i]))
		if err != nil{
			fmt.Fprintf(&out, "Error: %s\n", err)
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n",i, ins.fmtInstruction(def, operands))
		i += 1 + read 
	}

	return out.String()
}

// To display a single instruction
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string{
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
		case 1: 
			return fmt.Sprintf("%s %d", def.Name,operands[0])
	} 

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte

type Definition struct {
	Name string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant :{"OpConstant", []int{2}} ,
}

func Lookup(op Opcode) (*Definition, error) {
	def , ok := definitions[op]
	if !ok {
	return nil, fmt.Errorf("Opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands... int)Instructions{
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1 // 1 byte of opcode is already considered here
	for _, w := range def.OperandWidths{
		instructionLen += w
	}
	
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int,int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		{OpClosure, []int{65534,255}, []byte{byte(OpClosure), 255,254,255}},
	}
	for _, tt := range tests {
		instructions := Make(tt.op, tt.operands...)

		if len(instructions) != len(tt.expected) {
			t.Errorf("instructions has wrong lenght. want=%d, got=%d", len(tt.expected), len(instructions))
		}
		for i, b := range tt.expected {
			if instructions[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d, want=%d, got=%d", i, b, instructions[i])
			}
		}

	}
}

func TestInstructionString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
		Make(OpGetLocal, 5),
		Make(OpClosure, 65535, 255),
	}

	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpAdd
0010 OpGetLocal 5
0012 OpClosure 65535 255
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instruction wrongly formatted, want=%q, got=%q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		byteRead int
	}{
		{OpConstant, []int{65534}, 2},
		{OpGetLocal, []int{255}, 1},
		{OpClosure,[]int{65534, 255}, 3},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(tt.op)
		if err != nil {
			t.Errorf("definitation not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.byteRead {
			t.Errorf("n wrong, want=%d, got=%d", tt.byteRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong, want=%d, got=%d", want, operandsRead[i])
			}
		}

	}
}

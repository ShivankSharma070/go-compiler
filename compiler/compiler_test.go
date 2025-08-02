package compiler

import (
	"fmt"
	"testing"

	"github.com/ShivankSharma070/go-compiler/ast"
	"github.com/ShivankSharma070/go-compiler/code"
	"github.com/ShivankSharma070/go-compiler/lexer"
	"github.com/ShivankSharma070/go-compiler/object"
	"github.com/ShivankSharma070/go-compiler/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1+2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1-2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1*2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "2/1",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTest(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1<2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1>2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1==2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true==false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTest(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if(true){10}; 3333;
			`,
			expectedConstants: []any{10, 3333},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),              // 0000
				code.Make(code.OpJumpNotTruthy, 10), // 0001
				code.Make(code.OpConstant, 0),       // 0004
				code.Make(code.OpJump, 11),          // 0007
				code.Make(code.OpNull),              // 0010
				code.Make(code.OpPop),               // 0011
				code.Make(code.OpConstant, 1),       // 0012
				code.Make(code.OpPop),               // 0015
			},
		},
		{
			input: `
			if (true) {10} else {20}; 30;
			`,
			expectedConstants: []any{10, 20, 30},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),              //0000
				code.Make(code.OpJumpNotTruthy, 10), // 0001
				code.Make(code.OpConstant, 0),       // 0004
				code.Make(code.OpJump, 13),          // 0007
				code.Make(code.OpConstant, 1),       // 0010
				code.Make(code.OpPop),               // 0013
				code.Make(code.OpConstant, 2),       // 0014
				code.Make(code.OpPop),               //0017
			},
		},
	}

	runCompilerTest(t, tests)

}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: ` 
			let one = 1;
			let two = 2;
			`,
			expectedConstants : []any{1,2},
			expectedInstructions: []code.Instructions {
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input : `
			let one = 1;
			one;
			`,
			expectedConstants : []any{1},
			expectedInstructions : []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop), // because 'one;' is an expression
			},
		},
		{
			input : `
			let one = 1;
			let two = one;
			two;
			`,
			expectedConstants : []any{1},
			expectedInstructions : []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop), // because 'two;' is an expression
			},
		},
	}

	runCompilerTest(t, tests)
}

func runCompilerTest(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Errorf("Compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Errorf("testInstrucitons failed: %s", err)
		}
		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Errorf("testContants failed: %s", err)
		}
	}
}

func parse(input string) ast.Node {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testConstants(t *testing.T, expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("Wrong number of constants, want=%d, got=%d", len(expected), len(actual))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				t.Errorf("constant %d - testIntegerObject failed, %s", i, err)
			}
		}
	}

	return nil
}

func testInstructions(expectedInstructions []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expectedInstructions)
	if len(actual) != len(concatted) {
		return fmt.Errorf("Wrong instruction lenght, want=%q, got=%q", concatted, actual)
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("Wrong instruction at %d, want=%q, got=%q", i, ins, actual[i])
		}
	}

	return nil
}
func concatInstructions(instructions []code.Instructions) code.Instructions {
	var out code.Instructions
	for _, ins := range instructions {
		out = append(out, ins...)
	}

	return out
}

func testIntegerObject(value int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not integer, got=%T (%+v)", actual, actual)
	}

	if result.Value != value {
		return fmt.Errorf("object has wrong value, want=%d, got=%d", value, result.Value)
	}

	return nil
}

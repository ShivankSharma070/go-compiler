package vm

import (
	"fmt"
	"testing"

	"github.com/ShivankSharma070/go-compiler/ast"
	"github.com/ShivankSharma070/go-compiler/compiler"
	"github.com/ShivankSharma070/go-compiler/lexer"
	"github.com/ShivankSharma070/go-compiler/object"
	"github.com/ShivankSharma070/go-compiler/parser"
)

type vmTestCase struct {
	input    string
	expected any
}

func TestIntegerArigthmetic(t *testing.T){
		tests := []vmTestCase{
		{"1",1},
		{"2",2},
		{"1+2",3},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Errorf("Compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Errorf("vm error: %s", err)
		}

		stackElem := vm.StackTop()

		testExpectedObject(t, tt.expected, stackElem)
	}

}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerLiteral(actual, int64(expected))
		if err != nil {
			t.Errorf("testIntegerLiteral error: %s", err)
		}
	}
}
func testIntegerLiteral(obj object.Object, value int64) error {
	actual, ok := obj.(*object.Integer)
	if !ok {
		return fmt.Errorf("Object is not of type Integer, got=%T (%+v)", obj, obj)
	}

	if actual.Value != value {
		return fmt.Errorf("object has wrong value, want=%d, got=%d", value, actual.Value)
	}

	return nil
}

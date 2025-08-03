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

func TestIntegerArigthmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1+2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestStringExpression(t *testing.T){
	tests := []vmTestCase {
		{`"monkey"`, "monkey"},
		{`"mon"+"key"`, "monkey"},
		{`"mon"+"key"+"bannana"`, "monkeybannana"},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) {10}", Null},
		{"if (false) {10}", Null},
		{"!( if(false){10;} )", true},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}
	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase {
		{"let one = 1; one;", 1},
		{"let one = 1; let two=2; one+two;", 3},
		{"let one = 1; let two=one+one; one+two", 3},
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

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)
	}

}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected any, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerLiteral(actual, int64(expected))
		if err != nil {
			t.Errorf("testIntegerLiteral error: %s", err)
		}
	case bool:
		err := testBooleanLiteral(actual, expected)
		if err != nil {
			t.Errorf("testBooleanLiteral error: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not null: %T (%+v)", actual, actual)
		}
	}
}

func testBooleanLiteral(actual object.Object, value bool) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not boolean, got=%T (%+v)", actual, actual)
	}

	if result.Value != value {
		return fmt.Errorf("object has wrong value, want=%t, got=%t", value, result.Value)
	}
	return nil
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

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

func TestStringExpression(t *testing.T) {
	tests := []vmTestCase{
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
	tests := []vmTestCase{
		{"let one = 1; one;", 1},
		{"let one = 1; let two=2; one+two;", 3},
		{"let one = 1; let two=one+one; one+two", 3},
	}
	runVmTests(t, tests)
}

func TestArrayLiteral(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1,2,3]", []int{1, 2, 3}},
		{"[1+2,3-4,5*6]", []int{3, -1, 30}},
	}

	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{"{1:2, 2 :3}", map[object.HashKey]int64{
			(&object.Integer{Value: 1}).HashKey(): 2,
			(&object.Integer{Value: 2}).HashKey(): 3,
		}},
		{"{1+1 : 2*2, 3+3: 4*4}", map[object.HashKey]int64{
			(&object.Integer{Value: 2}).HashKey(): 4,
			(&object.Integer{Value: 6}).HashKey(): 16,
		}},
	}

	runVmTests(t, tests)
}

func TestIndexExpression(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{`fn(){5+10;}()`, 15},
		{`
			let fivePlusTen = fn() {5+10;}
			fivePlusTen()
			`, 15},
		{
			input: `
			let one = fn() { 1; };
			let two = fn() { 2; };
			one() + two()
			`,
			expected: 3,
		},
		{
			input: `
			let a = fn() { 1 };
			let b = fn() { a() + 1 };
			let c = fn() { b() + 1 };
			c();
			`,
			expected: 3,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionWithReturnStatements(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let earlyExit = fn() { return 99; 100; };
			earlyExit();
			`,
			expected: 99,
		},
		{
			input: `
			let earlyExit = fn() { return 99; return 100; };
			earlyExit();
			`,
			expected: 99,
		},
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
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object is not array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong num of elements: want=%d, got=%d", len(array.Elements), len(expected))
		}
		for i, expectedElement := range expected {
			err := testIntegerLiteral(array.Elements[i], int64(expectedElement))
			if err != nil {
				t.Errorf("testIntegerLiteral error: %s", err)
			}
		}

	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T", actual)
		}

		if len(hash.Pair) != len(expected) {
			t.Errorf("hash has wrong number of pairs, want=%d, got=%d", len(expected), len(hash.Pair))
			return
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pair[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in pairs. %d", expectedKey.Value)
			}

			err := testIntegerLiteral(pair.Value, expectedValue)
			if err != nil {
				t.Errorf("testIntegerLiteral err: %s", err)
			}

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

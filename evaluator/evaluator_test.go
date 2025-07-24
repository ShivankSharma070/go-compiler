package evaluator

import (
	"testing"

	"github.com/ShivankSharma070/go-interpreter/lexer"
	"github.com/ShivankSharma070/go-interpreter/object"
	"github.com/ShivankSharma070/go-interpreter/parser"
)

// ========= INTEGER ============
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		output := testEval(tt.input)
		testIntegerObject(t, output, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	integerObj, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("Object not of type object.Integer, got %T", obj)
		return false
	}

	if integerObj.Value != expected {
		t.Fatalf("integerObj.value is not %d, got %d", expected, integerObj.Value)
		return false
	}

	return true
}

// ======== BOOLEAN ==========
func TestEvalBooleanExpresion(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
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
	}

	for _, tt := range tests {
		output := testEval(tt.input)
		testBooleanObject(t, output, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("obj is not of type object.Boolean, got %T", obj)
		return false
	}

	if result.Value != expected {
		t.Fatalf("result.Value is not %t, got %t", expected, result.Value)
		return false
	}
	return true
}

// ========== PREFIX ==========
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		output := testEval(tt.input)
		testBooleanObject(t, output, tt.expected)
	}
}

// ================ Conditional Expressions ===========
func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Fatalf("object is not object.NULL, got %T (%+v)", obj, obj)
		return false
	}
	return true
}

// ================== RETURN ================
func TestReturnExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10
				}
				return 1
			}
			`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// =============== Let Statements ====================
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

// ==================== STRING ==================
func TestStringLiteral(t *testing.T) {
	input := `"hello world";`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Errorf("evaluated object is not of type object.String, got %T", evaluated)
	}

	if str.Value != "hello world" {
		t.Errorf("str.value is not %q, got %q", "hello world", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"hello"+"world"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("evaluated object is not of type object.String, got %T", evaluated)
	}

	if str.Value != "helloworld" {
		t.Errorf("str.value is not %q, got %q", "hello world", str.Value)
	}
}

// ==================== FUNCTION ==================
func TestFunctionExpression(t *testing.T) {
	input := ` fn (x) {x+2;}; `
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.FunctionLiteral)
	if !ok {
		t.Errorf("Evaluated value is not of type object.FunctionLiteral, got %T(%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("Not enough parameters, expected %d, got %d", 1, len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Errorf("parameter value is not %s, got %s", "x", fn.Parameters[0].String())
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Errorf("body vaulue is not %s, got %s", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClousers(t *testing.T) {
	input := `
	let adder = fn(x) {
	fn(y) { x+y; }
	}

	let newAdder = adder(2)
	newAdder(2)
	`

	testIntegerObject(t, testEval(input), 4)
}

// =============== ARRAY ====================
func TestArrayLiteral(t *testing.T) {
	input := "[1, 2*2, 3+3]"
	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Errorf("Evaluated object is not of type object.Array, got %T", evaluated)
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// =============== HASH ====================
func TestHashLiteral(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Errorf("evaluted object is not of type object.hash, got %T", evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pair) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs, got %d", len(result.Pair))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pair[expectedKey]
		if !ok {
			t.Errorf("no pair for give hash key")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{false:5}[false]`,
			5,
		},
		{
			`{true:5}[true]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// =============== Errors ====================
func TestBuiltInFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			err, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("evaluated object is not error, got %T", evaluated)
				continue
			}

			if err.Message != expected {
				t.Errorf("Error message is not %q, got %q", expected, err.Message)
			}
		}
	}

}

// =============== Errors ====================
func TestErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
		{
			`
if (10 > 1) {
if (10 > 1) {
return true + false;
}
return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar", "identifier not found: foobar",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errorObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No errror object retruned, got %T (%+v)", evaluated, evaluated)
			continue
		}

		if errorObj.Message != tt.expectedMessage {
			t.Errorf("ErrorObj.Message is not %s, got %s", tt.expectedMessage, errorObj.Message)
		}
	}

}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	env := object.NewEnvironment()

	program := p.ParseProgram()
	return Eval(program, env)
}

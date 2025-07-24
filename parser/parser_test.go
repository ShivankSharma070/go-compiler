package parser

import (
	"fmt"
	"testing"

	"github.com/ShivankSharma070/go-interpreter/ast"
	"github.com/ShivankSharma070/go-interpreter/lexer"
)

func TestReturnParser(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 1234;", 1234},
		{"return myvar;", "myvar"},
		{"return true;", true},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkForParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("Program does not container 1 statements. It contains %d", len(program.Statements))
			return
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not *ast.ReturnStatement, got %T", program.Statements[0])
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("TokenLiteral() is not 'return', got %s", returnStmt.TokenLiteral())
		}

		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)

	prog := p.ParseProgram()
	checkForParserErrors(t, p)

	if len(prog.Statements) != 1 {
		t.Fatalf("Now enough statement in program, got %d", len(prog.Statements))
	}

	for _, stmt := range prog.Statements {

		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not of *ast.ExpressionStatement, got %T", stmt)
		}

		if !testLiteralExpression(t, expStmt.Expression, "foobar") {
			t.FailNow()
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not have enough statements, got %d", len(program.Statements))
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	if !testLiteralExpression(t, expStmt.Expression, 5) {
		t.FailNow()
	}
}

func TestPrefixExpressionParsing(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)

		prog := parser.ParseProgram()

		checkForParserErrors(t, parser)
		if len(prog.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. Got %d", 1, len(prog.Statements))
		}

		expStmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("prog.Statement is not of type *ast.ExpressionStatement, got %T", prog.Statements[0])
		}

		prefixStmt, ok := expStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expStmt.Expression is not of type *ast.PrefixExpression, got %T", expStmt.Expression)
		}

		if prefixStmt.Operator != tt.operator {
			t.Fatalf("prefixStmt.Operator is not %s, got %s", tt.operator, prefixStmt.Operator)
		}

		if !testLiteralExpression(t, prefixStmt.Right, tt.integerValue) {
			t.FailNow()
		}
	}

}

func TestInfixExpressionParsing(t *testing.T) {
	tests := []struct {
		input    string
		left     any
		operator string
		right    any
	}{
		{"5+5", 5, "+", 5},
		{"5-5", 5, "-", 5},
		{"5*5", 5, "*", 5},
		{"5/5", 5, "/", 5},
		{"5>5", 5, ">", 5},
		{"5<5", 5, "<", 5},
		{"5==5", 5, "==", 5},
		{"5!=5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		prog := p.ParseProgram()
		checkForParserErrors(t, p)
		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain %d, got %d", 1, len(prog.Statements))
		}

		exp, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("prog.Statements is not ast.ExpressionStatement, got %T", prog.Statements[0])
		}

		if !testInfixExpression(t, exp.Expression, tt.left, tt.operator, tt.right) {
			t.FailNow()
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkForParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Fatalf("Expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBoolExpressionParsing(t *testing.T) {
	input := `false;`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program.Statements does not contains %d statements, got %d", 1, len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.Statement[0] is not ast.ExpressionStatemnet, got %T", program.Statements[0])
	}

	if !testLiteralExpression(t, exp.Expression, false) {
		t.FailNow()
	}
}

func TestIfStatementParsing(t *testing.T) {
	input := `if (x<y) {x};`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program does not container %d statements, got %d", 1, len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement is not of type ast.ExpressionStatement, got %T", program.Statements[0])
	}

	ifExp, ok := exp.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("exp.Expression is not of type ast.IfElseExpression, got %T", exp.Expression)
	}

	if !testInfixExpression(t, ifExp.Condition, "x", "<", "y") {
		t.FailNow()
		return
	}

	if len(ifExp.Consequence.Statements) != 1 {
		t.Fatalf("Consequences is not %d statements, got %d statements", 1, len(ifExp.Consequence.Statements))
	}

	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement, got %T", ifExp.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, consequence.Expression, "x") {
		t.FailNow()
	}

	if ifExp.Alternative != nil {
		t.Fatalf("Alternative is not nil, got %+v", ifExp.Alternative)
	}
}

func TestIfElseStatementParsing(t *testing.T) {
	input := `if (x<y) {x} else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program does not container %d statements, got %d", 1, len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement is not of type ast.ExpressionStatement, got %T", program.Statements[0])
	}

	ifExp, ok := exp.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("exp.Expression is not of type ast.IfElseExpression, got %T", exp.Expression)
	}

	if !testInfixExpression(t, ifExp.Condition, "x", "<", "y") {
		t.FailNow()
		return
	}

	if len(ifExp.Consequence.Statements) != 1 {
		t.Fatalf("Consequences is not %d statements, got %d statements", 1, len(ifExp.Consequence.Statements))
	}

	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement, got %T", ifExp.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, consequence.Expression, "x") {
		t.FailNow()
	}

	if len(ifExp.Alternative.Statements) != 1 {
		t.Fatalf("Alternative is not %d statements, got %d statements", 1, len(ifExp.Alternative.Statements))
	}

	alternative, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement[0] is not ast.ExpressionStatement, got %T", ifExp.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, alternative.Expression, "y") {
		t.FailNow()
	}
}

func TestFunctionStatementParsing(t *testing.T) {
	input := `fn(x, y) {x+y;}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not contain %d statements, got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Functionexpression, got %T", stmt.Expression)
	}

	if len(exp.Parameters) != 2 {
		t.Fatalf("Function does not container %d prameters, got %d", 2, len(exp.Parameters))
	}

	if !testLiteralExpression(t, exp.Parameters[0], "x") {
		return
	}
	if !testLiteralExpression(t, exp.Parameters[1], "y") {
		return
	}

	if len(exp.Body.Statements) != 1 {
		t.Fatalf("exp.Body.Statements does not contain %d statement, got %d", 1, len(exp.Body.Statements))
	}

	bodyStmt, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Body.Statements[0] is not ast.ExpressionStatement, got %T", exp.Body.Statements[0])
	}

	if !testInfixExpression(t, bodyStmt.Expression, "x", "+", "y") {
		return
	}
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn(){}", []string{}},
		{"fn(x){}", []string{"x"}},
		{"fn(x,y){}", []string{"x", "y"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkForParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		functionStmt := stmt.Expression.(*ast.FunctionExpression)
		if len(functionStmt.Parameters) != len(tt.expectedParams) {
			t.Fatalf("length of parameters is wrong, got %d", len(functionStmt.Parameters))
		}

		for i, iden := range tt.expectedParams {
			testLiteralExpression(t, functionStmt.Parameters[i], iden)
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp is not of type ast.StringLiteral, got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value is not %q, got %q", "hello world", literal.Value)
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2*3, 4+5)`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain %d statements, got %d", 1, len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	callExp, ok := exp.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf(" exp.Expression is not ast.CallExpression, got %T", exp.Expression)
	}

	if !testLiteralExpression(t, callExp.Function, "add") {
		return
	}

	if len(callExp.Argument) != 3 {
		t.Fatalf("Arguments are not %d, got %d", 3, len(callExp.Argument))
	}

	testLiteralExpression(t, callExp.Argument[0], 1)
	testInfixExpression(t, callExp.Argument[1], 2, "*", 3)
	testInfixExpression(t, callExp.Argument[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkForParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Argument) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Argument))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Argument[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Argument[i].String())
			}
		}
	}
}

func TestLetParser(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5", "x", 5},
		{"let y = true", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		p := New(lex)

		prog := p.ParseProgram()
		checkForParserErrors(t, p)
		if len(prog.Statements) != 1 {
			t.Fatalf("Returned Program does not container 1 statements got %d", len(prog.Statements))
		}
		stmt := prog.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		value := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, tt.expectedValue) {
			return
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2*3, 4+5]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkForParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 3)
	testInfixExpression(t, array.Elements[2], 4, "+", 5)
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1+1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkForParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Errorf("stmt.Expression is not ast.IndexExpression, got %T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiterals(t *testing.T){
	input := `{"one":1, "two":2, "three":3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkForParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash , ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not of type ast.HashLiteral, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("lenght os hash.Pairs is not %d, got %d",3, len(hash.Pairs))
	}
	
	expected := map[string]int64{
		"one": 1,
		"two": 2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key in hash.Pairs is not of type string, got %T", key)
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T){
	input := `{}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkForParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash , ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not of type ast.HashLiteral, got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("lenght os hash.Pairs is not %d, got %d",0, len(hash.Pairs))
	}
}

func TestParsingHashliteralWithExpression(t *testing.T) {
	input := `{"one": 0+1, "two": 10 - 3, "three": 15/5}`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkForParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash , ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not of type ast.HashLiteral, got %T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("lenght os hash.Pairs is not %d, got %d", 3, len(hash.Pairs))
	}
	
	tests := map[string]func(ast.Expression) {
		"one" : func(e ast.Expression) {
			testInfixExpression(t, e,0, "+", 1)
		},
		"two" : func(e ast.Expression) {
			testInfixExpression(t, e,10, "-", 3)
		},
		"three" : func(e ast.Expression) {
			testInfixExpression(t, e,15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key in hash.Pairs is not of type string, got %T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found.", literal.String())
		}

		testFunc(value)
	}
}

// ========== HELPER FUNCTIONS ================

func testIdentifier(t *testing.T, expStmt ast.Expression, value string) bool {
	idenStmt, ok := expStmt.(*ast.Identifier)
	if !ok {
		t.Fatalf("expStmt.Expression is not of *ast.Identifier, got %T", expStmt)
		return false
	}

	if idenStmt.Value != value {
		t.Fatalf("idenStmt.value is not equal to %s, got %s", value, idenStmt.Value)
		return false
	}

	if idenStmt.TokenLiteral() != value {
		t.Fatalf("idenStmt.tokenLiteral() is not equal to %s, got %s", value, idenStmt.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	intStmt, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Exp is not of type ast.IntegerLiteral, got %T", exp)
		return false
	}

	if intStmt.Value != value {
		t.Fatalf("intStmt.Value is not %d, got %d", value, intStmt.Value)
		return false
	}

	if intStmt.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("intStmt.TokenLiteral() is not equal to %s, got %s", fmt.Sprintf("%d", value), intStmt.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, value any) bool {
	switch v := value.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolExpression(t, exp, v)
	}

	t.Fatalf("type of exp not handled, got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("prog.Statements is not ast.ExpressionStatement, got %T", exp)
		return false
	}

	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}

	if infixExp.Operator != operator {
		t.Fatalf("infixExp.Operator is not %s, got %s", operator, infixExp.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}

func testBoolExpression(t *testing.T, exp ast.Expression, value bool) bool {
	boolExp, ok := exp.(*ast.BoolExpression)
	if !ok {
		t.Fatalf("exp.Expression is not ast.BoolExpression, got %T", exp)
		return false
	}

	if boolExp.String() != fmt.Sprint(value) {
		t.Fatalf("boolExp.String() does not return %s, got %s", fmt.Sprint(value), boolExp.String())
		return false
	}

	if boolExp.Value != value {
		t.Fatalf("boolExp.Value is not %t, got %t", value, boolExp.Value)
		return false
	}
	return true
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	// Check if it is let statement
	if stmt.TokenLiteral() != "let" {
		t.Errorf("Tokenliteral is not let, got %s\n", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement is not ast.LetStatement, got %T\n", stmt)
		return false
	}

	// Check if the identifier is what we expect
	if letStmt.Name.Value != name {
		t.Errorf("letstmt.name.value is not same as name, got %s", letStmt.Name.Value)
		return false
	}

	return true
}

func checkForParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	for _, err := range errors {
		t.Errorf("Parser error %s", err)

	}
	t.FailNow()
}

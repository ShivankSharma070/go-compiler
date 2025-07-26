package lexer

import (
	"testing"

	"github.com/ShivankSharma070/go-compiler/token"
)

func TestLexer(t *testing.T) {
	input := `let five = 5;
	let ten = 10;

	let sum = fn (a, b) {
		a + b;
	};

	let result = sum(five,ten);
	!-/*;
	< >;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	==
	!=
	"foobar"
	"foo bar"
	[1,2];
	{"foo":"bar"};
	`

	test := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDEN, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDEN, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDEN, "sum"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDEN, "a"},
		{token.COMMA, ","},
		{token.IDEN, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDEN, "a"},
		{token.PLUS, "+"},
		{token.IDEN, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDEN, "result"},
		{token.ASSIGN, "="},
		{token.IDEN, "sum"},
		{token.LPAREN, "("},
		{token.IDEN, "five"},
		{token.COMMA, ","},
		{token.IDEN, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERIK, "*"},
		{token.SEMICOLON, ";"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EQ, "=="},
		{token.NOT_EQ, "!="},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range test {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("Test_%d: Type mismatch Expected:%q Got:%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("Test_%d: Literal mismatch Expected:%q Got:%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

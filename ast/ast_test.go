package ast

import (
	"testing"

	"github.com/ShivankSharma070/go-compiler/token"
)

// let myvar = anotherVar;
func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDEN, Literal: "myvar"},
					Value: "myvar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDEN, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myvar = anotherVar;" {
		t.Errorf("program.String does not return let myvar = anotherVar, got %s", program.String())
	}

}

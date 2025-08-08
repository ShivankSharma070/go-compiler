package evaluator

import (
	"github.com/ShivankSharma070/go-compiler/object"
)

var builtins = map[string]*object.Builtin{
	"len":object.GetBuiltinByName("len") ,
	"exit": object.GetBuiltinByName("exit"),
	"first": object.GetBuiltinByName("first"),
	"last": object.GetBuiltinByName("last"),
	"rest": object.GetBuiltinByName("rest"),
	"push": object.GetBuiltinByName("push"),
	"puts":object.GetBuiltinByName("puts") ,
}

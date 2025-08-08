package object

import (
	"fmt"
	"os"
)

// ================== BUILT-IN FUNCTION ===================

type BuiltInFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltInFunction
}

func (bu *Builtin) Inspect() string  { return "builtin function" }
func (bu *Builtin) Type() ObjectType { return BUILTIN_OBJ }

var Builtins = []struct {
	Name    string
	Buitlin *Builtin
}{
	{
		"len",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		"puts",
		&Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return nil
			},
		},
	},
	{
		"first",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return nil
			},
		},
	},
	{
		"last",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[len(arr.Elements)-1]
				}
				return nil
			},
		},
	},
	{
		"rest",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					elements := make([]Object, length-1, length-1)
					copy(elements, arr.Elements[1:])
					return &Array{Elements: elements}
				}

				return nil
			},
		},
	},
	{
		"push",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)
				elements := make([]Object, length+1, length+1)
				copy(elements, arr.Elements)
				elements[length] = args[1]
				return &Array{Elements: elements}
			},
		},
	},
	{
		"exit",
		&Builtin{
			Fn: func(args ...Object) Object {
				os.Exit(1)
				return nil
			},
		},
	},
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Buitlin
		}
	}

	return nil
}

func newError(format string, args ...any) Object {
	return &Error{
		Message: fmt.Sprintf(format, args...),
	}
}

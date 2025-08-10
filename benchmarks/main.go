package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ShivankSharma070/go-compiler/compiler"
	"github.com/ShivankSharma070/go-compiler/evaluator"
	"github.com/ShivankSharma070/go-compiler/lexer"
	"github.com/ShivankSharma070/go-compiler/object"
	"github.com/ShivankSharma070/go-compiler/parser"
	"github.com/ShivankSharma070/go-compiler/vm"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

var input = `
	let fibonacci = fn(x) {
		if (x == 0) {
			0
		} else {
			if (x == 1) {
				return 1;
			} else {
				fibonacci(x - 1) + fibonacci(x - 2);
			}
		}
	};
	fibonacci(35);
`

func main() {
	flag.Parse()
	var duration time.Duration
	var result object.Object

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return 
		}

		machine := vm.New(comp.Bytecode())
		start := time.Now()
		err = machine.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)	
			return
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else {
		environment := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, environment)
		duration = time.Since(start)
	}

	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		*engine,
		result.Inspect(),
		duration,
	)

}

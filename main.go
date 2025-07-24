package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/ShivankSharma070/go-interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Welcome to Monkey Programming Language REPL", user.Username)
	fmt.Println("Feel free to type any command..")
	repl.Start(os.Stdin, os.Stdout)
}

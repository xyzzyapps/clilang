package main

import (
	"clilang/interpreter"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: clilang do=<command> arg=<value> ...")
		os.Exit(1)
	}

	interp := interpreter.NewInterpreter()
	
	// Register all core built-ins
	interpreter.RegisterBuiltins(interp)

	// Register custom extensions
	RegisterCustomCommands(interp)

	// Initialize state
	interp.Tokens = args

	// Evaluate
	res, err := interp.Eval()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print final result if it's not empty and no error occurred
	if res != "" {
		fmt.Println(res)
	}
}

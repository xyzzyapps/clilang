package main

import (
	"clilang/interpreter"
	"fmt"
	"os/exec"
	"strings"
)

// RegisterCustomCommands demonstrates how to extend clilang with Go functions.
func RegisterCustomCommands(i *interpreter.Interpreter) {
	// Example 1: String reversal
	i.Register("reverse", 1, func(interp *interpreter.Interpreter, args []string) (string, error) {
		runes := []rune(args[0])
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), nil
	})

	// Example 2: Execute system command
	i.Register("sys", 2, func(interp *interpreter.Interpreter, args []string) (string, error) {
		cmdName := args[0]
		cmdArgs := strings.Fields(args[1]) // simple splitting by space
		
		cmd := exec.Command(cmdName, cmdArgs...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("sys command failed: %w", err)
		}
		return string(out), nil
	})
}

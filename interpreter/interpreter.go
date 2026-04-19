package interpreter

import (
	"errors"
	"fmt"
	"strings"
)

// CommandFunc defines the signature for a command handler.
type CommandFunc func(i *Interpreter, args []string) (string, error)

// Command defines a command's arity and function.
type Command struct {
	Arity int
	Fn    CommandFunc
}

// Interpreter holds the evaluation state.
type Interpreter struct {
	Vars     map[string]string
	Tokens   []string
	Pos      int
	Commands map[string]Command
}

// NewInterpreter creates a new Interpreter instance.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		Vars:     make(map[string]string),
		Commands: make(map[string]Command),
	}
}

// Register adds a new command to the interpreter.
func (i *Interpreter) Register(name string, arity int, fn CommandFunc) {
	i.Commands[name] = Command{
		Arity: arity,
		Fn:    fn,
	}
}

// SubInterpreter creates a new interpreter for a block of tokens, sharing variables and commands.
func (i *Interpreter) SubInterpreter(blockArg string) *Interpreter {
	var tokens []string
	if blockArg != "" {
		tokens = strings.Split(blockArg, "\x00")
	}
	return &Interpreter{
		Vars:     i.Vars, // share global scope
		Tokens:   tokens,
		Pos:      0,
		Commands: i.Commands,
	}
}

// NextArgument evaluates and returns the next argument.
func (i *Interpreter) NextArgument() (string, error) {
	if i.Pos >= len(i.Tokens) {
		return "", errors.New("unexpected end of arguments")
	}

	token := i.Tokens[i.Pos]
	i.Pos++

	parts := strings.SplitN(token, "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format, expected key=value: %s", token)
	}

	key, val := parts[0], parts[1]

	switch key {
	case "arg":
		return val, nil
	case "var":
		v, ok := i.Vars[val]
		if !ok {
			return "", fmt.Errorf("undefined variable: %s", val)
		}
		return v, nil
	case "eval":
		if val != "begin" {
			return "", fmt.Errorf("expected eval=begin, got: %s", token)
		}
		return i.EvalUntil("eval=end")
	case "block":
		if val != "begin" {
			return "", fmt.Errorf("expected block=begin, got: %s", token)
		}
		blockTokens, err := i.collectBlock()
		if err != nil {
			return "", err
		}
		return strings.Join(blockTokens, "\x00"), nil
	default:
		return "", fmt.Errorf("unexpected token type in argument position: %s", key)
	}
}

func (i *Interpreter) collectBlock() ([]string, error) {
	var block []string
	depth := 1
	for i.Pos < len(i.Tokens) {
		token := i.Tokens[i.Pos]
		i.Pos++

		parts := strings.SplitN(token, "=", 2)
		if len(parts) == 2 && parts[0] == "block" {
			if parts[1] == "begin" {
				depth++
			} else if parts[1] == "end" {
				depth--
				if depth == 0 {
					return block, nil
				}
			}
		}
		block = append(block, token)
	}
	return nil, errors.New("unclosed block")
}

// EvalUntil evaluates commands until the specified endToken is found.
// If endToken is empty, it evaluates until EOF.
func (i *Interpreter) EvalUntil(endToken string) (string, error) {
	var lastResult string
	for i.Pos < len(i.Tokens) {
		token := i.Tokens[i.Pos]
		if endToken != "" && token == endToken {
			i.Pos++
			return lastResult, nil
		}

		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 || parts[0] != "do" {
			return "", fmt.Errorf("expected do=<cmd>, got: %s", token)
		}
		i.Pos++

		cmdName := parts[1]
		cmd, ok := i.Commands[cmdName]
		if !ok {
			return "", fmt.Errorf("unknown command: %s", cmdName)
		}

		var args []string
		for j := 0; j < cmd.Arity; j++ {
			arg, err := i.NextArgument()
			if err != nil {
				return "", fmt.Errorf("error reading argument %d for command %s: %w", j+1, cmdName, err)
			}
			args = append(args, arg)
		}

		res, err := cmd.Fn(i, args)
		if err != nil {
			return "", fmt.Errorf("error executing %s: %w", cmdName, err)
		}
		lastResult = res
	}

	if endToken != "" {
		return "", fmt.Errorf("unexpected EOF, looking for: %s", endToken)
	}
	return lastResult, nil
}

// Eval evaluates all tokens.
func (i *Interpreter) Eval() (string, error) {
	return i.EvalUntil("")
}

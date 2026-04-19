package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

// RegisterBuiltins registers all built-in commands into the interpreter.
func RegisterBuiltins(i *Interpreter) {
	i.Register("set", 2, cmdSet)
	i.Register("puts", 1, cmdPuts)

	i.Register("add", 2, cmdMath(func(a, b int) int { return a + b }))
	i.Register("sub", 2, cmdMath(func(a, b int) int { return a - b }))
	i.Register("mul", 2, cmdMath(func(a, b int) int { return a * b }))
	i.Register("div", 2, cmdMath(func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	}))

	i.Register("eq", 2, cmdLogic(func(a, b string) bool { return a == b }))
	i.Register("neq", 2, cmdLogic(func(a, b string) bool { return a != b }))
	i.Register("lt", 2, cmdLogicInt(func(a, b int) bool { return a < b }))
	i.Register("gt", 2, cmdLogicInt(func(a, b int) bool { return a > b }))

	i.Register("if", 2, cmdIf)
	i.Register("ifelse", 3, cmdIfElse)
	i.Register("while", 2, cmdWhile)

	i.Register("def", 3, cmdDef)
}

func cmdSet(i *Interpreter, args []string) (string, error) {
	i.Vars[args[0]] = args[1]
	return args[1], nil
}

func cmdPuts(i *Interpreter, args []string) (string, error) {
	fmt.Println(args[0])
	return "", nil
}

func cmdMath(op func(int, int) int) CommandFunc {
	return func(i *Interpreter, args []string) (string, error) {
		a, err := strconv.Atoi(args[0])
		if err != nil {
			return "", fmt.Errorf("invalid math argument: %s", args[0])
		}
		b, err := strconv.Atoi(args[1])
		if err != nil {
			return "", fmt.Errorf("invalid math argument: %s", args[1])
		}
		return strconv.Itoa(op(a, b)), nil
	}
}

func cmdLogic(op func(string, string) bool) CommandFunc {
	return func(i *Interpreter, args []string) (string, error) {
		if op(args[0], args[1]) {
			return "1", nil
		}
		return "0", nil
	}
}

func cmdLogicInt(op func(int, int) bool) CommandFunc {
	return func(i *Interpreter, args []string) (string, error) {
		a, err := strconv.Atoi(args[0])
		if err != nil {
			return "", fmt.Errorf("invalid logic argument: %s", args[0])
		}
		b, err := strconv.Atoi(args[1])
		if err != nil {
			return "", fmt.Errorf("invalid logic argument: %s", args[1])
		}
		if op(a, b) {
			return "1", nil
		}
		return "0", nil
	}
}

func isTruthy(s string) bool {
	return s != "0" && s != "" && s != "false"
}

func cmdIf(i *Interpreter, args []string) (string, error) {
	condBlock := args[0]
	bodyBlock := args[1]

	condInterpreter := i.SubInterpreter(condBlock)
	condRes, err := condInterpreter.Eval()
	if err != nil {
		return "", err
	}

	if isTruthy(condRes) {
		bodyInterpreter := i.SubInterpreter(bodyBlock)
		return bodyInterpreter.Eval()
	}
	return "", nil
}

func cmdIfElse(i *Interpreter, args []string) (string, error) {
	condBlock := args[0]
	trueBlock := args[1]
	falseBlock := args[2]

	condInterpreter := i.SubInterpreter(condBlock)
	condRes, err := condInterpreter.Eval()
	if err != nil {
		return "", err
	}

	var targetBlock string
	if isTruthy(condRes) {
		targetBlock = trueBlock
	} else {
		targetBlock = falseBlock
	}

	bodyInterpreter := i.SubInterpreter(targetBlock)
	return bodyInterpreter.Eval()
}

func cmdWhile(i *Interpreter, args []string) (string, error) {
	condBlock := args[0]
	bodyBlock := args[1]

	var lastRes string
	for {
		condInterpreter := i.SubInterpreter(condBlock)
		condRes, err := condInterpreter.Eval()
		if err != nil {
			return "", err
		}

		if !isTruthy(condRes) {
			break
		}

		bodyInterpreter := i.SubInterpreter(bodyBlock)
		res, err := bodyInterpreter.Eval()
		if err != nil {
			return "", err
		}
		lastRes = res
	}

	return lastRes, nil
}

func cmdDef(i *Interpreter, args []string) (string, error) {
	name := args[0]
	paramsBlock := args[1]
	bodyBlock := args[2]

	var paramNames []string
	if paramsBlock != "" {
		for _, token := range strings.Split(paramsBlock, "\x00") {
			parts := strings.SplitN(token, "=", 2)
			if len(parts) == 2 && parts[0] == "arg" {
				paramNames = append(paramNames, parts[1])
			} else {
				return "", fmt.Errorf("invalid token in params block, expected arg=<name>: %s", token)
			}
		}
	}

	i.Register(name, len(paramNames), func(callerInterpreter *Interpreter, callerArgs []string) (string, error) {
		bodyInterpreter := callerInterpreter.SubInterpreter(bodyBlock)
		
		// create local scope but since Variables are passed by map ref, we need to create a new map
		// wait, if we copy Vars, we have closures? TCL doesn't have true lexical closures by default, usually explicit.
		// we'll copy the global variables into the new scope.
		newVars := make(map[string]string)
		for k, v := range callerInterpreter.Vars {
			newVars[k] = v
		}

		// bind arguments
		for idx, pName := range paramNames {
			newVars[pName] = callerArgs[idx]
		}
		
		bodyInterpreter.Vars = newVars

		return bodyInterpreter.Eval()
	})

	return "", nil
}

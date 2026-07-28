# Clilang

Clilang is a TCL-like programming language implemented in Go, specifically designed to be written entirely as command-line arguments without running into shell conflict issues.

It achieves this by using a strict `key=value` syntax (inspired by the Linux `dd` command), completely eliminating the need for brackets, braces, and quotes that normally confuse shells like `bash` or `PowerShell`.

## Features
- **Conflict-Free Syntax**: No shell metacharacters needed.
- **TCL-like capabilities**: Supports sequential evaluation, variables, deferred blocks, and command substitution.
- **Extensible**: Easily add your own commands via Go functions.

## Syntax Rules
Every argument must be in the format `key=value`.
- `do=<command>`: Start a command execution.
- `arg=<value>`: Pass a literal string argument to the command.
- `var=<name>`: Automatically looks up a variable and passes its value as an argument.
- `eval=begin` / `eval=end`: Evaluates the inner sequence of commands and returns the final result as a single argument.
- `block=begin` / `block=end`: Defers evaluation, passing the sequence of tokens to control flow structures (like `if` or `while`).

## Built-in Commands
- `set`, `puts`
- Math: `add`, `sub`, `mul`, `div`
- Logic: `eq`, `neq`, `lt`, `gt`
- Control Flow: `if`, `ifelse`, `while`
- Functions: `def`

## Example Usage

**Reversing a string**
```bash
./clilang do=set arg=msg arg="Hello World" do=puts eval=begin do=reverse var=msg eval=end
```

**Looping**
*(equivalent to `set i 0; while {$i < 3} { puts $i; set i [add $i 1] }`)*
```bash
./clilang \
  do=set arg=i arg=0 \
  do=while \
    block=begin do=lt var=i arg=3 block=end \
    block=begin \
      do=puts var=i \
      do=set arg=i eval=begin do=add var=i arg=1 eval=end \
    block=end
```

## Adding Custom Commands
Open `custom_commands.go` and use the `Interpreter.Register` method to bind your own logic. See the provided `reverse` and `sys` commands for examples of handling arguments and returning values.

## Building
```bash
go build -o clilang.exe .
```
## License

MIT
 
## Signature

Original Research by Xyzzy, built with assistance from **Qwen 3.5**.   

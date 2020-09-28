package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// parser parses an .vm file and provides convenient access to the command's arguments.
// In addition, removes all white space and comments.
type parser struct {
	scanner *bufio.Scanner
	command []string
}

// newParser creates new parser.
func newParser(file io.Reader) (*parser, error) {
	return &parser{bufio.NewScanner(file), []string{}}, nil
}

// hasMoreCommands returns true if there are more commands in the input.
func (p *parser) hasMoreCommands() bool { return p.scanner.Scan() }

// advance reads the next command from the input and makes it the current command.
// Should be called only if HasMoreCommands() is true.
// Initially there is no current command.
func (p *parser) advance() {
	command := p.scanner.Text()

	for ignoreCommand(command) && p.hasMoreCommands() {
		command = p.scanner.Text()
	}

	p.command = strings.Fields(command)

	// Remove comments after the command
	for i, field := range p.command {
		if field == "//" {
			p.command = p.command[:i]
			return
		}
	}
}

// ignoreCommand defines which types of commands to ignore
func ignoreCommand(command string) bool {
	return command == "" || strings.HasPrefix(command, "//") || len(strings.Fields(command)) == 0
}

// commandType return one of defined command type constants
func (p *parser) commandType() commandType {
	switch p.command[0] {
	case "push":
		return pushCmd
	case "pop":
		return popCmd
	case "label":
		return labelCmd
	case "goto":
		return gotoCmd
	case "if-goto":
		return ifCmd
	case "function":
		return functionCmd
	case "call":
		return callCmd
	case "return":
		return returnCmd
	default:
		return arithmeticCmd
	}
}

// firstArgument returns the first argument of the current command.
// In the case of Arithmetic the command itself (add, sub, etc.) is returned.
// Should not be called if the current command is Return.
func (p *parser) firstArgument() string {
	if len(p.command) == 1 {
		return p.command[0]
	}

	return p.command[1]
}

// secondArgument returns the second argument of the command.
// Should be called only if the current command is Push, Pop, Function, or Call.
func (p *parser) secondArgument() int {
	i, err := strconv.ParseUint(p.command[2], 10, 15)
	if err != nil {
		return 0
	}

	return int(i)
}

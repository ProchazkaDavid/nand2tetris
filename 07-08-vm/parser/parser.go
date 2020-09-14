package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/ProchazkaDavid/nand2tetris/07-08-vm/command"
)

// Parser parses an .vm file and provides convenient access to the command's arguments.
// In addition, removes all white space and comments.
type Parser struct {
	scanner *bufio.Scanner
	command []string
}

// New opens the input file and gets ready to parse it.
func New(path string) (*Parser, error) {
	f, err := os.Open(path)
	return &Parser{bufio.NewScanner(f), []string{}}, err
}

// HasMoreCommands returns true if there are more commands in the input.
func (p *Parser) HasMoreCommands() bool { return p.scanner.Scan() }

// Advance reads the next command from the input and makes it the current command.
// Should be called only if HasMoreCommands() is true.
// Initially there is no current command.
func (p *Parser) Advance() {
	command := p.scanner.Text()

	for ignoreCommand(command) && p.HasMoreCommands() {
		command = p.scanner.Text()
	}

	p.command = strings.Fields(command)
}

// ignoreCommand defines which types of commands to ignore
func ignoreCommand(command string) bool {
	return command == "" || strings.HasPrefix(command, "//") || len(strings.Fields(command)) == 0
}

// CommandType return one of defined command type constants
func (p *Parser) CommandType() command.Type {
	if len(p.command) == 1 {
		return command.Arithmetic
	}

	// Push or Pop
	if p.command[0] == "push" {
		return command.Push
	}
	return command.Pop
}

// Arg1 returns the first argument of the current command.
// In the case of Arithmetic the command itself (add, sub, etc.) is returned.
// Should not be called if the current command is Return.
func (p *Parser) Arg1() string {
	if len(p.command) == 1 {
		return p.command[0]
	}

	return p.command[1]
}

// Arg2 returns the second argument of the command.
// Should be called only if the current command is Push, Pop, Function, or Call.
func (p *Parser) Arg2() int {
	i, err := strconv.ParseUint(p.command[2], 10, 15)
	if err != nil {
		return 0
	}

	return int(i)
}

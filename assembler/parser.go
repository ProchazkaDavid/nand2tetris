package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CommandType represents command type
type CommandType int

const (
	// ACommand is an addressing instruction
	//   Format: @value
	// Where value is either a non-negative decimal number
	// or a symbol referring to such number.
	ACommand CommandType = iota
	// CCommand is a compute instruction
	//   Format: dest=comp;jump
	// Either the dest or jump fields may be empty.
	// If dest is empty, the '=' is omitted;
	// If jump is empty, the ';' is omitted.
	CCommand
	// LCommand pseudo-command binds the Symbol to the memory location into which
	// the next command in the program will be stored.
	//   Format: (Symbol)
	LCommand
)

// firstRAMAddress represents first available RAM address that is used to store variables
const firstRAMAddress = 0x10

// Parser reads an assembly language command, parses it,
// and provides convenient access to the command's components (fields and symbols).
// In addition, removes all white space and comments.
type Parser struct {
	file       *os.File
	scanner    *bufio.Scanner
	command    string
	ramAddress uint16
}

// newParser opens the input file and gets ready to parse it.
func newParser(filename string) (*Parser, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("parser can't open the file: %w", err)
	}

	return &Parser{f, bufio.NewScanner(f), "", firstRAMAddress}, nil
}

// HasMoreCommands returns true if there are more commands in the input.
func (p *Parser) HasMoreCommands() bool { return p.scanner.Scan() }

// Advance reads the next command from the input and makes it the current command.
// Should be called only if HasMoreCommands() is true.
// Initially there is no current command.
func (p *Parser) Advance() {
	p.command = p.scanner.Text()

	for ignoreCommand(p.command) && p.HasMoreCommands() {
		p.command = p.scanner.Text()
	}

	// Removes whitespaces before the command and whitespaces and comments after the command
	p.command = strings.Fields(p.command)[0]
}

// ignoreCommand defines which types of commands to ignore
func ignoreCommand(command string) bool {
	return command == "" || strings.HasPrefix(command, "//") || len(strings.Fields(command)) == 0
}

// commandType returns the type of the current command:
//   ACommand for @Xxx where Xxx is either a symbol or a decimal number
//   CCommand for dest=comp;jump
//   LCommand (actually, pseudo-command) for (Xxx) where Xxx is a symbol.
func (p *Parser) commandType() CommandType {
	switch {
	case strings.HasPrefix(p.command, "@"):
		return ACommand
	case strings.HasPrefix(p.command, "("):
		return LCommand
	default:
		return CCommand
	}
}

// symbol returns the symbol or decimal Xxx of the current command @Xxx or (Xxx).
// Should be called only when CommandType() is ACommand or LCommand.
func (p *Parser) symbol() string {
	if strings.HasPrefix(p.command, "@") {
		return p.command[1:]
	}

	return p.command[1 : len(p.command)-1]
}

// dest returns the dest mnemonic in the current C-command (8 possibilities).
// Should be called only when CommandType() is CCommand.
func (p *Parser) dest() string {
	if i := strings.IndexRune(p.command, '='); i != -1 {
		return p.command[:i]
	}

	return ""
}

// comp returns the comp mnemonic in the current C-command (28 possibilities).
// Should be called only when CommandType() is CCommand.
func (p *Parser) comp() string {
	from, to := 0, len(p.command)

	if i := strings.IndexRune(p.command, '='); i != -1 {
		from = i + 1
	}
	if i := strings.IndexRune(p.command, ';'); i != -1 {
		to = i
	}

	return p.command[from:to]
}

// jump returns the jump mnemonic in the current C-command (8 possibilities).
// Should be called only when CommandType() is CCommand.
func (p *Parser) jump() string {
	if i := strings.IndexRune(p.command, ';'); i != -1 {
		return p.command[i+1:]
	}

	return ""
}

// getFreeRAMAddress return next free RAM address which is used for storing variables
func (p *Parser) getFreeRAMAddress() uint16 {
	address := p.ramAddress
	p.ramAddress++
	return address
}

// parseSymbols returns populated symboltable.SymbolTable with parser.LCommands
// and prepares parser for another file scan.
func (p *Parser) parseSymbols() (symbolTable, error) {
	table := newSymbolTable()
	line := uint16(0)

	for p.HasMoreCommands() {
		p.Advance()
		switch p.commandType() {
		case ACommand, CCommand:
			line++
		case LCommand:
			table[p.symbol()] = line
		}
	}

	if _, err := p.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("parser: can't seek to the beginning of the file: %w", err)
	}

	p.scanner = bufio.NewScanner(p.file)

	return table, nil
}

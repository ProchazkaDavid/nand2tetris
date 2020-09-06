package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ProchazkaDavid/nand2tetris/06-assembler/code"
	"github.com/ProchazkaDavid/nand2tetris/06-assembler/parser"
	"github.com/ProchazkaDavid/nand2tetris/06-assembler/symboltable"
)

// writeACommand writes parser's current ACommand address value to the f based on the table
func writeACommand(f *os.File, p *parser.Parser, table symboltable.SymbolTable) error {
	symbol := p.Symbol()

	// Decimal address
	if '0' <= symbol[0] && symbol[0] <= '9' {
		address, err := strconv.ParseUint(symbol, 10, 15)
		if err != nil {
			return err
		}

		return writeAddress(f, uint16(address))
	}

	// Known symbol or variable in SymbolTable
	if address, ok := table[symbol]; ok {
		return writeAddress(f, address)
	}

	// Unknown symbol, declaration of a new variable
	table[symbol] = p.GetFreeRAMAddress()
	return writeAddress(f, table[symbol])
}

// writeAddress writes address to the f
func writeAddress(f *os.File, address uint16) error {
	if _, err := fmt.Fprintf(f, "0%015b\n", address); err != nil {
		return err
	}
	return nil
}

// writeCCommand writes parser's current CCommand to the f
func writeCCommand(f *os.File, p *parser.Parser) error {
	if _, err := fmt.Fprintf(f, "111%s%s%s\n", code.Comp(p.Comp()), code.Dest(p.Dest()), code.Jump(p.Jump())); err != nil {
		return err
	}
	return nil
}

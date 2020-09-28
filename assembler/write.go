package main

import (
	"fmt"
	"os"
	"strconv"
)

// writeACommand writes parser's current ACommand address value to the f based on the table
func writeACommand(f *os.File, p *Parser, table symbolTable) error {
	symbol := p.symbol()

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
	table[symbol] = p.getFreeRAMAddress()
	return writeAddress(f, table[symbol])
}

// writeAddress writes address to the f
func writeAddress(f *os.File, address uint16) error {
	_, err := fmt.Fprintf(f, "0%015b\n", address)
	return err
}

// writeCCommand writes parser's current CCommand to the f
func writeCCommand(f *os.File, p *Parser) error {
	_, err := fmt.Fprintf(f, "111%s%s%s\n", getCompBinary(p.comp()), getDestBinary(p.dest()), getJumpBinary(p.jump()))
	return err
}

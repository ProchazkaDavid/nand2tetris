package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("expected single .asm file")
	}

	if filepath.Ext(os.Args[1]) != ".asm" {
		log.Fatalln("expected .asm file")
	}

	if err := run(os.Args[1]); err != nil {
		log.Fatalln(err)
	}
}

func run(filename string) error {
	parser, err := newParser(filename)
	if err != nil {
		return fmt.Errorf("can't create a parser: %w", err)
	}

	// First pass - populates the symbol table
	table, err := parser.parseSymbols()
	if err != nil {
		return err
	}

	hackFile, err := os.Create(strings.TrimSuffix(filename, filepath.Ext(filename)) + ".hack")
	if err != nil {
		return fmt.Errorf("can't create .hack file: %w", err)
	}
	defer hackFile.Close()

	// Second pass - writes the final binary code into the .hack file
	for parser.HasMoreCommands() {
		parser.Advance()

		switch parser.commandType() {
		case ACommand:
			err = writeACommand(hackFile, parser, table)
		case CCommand:
			err = writeCCommand(hackFile, parser)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

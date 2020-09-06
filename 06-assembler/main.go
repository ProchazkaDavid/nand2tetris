package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProchazkaDavid/nand2tetris/06-assembler/parser"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("expected one argument")
	}

	// Create a parser
	filename := os.Args[1]
	p, err := parser.New(filename)
	if err != nil {
		log.Fatalln(err)
	}

	// First pass - populates the symbol table
	table, err := p.ParseSymbols()
	if err != nil {
		log.Fatalln(err)
	}

	// Creates .hack file
	hackFile, err := os.Create(strings.TrimSuffix(filename, filepath.Ext(filename)) + ".hack")
	if err != nil {
		log.Fatalf("can't create the .hack file: %v", err)
	}

	// Handle errors and closing of the .hack file
	defer func() {
		if err != nil {
			log.Println(err)
		}

		if closeError := hackFile.Close(); closeError != nil {
			log.Fatalf("can't save the .hack file: %v", closeError)
		}

		if err != nil {
			os.Exit(1)
		}
	}()

	// Second pass - writes the final binary code into the .hack file
	for p.HasMoreCommands() {
		p.Advance()

		switch p.CommandType() {
		case parser.ACommand:
			err = writeACommand(hackFile, p, table)
		case parser.CCommand:
			err = writeCCommand(hackFile, p)
		}

		if err != nil {
			return
		}
	}
}

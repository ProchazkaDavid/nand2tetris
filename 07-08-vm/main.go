package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ProchazkaDavid/nand2tetris/07-08-vm/code"
	"github.com/ProchazkaDavid/nand2tetris/07-08-vm/command"
	"github.com/ProchazkaDavid/nand2tetris/07-08-vm/parser"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("expected one argument")
	}

	fPath := os.Args[1]

	// Create a parser
	p, err := parser.New(fPath)
	if err != nil {
		log.Fatalln("couldn't setup the parser:", err)
	}

	// Prepare .asm file
	writer, err := code.New(fPath)
	if err != nil {
		log.Fatalln("couldn't setup the code writer:", err)
	}

	// Handle errors and closing of the .asm file
	defer func() {
		if err != nil {
			log.Println(err)
		}

		if closeError := writer.Close(); closeError != nil {
			log.Fatalf("couldn't save the .asm file: %v", closeError)
		}

		if err != nil {
			os.Exit(1)
		}
	}()

	filename := strings.TrimSuffix(path.Base(fPath), filepath.Ext(fPath))

	for p.HasMoreCommands() {
		p.Advance()

		switch p.CommandType() {
		case command.Push:
			err = writer.WritePush(p.Arg1(), p.Arg2(), filename)
		case command.Pop:
			err = writer.WritePop(p.Arg1(), p.Arg2(), filename)
		default:
			err = writer.WriteArithmetic(p.Arg1())
		}

		if err != nil {
			return
		}
	}
}

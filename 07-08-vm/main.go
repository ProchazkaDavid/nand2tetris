package main

import (
	"log"
	"os"
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

	path := os.Args[1]

	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("couldn't get info about the input: %v", err)
	}

	// Setup output file and files for parser
	output := strings.TrimSuffix(path, filepath.Ext(path)) + ".asm"
	files := []string{path}

	if info.IsDir() {
		files, err = filepath.Glob(filepath.Join(path, "*.vm"))
		if err != nil {
			log.Fatalf("couldn't get input files: %v", err)
		}

		output = filepath.Join(path, filepath.Base(path)+".asm")
	}

	writer, err := code.New(output)
	if err != nil {
		log.Fatalf("couldn't setup the code writer: %v", err)
	}

	if info.IsDir() {
		if err := writer.WriteInit(); err != nil {
			log.Fatalf("couldn't write the bootstrap code: %v", err)
		}
	}

	for _, file := range files {
		writer.SetFilename(file)
		if err := parse(file, writer); err != nil {
			log.Fatalf("couldn't parse %s: %v", file, err)
		}
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("couldn't save the .asm file: %v", err)
	}
}

func parse(file string, writer *code.Writer) error {
	p, err := parser.New(file)
	if err != nil {
		return err
	}

	currentFunction := ""
	for p.HasMoreCommands() {
		p.Advance()

		switch p.CommandType() {
		case command.Push:
			err = writer.WritePush(p.Arg1(), p.Arg2())
		case command.Pop:
			err = writer.WritePop(p.Arg1(), p.Arg2())
		case command.Label:
			err = writer.WriteLabel(p.Arg1(), currentFunction)
		case command.Goto:
			err = writer.WriteGoto(p.Arg1(), currentFunction)
		case command.If:
			err = writer.WriteIf(p.Arg1(), currentFunction)
		case command.Function:
			err = writer.WriteFunction(p.Arg1(), p.Arg2())
			currentFunction = p.Arg1()
		case command.Call:
			err = writer.WriteCall(p.Arg1(), p.Arg2())
		case command.Return:
			err = writer.WriteReturn()
		default:
			err = writer.WriteArithmetic(p.Arg1())
		}

		if err != nil {
			return err
		}
	}
	return nil
}

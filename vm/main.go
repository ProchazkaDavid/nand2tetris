package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProchazkaDavid/nand2tetris/vm/code"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("expected one argument - file or folder")
	}

	if err := run(os.Args[1]); err != nil {
		log.Fatalln(err)
	}
}

// run translates given file or folder
func run(path string) error {
	inputFileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("can't get info about the input: %w", err)
	}

	// Setup output file and files for parser
	ouputFilename := strings.TrimSuffix(path, filepath.Ext(path)) + ".asm"
	files := []string{path}

	inputIsDirectory := inputFileInfo.IsDir()

	// Given that the argument is a folder, setup parsing for every
	// file in this directory
	if inputIsDirectory {
		files, err = filepath.Glob(filepath.Join(path, "*.vm"))
		if err != nil {
			return fmt.Errorf("can't get input files: %w", err)
		}

		ouputFilename = filepath.Join(path, filepath.Base(path)+".asm")
	}

	outputFile, err := os.Create(ouputFilename)
	if err != nil {
		return fmt.Errorf("can't open the output file: %w", err)
	}
	defer outputFile.Close()

	writer := code.NewWriter(outputFile, ouputFilename)

	if inputIsDirectory {
		if err := writer.WriteInit(); err != nil {
			return fmt.Errorf("can't write the bootstrap code: %w", err)
		}
	}

	for _, file := range files {
		writer.SetFilename(file)

		if err := parseVMFile(file, writer); err != nil {
			return fmt.Errorf("can't parse %s: %w", file, err)
		}
	}

	return nil
}

func parseVMFile(file string, writer *code.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	vmParser, err := newParser(f)
	if err != nil {
		return err
	}

	currentFunction := ""
	for vmParser.hasMoreCommands() {
		vmParser.advance()

		switch vmParser.commandType() {
		case pushCmd:
			err = writer.WritePush(vmParser.firstArgument(), vmParser.secondArgument())
		case popCmd:
			err = writer.WritePop(vmParser.firstArgument(), vmParser.secondArgument())
		case labelCmd:
			err = writer.WriteLabel(vmParser.firstArgument(), currentFunction)
		case gotoCmd:
			err = writer.WriteGoto(vmParser.firstArgument(), currentFunction)
		case ifCmd:
			err = writer.WriteIf(vmParser.firstArgument(), currentFunction)
		case functionCmd:
			err = writer.WriteFunction(vmParser.firstArgument(), vmParser.secondArgument())
			currentFunction = vmParser.firstArgument()
		case callCmd:
			err = writer.WriteCall(vmParser.firstArgument(), vmParser.secondArgument())
		case returnCmd:
			err = writer.WriteReturn()
		default:
			err = writer.WriteArithmetic(vmParser.firstArgument())
		}

		if err != nil {
			return err
		}
	}

	return nil
}

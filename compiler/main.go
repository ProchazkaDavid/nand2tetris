package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProchazkaDavid/nand2tetris/compiler/compilation"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("expected one argument - file or folder")
	}

	if err := run(os.Args[1]); err != nil {
		log.Fatalln(err)
	}
}

// run compiles given file or folder
func run(path string) error {
	input, err := os.Open(path)
	if err != nil {
		return err
	}

	fileInfo, err := input.Stat()
	if err != nil {
		return fmt.Errorf("couldn't get info about the input: %w", err)
	}
	defer input.Close()

	files := []string{path}

	// Given that the argument is a folder, compile every file in this directory
	if fileInfo.IsDir() {
		files, err = filepath.Glob(filepath.Join(path, "*.jack"))
		if err != nil {
			return fmt.Errorf("couldn't get input files: %w", err)
		}
	}

	for _, file := range files {
		inputFile, err := os.Open(file)
		if err != nil {
			return err
		}

		vmFilename := strings.TrimSuffix(file, filepath.Ext(file)) + ".vm"
		vmOutput, err := os.Create(vmFilename)
		if err != nil {
			return errors.New("can't create the ouput file")
		}

		compilation.NewEngine(inputFile, vmOutput).CompileClass()

		vmOutput.Close()
		inputFile.Close()
	}

	return nil
}

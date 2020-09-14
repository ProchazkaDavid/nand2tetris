package code

import (
	"os"
	"path/filepath"
	"strings"
)

// Writer represents the
type Writer struct {
	file *os.File
}

// New opens the output file and gets ready to write into it.
func New(path string) (*Writer, error) {
	file, err := os.Create(strings.TrimSuffix(path, filepath.Ext(path)) + ".asm")
	return &Writer{file}, err
}

// WriteArithmetic writes to the output file the assembly code that implements the given arithmetic command.
func (cw *Writer) WriteArithmetic(operation string) error {
	switch operation {
	case "add", "sub", "and", "or":
		return cw.write(binaryOperation(operation))
	case "lt", "gt":
		return cw.write(compare(operation))
	case "neg", "not":
		return cw.write(unaryOperation(operation))
	default:
		return cw.write(eqInstructions())
	}
}

// WritePush writes to the output file the assembly code that implements Push/Pop command.
func (cw *Writer) WritePush(segment string, index int, filename string) error {
	switch segment {
	case "constant":
		return cw.write(constant(index))
	case "static":
		return cw.write(pushStatic(index, filename))
	case "temp":
		return cw.write(pushTemp(index))
	case "pointer":
		return cw.write(pushPointer(index))
	default:
		return cw.write(push(segment, index))
	}
}

// WritePop writes to the output file the assembly code that implements Push/Pop command.
func (cw *Writer) WritePop(segment string, index int, filename string) error {
	switch segment {
	case "static":
		return cw.write(popStatic(index, filename))
	case "temp":
		return cw.write(popTemp(index))
	case "pointer":
		return cw.write(popPointer(index))
	default:
		return cw.write(pop(segment, index))
	}
}

func (cw *Writer) write(instructions []string) error {
	var builder strings.Builder
	for _, inst := range instructions {
		builder.WriteString(inst)
		builder.WriteRune('\n')
	}

	_, err := cw.file.WriteString(builder.String())
	return err
}

// Close closes the ouput file
func (cw *Writer) Close() error { return cw.file.Close() }

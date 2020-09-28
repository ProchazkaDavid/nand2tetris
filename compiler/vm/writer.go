package vm

import (
	"fmt"
	"io"
)

// Writer produces stack based commands for the VM.
type Writer struct {
	output io.StringWriter
}

// NewWriter create a new output .vm file and prepares it for writing.
func NewWriter(output io.StringWriter) *Writer { return &Writer{output} }

// WritePush writes a VM push command.
func (w *Writer) WritePush(segment Segment, index int) {
	w.write(fmt.Sprintf("push %s %d", segment, index))
}

// WritePop writes a VM pop command.
func (w *Writer) WritePop(segment Segment, index int) {
	if segment == Constant {
		panic("can't pop the const segment")
	}

	w.write(fmt.Sprintf("pop %s %d", segment, index))
}

// WriteArithmetic writes a VM arithmetic-logical command.
func (w *Writer) WriteArithmetic(command string) {
	if command == "*" {
		w.WriteCall("Math.multiply", 2)
		return
	}
	if command == "/" {
		w.WriteCall("Math.divide", 2)
		return
	}

	value, ok := map[string]string{
		"+": "add",
		"-": "sub",
		"=": "eq",
		"<": "lt",
		">": "gt",
		"&": "and",
		"|": "or",
		"~": "not",
	}[command]
	if !ok {
		panic("unknown arithmetic command")
	}

	w.write(value)
}

// WriteUnaryOperation writes a VM '-' and '~' commands.
func (w *Writer) WriteUnaryOperation(operation string) {
	command := "neg"
	if operation != "-" {
		command = "not"
	}

	w.write(command)
}

// WriteLabel writes a VM label command.
func (w *Writer) WriteLabel(label string) {
	w.write(fmt.Sprintf("label %s", label))
}

// WriteGoto writes a VM goto command.
func (w *Writer) WriteGoto(label string) {
	w.write(fmt.Sprintf("goto %s", label))
}

// WriteIf writes a VM if-goto command.
func (w *Writer) WriteIf(label string) {
	w.write(fmt.Sprintf("if-goto %s", label))
}

// WriteCall writes a VM call command.
func (w *Writer) WriteCall(name string, args int) {
	w.write(fmt.Sprintf("call %s %d", name, args))
}

// WriteFunction writes a VM function command.
func (w *Writer) WriteFunction(name string, locals int) {
	w.write(fmt.Sprintf("function %s %d", name, locals))
}

// WriteReturn writes a VM return command.
func (w *Writer) WriteReturn() { w.write("return") }

// WriteString writes a string constant.
func (w *Writer) WriteString(input string) {
	w.WritePush(Constant, len(input))
	w.WriteCall("String.new", 1)

	for _, char := range input {
		w.WritePush(Constant, int(char))
		w.WriteCall("String.appendChar", 2)
	}
}

func (w *Writer) write(line string) {
	if _, err := w.output.WriteString(line + "\n"); err != nil {
		panic(fmt.Errorf("can't write to the file: %w", err))
	}
}

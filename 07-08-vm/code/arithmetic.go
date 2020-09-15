package code

import "fmt"

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

// instruction maps .vm arithmetic operation to the corresponding .asm operation
var instruction = map[string]string{
	"add": "+",
	"sub": "-",
	"neg": "-",
	"gt":  "GT",
	"lt":  "LT",
	"and": "&",
	"or":  "|",
	"not": "!",
}

// Generates instructions for neg, not
func unaryOperation(operation string) []string {
	return []string{
		fmt.Sprintf("// %s", operation),
		"@SP",
		"A=M-1",
		fmt.Sprintf("M=%sM", instruction[operation]),
	}
}

// Generates instructions for add, sub, and, or
func binaryOperation(operation string) []string {
	return []string{
		fmt.Sprintf("// %s", operation),
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		"A=A-1",
		fmt.Sprintf("M=M%sD", instruction[operation]),
	}
}

// counters for keeping the eq and compare operations unique across the .asm file
var eqCounter = 0
var compareCounter = 0

// Generates instructions for eq
func eqInstructions() []string {
	eqCounter++
	return []string{
		"// eq",
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		"A=A-1",
		"D=M-D",
		"M=-1",
		fmt.Sprintf("@EQ_%d", eqCounter),
		"D;JEQ",
		"@SP",
		"A=M-1",
		"M=0",
		fmt.Sprintf("(EQ_%d)", eqCounter),
	}
}

// Generates instructions for gt, lt
func compare(operation string) []string {
	compareCounter++
	return []string{
		fmt.Sprintf("// %s", operation),
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		"A=A-1",
		"D=M-D",
		"M=-1",
		fmt.Sprintf("@COMP_%d", compareCounter),
		fmt.Sprintf("D;J%s", instruction[operation]),
		"@SP",
		"A=M-1",
		"M=0",
		fmt.Sprintf("(COMP_%d)", compareCounter),
	}
}

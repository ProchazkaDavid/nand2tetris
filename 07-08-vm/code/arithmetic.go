package code

import "fmt"

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

// counters for keeping the compare and eq operations unique across the .vm file
var eqCounter = 0
var compareCounter = 0

// Generates instructions for neg, not
func unaryOperation(operation string) []string {
	return []string{
		fmt.Sprintf("// %sx", instruction[operation]),
		"@SP",
		"A=M-1",
		fmt.Sprintf("M=%sM", instruction[operation]),
	}
}

// Generates instructions for add, sub, and, or
func binaryOperation(operation string) []string {
	return []string{
		fmt.Sprintf("// x %s y", instruction[operation]),
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		"A=A-1",
		fmt.Sprintf("M=M%sD", instruction[operation]),
	}
}

// Generates instructions for eq
func eqInstructions() []string {
	eqCounter++
	return []string{
		"// x == y",
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
		fmt.Sprintf("// x %s y", instruction[operation]),
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

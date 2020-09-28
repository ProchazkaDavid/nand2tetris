package code

import "fmt"

// counter for keeping the call operations unique across the .asm file
var callCounter = 0

// WriteFunction writes function command to the the assembly file.
func (cw *Writer) WriteFunction(name string, variables int) error {
	instructions := []string{
		fmt.Sprintf("// function %s %d", name, variables),
		fmt.Sprintf("(%s)", name),
	}

	if variables > 0 {
		instructions = append(instructions, []string{
			"@0",
			"D=A",
		}...)
	}

	for i := 0; i < variables; i++ {
		instructions = append(instructions, []string{
			"@SP",
			"M=M+1",
			"A=M",
			"A=A-1",
			"M=D",
		}...)
	}

	return cw.write(instructions)
}

// WriteCall writes call command to the the assembly file.
func (cw *Writer) WriteCall(function string, arguments int) error {
	callCounter++

	instructions := []string{
		fmt.Sprintf("// call %s %d", function, arguments),
		"@SP",
		"D=M",
		"@R13",
		"M=D",
		fmt.Sprintf("@%s$ret.%d", function, callCounter),
		"D=A",
		"@SP",
		"M=M+1",
		"A=M",
		"A=A-1",
		"M=D",
	}

	for _, segment := range [...]string{"LCL", "ARG", "THIS", "THAT"} {
		instructions = append(instructions, []string{
			fmt.Sprintf("@%s", segment),
			"D=M",
			"@SP",
			"M=M+1",
			"A=M",
			"A=A-1",
			"M=D",
		}...)
	}

	return cw.write(append(instructions, []string{
		"@R13",
		"D=M",
		fmt.Sprintf("@%d", arguments),
		"D=D-A",
		"@ARG",
		"M=D",
		"@SP",
		"D=M",
		"@LCL",
		"M=D",
		fmt.Sprintf("@%s", function),
		"0;JMP",
		fmt.Sprintf("(%s$ret.%d)", function, callCounter),
	}...))
}

// WriteReturn writes return command to the the assembly file.
func (cw *Writer) WriteReturn() error {
	instructions := []string{
		"// return",
		"@LCL",
		"D=M",
		"@R13",
		"M=D",
		"@5",
		"D=D-A",
		"A=D",
		"D=M",
		"@R14",
		"M=D",
		"@SP",
		"A=M-1",
		"D=M",
		"@ARG",
		"A=M",
		"M=D",
		"@ARG",
		"D=M+1",
		"@SP",
		"M=D",
	}

	for _, segment := range [...]string{"THAT", "THIS", "ARG", "LCL"} {
		instructions = append(instructions, []string{
			"@R13",
			"M=M-1",
			"A=M",
			"D=M",
			fmt.Sprintf("@%s", segment),
			"M=D",
		}...)
	}

	return cw.write(append(instructions, []string{
		"@R14",
		"A=M",
		"0;JMP",
	}...))
}

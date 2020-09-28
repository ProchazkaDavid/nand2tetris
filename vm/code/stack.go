package code

import "fmt"

// segments maps .vm segment to the .asm address
var segments = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
}

// pointers maps .vm pointer to the .asm address
var pointers = map[int]string{
	0: "THIS",
	1: "THAT",
}

// WritePush writes to the output file the assembly code that implements Push/Pop command.
func (cw *Writer) WritePush(segment string, index int) error {
	switch segment {
	case "constant":
		return cw.write(constant(index))
	case "static":
		return cw.write(pushStatic(index, cw.filename))
	case "temp":
		return cw.write(pushTemp(index))
	case "pointer":
		return cw.write(pushPointer(index))
	default:
		return cw.write(push(segment, index))
	}
}

// WritePop writes to the output file the assembly code that implements Push/Pop command.
func (cw *Writer) WritePop(segment string, index int) error {
	switch segment {
	case "static":
		return cw.write(popStatic(index, cw.filename))
	case "temp":
		return cw.write(popTemp(index))
	case "pointer":
		return cw.write(popPointer(index))
	default:
		return cw.write(pop(segment, index))
	}
}

// Generates push instructions for local, argument, this, that
func push(segment string, index int) []string {
	return []string{
		fmt.Sprintf("// push %s %d", segment, index),
		fmt.Sprintf("@%d", index),
		"D=A",
		fmt.Sprintf("@%s", segments[segment]),
		"A=M",
		"A=D+A",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M",
		"A=A-1",
		"M=D",
	}
}

// Generates pop instructions for local, argument, this, that
func pop(segment string, index int) []string {
	return []string{
		fmt.Sprintf("// pop %s %d", segment, index),
		fmt.Sprintf("@%s", segments[segment]),
		"D=M",
		fmt.Sprintf("@%d", index),
		"D=D+A",
		"@R13",
		"M=D",
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		"@R13",
		"A=M",
		"M=D",
	}
}

// Generates push instructions for constant
func constant(index int) []string {
	return []string{
		fmt.Sprintf("// push constant %d", index),
		fmt.Sprintf("@%d", index),
		"D=A",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
}

// Generates push instructions for static
func pushStatic(index int, filename string) []string {
	return []string{
		fmt.Sprintf("// push static %d", index),
		fmt.Sprintf("@%s.%d", filename, index),
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
}

// Generates pop instructions for static
func popStatic(index int, filename string) []string {
	return []string{
		fmt.Sprintf("// pop static %d", index),
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		fmt.Sprintf("@%s.%d", filename, index),
		"M=D",
	}
}

// Generates push instructions for temp
func pushTemp(index int) []string {
	return []string{
		fmt.Sprintf("// push temp %d", index),
		"@5",
		"D=A",
		fmt.Sprintf("@%d", index),
		"A=D+A",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M",
		"A=A-1",
		"M=D",
	}
}

// Generates pop instructions for temp
func popTemp(index int) []string {
	return []string{
		fmt.Sprintf("// pop temp %d", index),
		"@5",
		"D=A",
		fmt.Sprintf("@%d", index),
		"D=D+A",
		"@R13",
		"M=D",
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		"@R13",
		"A=M",
		"M=D",
	}
}

// Generates push instructions for pointer
func pushPointer(index int) []string {
	return []string{
		fmt.Sprintf("// push pointer %d", index),
		fmt.Sprintf("@%s", pointers[index]),
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
}

// Generates pop instructions for pointer
func popPointer(index int) []string {
	return []string{
		fmt.Sprintf("// pop pointer %d", index),
		fmt.Sprintf("@%s", pointers[index]),
		"D=A",
		"@R13",
		"M=D",
		"@SP",
		"M=M-1",
		"@SP",
		"A=M",
		"D=M",
		"@R13",
		"A=M",
		"M=D",
	}
}

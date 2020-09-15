package code

import "fmt"

// WriteLabel writes label command to the the assembly file.
func (cw *Writer) WriteLabel(label, function string) error {
	return cw.write([]string{
		fmt.Sprintf("// label %s$%s", function, label),
		fmt.Sprintf("(%s$%s)", function, label),
	})
}

// WriteGoto writes goto command to the the assembly file.
func (cw *Writer) WriteGoto(label, function string) error {
	return cw.write([]string{
		fmt.Sprintf("// goto %s$%s", function, label),
		fmt.Sprintf("@%s$%s", function, label),
		"0;JMP",
	})
}

// WriteIf writes if-goto command to the the assembly file.
func (cw *Writer) WriteIf(label, function string) error {
	return cw.write([]string{
		fmt.Sprintf("// if-goto %s$%s", function, label),
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		fmt.Sprintf("@%s$%s", function, label),
		"D;JNE",
	})
}

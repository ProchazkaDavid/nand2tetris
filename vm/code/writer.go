package code

import (
	"io"
	"path"
	"path/filepath"
	"strings"
)

// Writer writes corresponding .asm instructions for VM commands
type Writer struct {
	output   io.StringWriter
	filename string
}

// NewWriter opens the output file and gets ready to write into it.
func NewWriter(output io.StringWriter, filename string) *Writer {
	return &Writer{
		output:   output,
		filename: strings.TrimSuffix(path.Base(filename), filepath.Ext(filename)),
	}
}

// SetFilename sets a new filename
func (cw *Writer) SetFilename(filename string) {
	cw.filename = strings.TrimSuffix(path.Base(filename), filepath.Ext(filename))
}

// WriteInit writes bootstrap code to the output file
func (cw *Writer) WriteInit() error {
	if err := cw.write([]string{
		"// Bootstrap code",
		"@256",
		"D=A",
		"@SP",
		"M=D",
	}); err != nil {
		return err
	}

	return cw.WriteCall("Sys.init", 0)
}

// write writes instructions to the file
func (cw *Writer) write(instructions []string) error {
	var builder strings.Builder

	for _, inst := range instructions {
		builder.WriteString(inst)
		builder.WriteRune('\n')
	}

	_, err := cw.output.WriteString(builder.String())
	return err
}

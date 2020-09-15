package code

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Writer represents the
type Writer struct {
	file     *os.File
	filename string
}

// New opens the output file and gets ready to write into it.
func New(fPath string) (*Writer, error) {
	file, err := os.Create(fPath)
	if err != nil {
		return nil, err
	}

	return &Writer{
		file:     file,
		filename: strings.TrimSuffix(path.Base(fPath), filepath.Ext(fPath)),
	}, err
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

	_, err := cw.file.WriteString(builder.String())
	return err
}

// Close closes the ouput file
func (cw *Writer) Close() error { return cw.file.Close() }

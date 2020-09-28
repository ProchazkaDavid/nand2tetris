package tokenizer

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
)

func isChararacter(char byte) bool { return unicode.IsLetter(rune(char)) }
func isNumber(char byte) bool      { return unicode.IsNumber(rune(char)) }
func isSpecial(char byte) bool     { return bytes.Contains(specialChars, []byte{char}) }

// parse parses at least one byte. The parsing stops when the isValid func
// returns false.
func (t *Tokenizer) parse(isValid func(char byte) bool) (string, error) {
	char, err := t.scanner.ReadByte()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteByte(char)

	chars, err := t.scanner.Peek(1)
	if err != nil {
		return "", err
	}

	for isValid(chars[0]) {
		chars[0], err = t.scanner.ReadByte()
		if err != nil {
			return "", err
		}

		builder.WriteByte(chars[0])

		chars, err = t.scanner.Peek(1)
		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}

// parseIndetifier parses identifier.
func (t *Tokenizer) parseIndetifier() (string, error) {
	return t.parse(func(char byte) bool {
		return isChararacter(char) || isNumber(char) || char == '_'
	})
}

// parseNumber parses number.
func (t *Tokenizer) parseNumber() (string, error) { return t.parse(isNumber) }

// peek peeks at one byte.
func peek(scanner *bufio.Reader) (byte, error) {
	chars, err := scanner.Peek(1)
	if err != nil {
		return 0, err
	}

	return chars[0], nil
}

// readByte reads a byte.
func readByte(scanner *bufio.Reader) error {
	_, err := scanner.ReadByte()
	return err
}

// eatByte reads a byte and returns if the reading was successful.
func eatByte(scanner *bufio.Reader) bool {
	return readByte(scanner) == nil
}

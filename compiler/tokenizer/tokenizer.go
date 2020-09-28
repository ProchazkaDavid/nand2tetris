package tokenizer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ProchazkaDavid/nand2tetris/compiler/token"
)

var specialChars = []byte{' ', '\n', '\r', '\t'}

// Tokenizer tokenizes .jack file, removes all white space and comments.
type Tokenizer struct {
	scanner *bufio.Reader
	token   string
}

// New opens the input .jack file and gets ready to tokenize it.
func New(file io.Reader) *Tokenizer { return &Tokenizer{bufio.NewReader(file), ""} }

// HasMoreTokens returns true if there are more tokens in the input.
func (t *Tokenizer) HasMoreTokens() bool {
	chars, err := t.scanner.Peek(1)
	if err != nil {
		return false
	}

	char := chars[0]

	var isComment bool
	var isDocsComment bool
	var previousChar byte
	var beforePreviousChar byte

	for isComment || isDocsComment || isSpecial(char) || char == '/' || char == '*' {
		if isComment {
			if char == '\n' {
				isComment = false
			}
		} else if isDocsComment {
			if previousChar == '*' && char == '/' {
				isDocsComment = false
			}
		} else {
			if beforePreviousChar == '/' && previousChar == '*' && char == '*' {
				isDocsComment = true
			} else if previousChar == '/' && char == '/' {
				isComment = true
			} else if char == '/' {
				chars, err := t.scanner.Peek(2)
				if err != nil {
					return false
				}
				if previousChar != '/' && !(chars[1] == '/' || chars[1] == '*') {
					break
				}
			} else if char == '*' {
				if previousChar != '/' {
					break
				}
			}
		}

		if !eatByte(t.scanner) {
			return false
		}

		beforePreviousChar = previousChar
		previousChar = char

		char, err = peek(t.scanner)
		if err != nil {
			return false
		}
	}

	return true
}

// Advance reads the next token from the input and makes it the current token.
// Should be called only if HasMoreTokens() is true.
// Initially there is no current command.
func (t *Tokenizer) Advance() error {
	chars, err := t.scanner.Peek(1)
	if err != nil {
		return err
	}

	if token.IsSymbol(string(chars[0])) {
		t.token = string(chars[0])
		return readByte(t.scanner)
	}

	switch {
	case isChararacter(chars[0]):
		id, err := t.parseIndetifier()
		if err != nil {
			return err
		}

		t.token = id

	case isNumber(chars[0]):
		number, err := t.parseNumber()
		if err != nil {
			return err
		}

		t.token = number

	case chars[0] == '"':
		if err := readByte(t.scanner); err != nil {
			return err
		}

		text, err := t.scanner.ReadString('"')
		if err != nil {
			return err
		}

		t.token = `"` + text

	default:
		return errors.New("unknown character")
	}

	return nil
}

// TokenType returns the type of the current token.
func (t *Tokenizer) TokenType() token.Type {
	_, err := strconv.ParseUint(t.token, 10, 15)

	switch {
	case token.IsKeyword(t.token):
		return token.Keyword
	case token.IsSymbol(t.token):
		return token.Symbol
	case strings.HasPrefix(t.token, `"`) && strings.HasSuffix(t.token, `"`):
		return token.StringConstant
	case err != nil:
		return token.Identifier
	default:
		return token.IntegerConstant
	}
}

// Keyword returns the keyword constant which is the current token.
// This method should be called only if TokenType() is token.Keyword.
func (t *Tokenizer) Keyword() token.KeywordType {
	for _, keyword := range token.Keywords {
		if t.token == string(keyword) {
			return keyword
		}
	}

	return token.Unknown
}

// Symbol returns the character which is the current token.
// This method should be called only if TokenType() is token.Symbol.
func (t *Tokenizer) Symbol() string { return t.token[:1] }

// Identifier returns the identifier which is the current token.
// This method should be called only if TokenType() is token.Identifier.
func (t *Tokenizer) Identifier() string { return t.token }

// IntValue returns the integer value of the current token.
// This method should be called only if TokenType() is token.IntConstant.
func (t *Tokenizer) IntValue() int {
	value, err := strconv.ParseUint(t.token, 10, 15)
	if err != nil {
		panic(fmt.Errorf("can't parse the integer value: %w", err))
	}

	return int(value)
}

// StringValue returns the string value of the current token,
// without the two enclosing double quotes.
// This method should be called only if TokenType() is token.StringConstant.
func (t *Tokenizer) StringValue() string { return t.token[1 : len(t.token)-1] }
